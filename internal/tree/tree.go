package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hosseinghorbani0/gowtree/internal/format"
	"github.com/hosseinghorbani0/gowtree/internal/platform"
	"github.com/hosseinghorbani0/gowtree/internal/style"
)

// Options controls tree traversal and display.
type Options struct {
	ShowAll        bool
	MaxDepth       int
	DirsOnly       bool
	FilesOnly      bool
	ShowSize       bool
	Ignore         []string
	SortBy         string
	SortRev        bool
	TimeMode       string
	FilesFirst     bool
	Regex          *regexp.Regexp
	FollowSymlinks bool
	DU             bool
	Verbose        bool
	Theme          style.Theme
}

// Node is the JSON-serializable tree structure.
type Node struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Size     int64   `json:"size,omitempty"`
	ModTime  string  `json:"mod_time,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

// Stats holds directory/file counts from a walk.
type Stats struct {
	Dirs  int32
	Files int32
}

func matchesRegex(regex *regexp.Regexp, name string) bool {
	if regex == nil {
		return true
	}
	return regex.MatchString(name)
}

func shouldInclude(name, fullPath string, opts Options) bool {
	if !opts.ShowAll && style.IsHidden(name, fullPath) {
		return false
	}
	if !matchesRegex(opts.Regex, name) {
		return false
	}
	for _, pat := range opts.Ignore {
		if m, _ := filepath.Match(pat, name); m {
			return false
		}
	}
	return true
}

type sortableEntry struct {
	entry os.DirEntry
	info  os.FileInfo
}

func SortEntries(entries []os.DirEntry, opts Options) []os.DirEntry {
	list := make([]sortableEntry, len(entries))
	for i, e := range entries {
		info, _ := e.Info()
		list[i] = sortableEntry{entry: e, info: info}
	}

	dirsFirst := !opts.FilesFirst
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].entry.IsDir() != list[j].entry.IsDir() {
			if dirsFirst {
				return list[i].entry.IsDir()
			}
			return !list[i].entry.IsDir()
		}
		switch opts.SortBy {
		case "size":
			si, sj := int64(0), int64(0)
			if list[i].info != nil {
				si = list[i].info.Size()
			}
			if list[j].info != nil {
				sj = list[j].info.Size()
			}
			if si != sj {
				return si < sj
			}
		case "time":
			ti, tj := time.Time{}, time.Time{}
			if list[i].info != nil {
				ti = platform.FileTime(list[i].info, opts.TimeMode)
			}
			if list[j].info != nil {
				tj = platform.FileTime(list[j].info, opts.TimeMode)
			}
			if !ti.Equal(tj) {
				return ti.Before(tj)
			}
		}
		return strings.ToLower(list[i].entry.Name()) < strings.ToLower(list[j].entry.Name())
	})

	if opts.SortRev {
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
	}

	result := make([]os.DirEntry, len(list))
	for i, se := range list {
		result[i] = se.entry
	}
	return result
}

func IsExecutable(name string) bool {
	name = strings.ToLower(name)
	return strings.HasSuffix(name, ".exe") ||
		strings.HasSuffix(name, ".bat") ||
		strings.HasSuffix(name, ".cmd") ||
		strings.HasSuffix(name, ".sh") ||
		strings.HasSuffix(name, ".ps1")
}

func DirSize(path string, verbose bool) int64 {
	var size int64
	filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func CountItems(root string, opts Options) int64 {
	var total int64
	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !shouldInclude(info.Name(), p, opts) {
			return nil
		}
		if opts.DirsOnly && !info.IsDir() {
			return nil
		}
		if opts.FilesOnly && info.IsDir() {
			return nil
		}
		total++
		return nil
	})
	return total
}

func BuildJSON(path string, depth int, visited map[string]bool, opts Options) *Node {
	if visited == nil {
		visited = make(map[string]bool)
	}
	absPath, _ := filepath.Abs(path)
	if visited[absPath] {
		return &Node{Name: filepath.Base(path), Type: "symlink-loop"}
	}
	visited[absPath] = true

	node := &Node{Name: filepath.Base(path)}
	info, err := os.Stat(path)
	if err != nil {
		node.Type = "error"
		return node
	}

	if info.IsDir() {
		node.Type = "directory"
		if opts.DU {
			node.Size = DirSize(path, opts.Verbose)
		} else {
			node.Size = info.Size()
		}
		if opts.TimeMode != "" {
			node.ModTime = platform.FileTime(info, opts.TimeMode).Format(opts.Theme.TimeFormat)
		}
		if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
			return node
		}
		entries, _ := os.ReadDir(path)
		for _, e := range entries {
			if !shouldInclude(e.Name(), filepath.Join(path, e.Name()), opts) {
				continue
			}
			if opts.DirsOnly && !e.IsDir() {
				continue
			}
			if opts.FilesOnly && e.IsDir() {
				continue
			}

			entryPath := filepath.Join(path, e.Name())
			if opts.FollowSymlinks && e.Type()&os.ModeSymlink != 0 {
				resolved, err := filepath.EvalSymlinks(entryPath)
				if err == nil {
					info2, err2 := os.Stat(resolved)
					if err2 == nil && info2.IsDir() {
						child := BuildJSON(resolved, depth+1, visited, opts)
						child.Name = e.Name()
						node.Children = append(node.Children, child)
						continue
					}
				}
			}
			child := BuildJSON(entryPath, depth+1, visited, opts)
			node.Children = append(node.Children, child)
		}
		sortChildren(node, opts)
	} else {
		node.Type = "file"
		node.Size = info.Size()
		if opts.TimeMode != "" {
			node.ModTime = platform.FileTime(info, opts.TimeMode).Format(opts.Theme.TimeFormat)
		}
	}
	return node
}

func sortChildren(node *Node, opts Options) {
	type childItem struct {
		child *Node
		time  time.Time
	}
	var childList []childItem
	for _, c := range node.Children {
		var t time.Time
		if c.ModTime != "" {
			t, _ = time.Parse(opts.Theme.TimeFormat, c.ModTime)
		}
		childList = append(childList, childItem{c, t})
	}
	sort.SliceStable(childList, func(i, j int) bool {
		ci, cj := childList[i].child, childList[j].child
		dirsFirst := !opts.FilesFirst
		if ci.Type != cj.Type {
			if dirsFirst {
				return ci.Type == "directory"
			}
			return ci.Type != "directory"
		}
		switch opts.SortBy {
		case "size":
			if ci.Size != cj.Size {
				return ci.Size < cj.Size
			}
		case "time":
			if !childList[i].time.Equal(childList[j].time) {
				return childList[i].time.Before(childList[j].time)
			}
		}
		return strings.ToLower(ci.Name) < strings.ToLower(cj.Name)
	})
	if opts.SortRev {
		for i, j := 0, len(childList)-1; i < j; i, j = i+1, j-1 {
			childList[i], childList[j] = childList[j], childList[i]
		}
	}
	node.Children = nil
	for _, item := range childList {
		node.Children = append(node.Children, item.child)
	}
}

// Print writes an ASCII/Unicode tree to w and updates stats.
func Print(w io.Writer, path, prefix string, depth int, progCount *int64, visited map[string]bool, stats *Stats, opts Options) {
	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		return
	}
	absPath, _ := filepath.Abs(path)
	if visited[absPath] {
		return
	}
	visited[absPath] = true

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	var filtered []os.DirEntry
	for _, e := range entries {
		if !shouldInclude(e.Name(), filepath.Join(path, e.Name()), opts) {
			continue
		}
		if opts.DirsOnly && !e.IsDir() {
			continue
		}
		if opts.FilesOnly && e.IsDir() {
			continue
		}
		filtered = append(filtered, e)
	}

	sorted := SortEntries(filtered, opts)
	br := opts.Theme.Branches
	col := opts.Theme.Colors

	for i, entry := range sorted {
		if progCount != nil {
			atomic.AddInt64(progCount, 1)
		}

		isLast := i == len(sorted)-1
		connector := br.Middle
		nextPrefix := prefix + br.Pipe
		if isLast {
			connector = br.Last
			nextPrefix = prefix + br.Space
		}

		display := entry.Name()
		info, infoErr := entry.Info()

		if opts.Theme.UseIcons {
			display = style.IconFor(entry.Name(), entry.IsDir(), entry.Type()&os.ModeSymlink != 0) + display
		}

		if opts.ShowSize && !entry.IsDir() && infoErr == nil {
			display += " (" + format.Size(info.Size()) + ")"
		}
		if opts.TimeMode != "" && infoErr == nil {
			timeStr := platform.FileTime(info, opts.TimeMode).Format(opts.Theme.TimeFormat)
			display += " " + col.Yellow + timeStr + col.Reset
		}
		if opts.DU && entry.IsDir() {
			dirSize := DirSize(filepath.Join(path, entry.Name()), opts.Verbose)
			display += " (" + format.Size(dirSize) + ")"
		}

		if entry.IsDir() {
			display = col.Blue + display + col.Reset
		} else if entry.Type()&os.ModeSymlink != 0 {
			display = col.Cyan + display + col.Reset
		} else if IsExecutable(entry.Name()) {
			display = col.Green + display + col.Reset
		}

		_, _ = fmt.Fprintln(w, prefix+connector+display)

		isSymlinkDir := false
		var targetPath string
		if opts.FollowSymlinks && entry.Type()&os.ModeSymlink != 0 {
			fullPath := filepath.Join(path, entry.Name())
			resolved, err := filepath.EvalSymlinks(fullPath)
			if err == nil {
				info2, err2 := os.Stat(resolved)
				if err2 == nil && info2.IsDir() {
					isSymlinkDir = true
					targetPath = resolved
				}
			}
		}

		if isSymlinkDir {
			atomic.AddInt32(&stats.Dirs, 1)
			Print(w, targetPath, nextPrefix, depth+1, progCount, visited, stats, opts)
		} else if entry.IsDir() {
			atomic.AddInt32(&stats.Dirs, 1)
			Print(w, filepath.Join(path, entry.Name()), nextPrefix, depth+1, progCount, visited, stats, opts)
		} else {
			atomic.AddInt32(&stats.Files, 1)
		}
	}
}
