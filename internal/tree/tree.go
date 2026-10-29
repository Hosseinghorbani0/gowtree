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
	SizeCache      map[string]int64
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

type entryKind string

const (
	entryFile    entryKind = "file"
	entryDir     entryKind = "directory"
	entrySymlink entryKind = "symlink"
)

type classifiedEntry struct {
	kind          entryKind
	info          os.FileInfo
	targetInfo    os.FileInfo
	canonicalPath string
}

func classifyEntry(path string) (classifiedEntry, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return classifiedEntry{}, err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return classifiedEntry{
				kind:          entrySymlink,
				info:          info,
				canonicalPath: path,
			}, nil
		}
		targetInfo, err := os.Stat(resolved)
		if err != nil {
			return classifiedEntry{
				kind:          entrySymlink,
				info:          info,
				canonicalPath: resolved,
			}, nil
		}
		return classifiedEntry{
			kind:          entrySymlink,
			info:          info,
			targetInfo:    targetInfo,
			canonicalPath: resolved,
		}, nil
	}

	kind := entryFile
	if info.IsDir() {
		kind = entryDir
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	return classifiedEntry{
		kind:          kind,
		info:          info,
		canonicalPath: absPath,
	}, nil
}

func SortEntries(entries []os.DirEntry, opts Options) []os.DirEntry {
	sorted, _ := SortEntriesStrict(entries, opts)
	return sorted
}

// SortEntriesStrict sorts entries and reports metadata lookup failures.
func SortEntriesStrict(entries []os.DirEntry, opts Options) ([]os.DirEntry, error) {
	list := make([]sortableEntry, len(entries))
	for i, e := range entries {
		info, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("read metadata for %q: %w", e.Name(), err)
		}
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
	return result, nil
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
	size, _ := DirSizeWithError(path)
	return size
}

func dirSize(path string, cache map[string]int64) int64 {
	size, _ := dirSizeWithError(path, cache)
	return size
}

// DirSizeWithError calculates the size of regular files below path.
// It does not follow directory symlinks, matching filepath.Walk semantics.
func DirSizeWithError(path string) (int64, error) {
	return dirSizeWithError(path, nil)
}

func dirSizeWithError(path string, cache map[string]int64) (int64, error) {
	cacheKey, err := filepath.Abs(path)
	if err != nil {
		cacheKey = path
	}
	if cache != nil {
		if size, ok := cache[cacheKey]; ok {
			return size, nil
		}
	}
	var size int64
	err = filepath.Walk(path, func(walkPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info == nil {
			return fmt.Errorf("missing metadata for %q", walkPath)
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("calculate directory size for %q: %w", path, err)
	}
	if cache != nil {
		cache[cacheKey] = size
	}
	return size, nil
}

func CountItems(root string, opts Options) int64 {
	total, _ := CountItemsStrict(root, opts)
	return total
}

// CountItemsStrict counts exactly the entries traversed by Print.
func CountItemsStrict(root string, opts Options) (int64, error) {
	return countItems(root, 0, make(map[string]bool), opts)
}

func countItems(path string, depth int, visited map[string]bool, opts Options) (int64, error) {
	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		return 0, nil
	}
	entry, err := classifyEntry(path)
	if err != nil {
		return 0, fmt.Errorf("classify %q: %w", path, err)
	}
	if entry.kind == entrySymlink && !opts.FollowSymlinks {
		return 0, nil
	}
	if visited[entry.canonicalPath] {
		return 0, nil
	}
	visited[entry.canonicalPath] = true

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("read directory %q: %w", path, err)
	}

	var total int64
	for _, e := range entries {
		entryPath := filepath.Join(path, e.Name())
		if !shouldInclude(e.Name(), entryPath, opts) {
			continue
		}
		if opts.DirsOnly && !e.IsDir() {
			continue
		}
		if opts.FilesOnly && e.IsDir() {
			continue
		}

		total++
		classified, classifyErr := classifyEntry(entryPath)
		if classifyErr != nil {
			continue
		}
		isDirectory := classified.kind == entryDir ||
			(classified.kind == entrySymlink && opts.FollowSymlinks &&
				classified.targetInfo != nil && classified.targetInfo.IsDir())
		if isDirectory && e.IsDir() {
			childTotal, childErr := countItems(entryPath, depth+1, visited, opts)
			if childErr != nil {
				return total, childErr
			}
			total += childTotal
		} else if isDirectory && classified.kind == entrySymlink {
			childTotal, childErr := countItems(classified.canonicalPath, depth+1, visited, opts)
			if childErr != nil {
				return total, childErr
			}
			total += childTotal
		}
	}
	return total, nil
}

func BuildJSON(path string, depth int, visited map[string]bool, opts Options) *Node {
	node, _ := BuildJSONStrict(path, depth, visited, opts)
	return node
}

// BuildJSONStrict builds a JSON tree and reports filesystem traversal errors.
func BuildJSONStrict(path string, depth int, visited map[string]bool, opts Options) (*Node, error) {
	if visited == nil {
		visited = make(map[string]bool)
	}
	node := &Node{Name: filepath.Base(path)}
	entry, err := classifyEntry(path)
	if err != nil {
		return nil, fmt.Errorf("classify %q: %w", path, err)
	}
	if entry.kind == entrySymlink && !opts.FollowSymlinks {
		node.Type = string(entrySymlink)
		return node, nil
	}
	if entry.kind == entrySymlink && entry.targetInfo == nil {
		node.Type = "broken-symlink"
		return node, nil
	}

	canonicalPath := entry.canonicalPath
	if visited[canonicalPath] {
		return &Node{Name: filepath.Base(path), Type: "symlink-loop"}, nil
	}
	visited[canonicalPath] = true

	info := entry.info
	if entry.kind == entrySymlink && entry.targetInfo != nil {
		info = entry.targetInfo
	}

	if info.IsDir() {
		node.Type = "directory"
		if opts.DU {
			node.Size, err = dirSizeWithError(path, opts.SizeCache)
			if err != nil {
				return nil, fmt.Errorf("calculate directory size for %q: %w", path, err)
			}
		} else {
			node.Size = info.Size()
		}
		if opts.TimeMode != "" {
			node.ModTime = platform.FileTime(info, opts.TimeMode).Format(opts.Theme.TimeFormat)
		}
		if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
			return node, nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("read directory %q: %w", path, err)
		}
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
			classified, err := classifyEntry(entryPath)
			if err != nil {
				return nil, fmt.Errorf("classify %q: %w", entryPath, err)
			}
			if classified.kind == entrySymlink && !opts.FollowSymlinks {
				node.Children = append(node.Children, &Node{Name: e.Name(), Type: string(entrySymlink)})
				continue
			}
			if classified.kind == entrySymlink && classified.targetInfo == nil {
				node.Children = append(node.Children, &Node{Name: e.Name(), Type: "broken-symlink"})
				continue
			}
			childPath := entryPath
			child, err := BuildJSONStrict(childPath, depth+1, visited, opts)
			if err != nil {
				return nil, err
			}
			if classified.kind == entrySymlink {
				child.Name = e.Name()
			}
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
	return node, nil
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
func Print(w io.Writer, path, prefix string, depth int, progCount *int64, visited map[string]bool, stats *Stats, opts Options) error {
	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		return nil
	}
	rootEntry, err := classifyEntry(path)
	if err != nil {
		return fmt.Errorf("classify %q: %w", path, err)
	}
	if rootEntry.kind == entrySymlink && !opts.FollowSymlinks {
		return nil
	}
	canonicalPath := rootEntry.canonicalPath
	if visited[canonicalPath] {
		return nil
	}
	visited[canonicalPath] = true

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", path, err)
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

	sorted, err := SortEntriesStrict(filtered, opts)
	if err != nil {
		return err
	}
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
		classified, classifyErr := classifyEntry(filepath.Join(path, entry.Name()))
		info, infoErr := entry.Info()
		if infoErr != nil {
			return fmt.Errorf("read metadata for %q: %w", entry.Name(), infoErr)
		}

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
			size, err := dirSizeWithError(filepath.Join(path, entry.Name()), opts.SizeCache)
			if err != nil {
				return fmt.Errorf("calculate directory size for %q: %w", entry.Name(), err)
			}
			display += " (" + format.Size(size) + ")"
		}

		if classified.kind == entryDir {
			display = col.Blue + display + col.Reset
		} else if classified.kind == entrySymlink {
			display = col.Cyan + display + col.Reset
		} else if IsExecutable(entry.Name()) {
			display = col.Green + display + col.Reset
		}

		if _, err := fmt.Fprintln(w, prefix+connector+display); err != nil {
			return fmt.Errorf("write entry %q: %w", entry.Name(), err)
		}

		isSymlinkDir := classifyErr == nil &&
			classified.kind == entrySymlink &&
			opts.FollowSymlinks &&
			classified.targetInfo != nil &&
			classified.targetInfo.IsDir()
		if isSymlinkDir {
			atomic.AddInt32(&stats.Dirs, 1)
			if err := Print(w, classified.canonicalPath, nextPrefix, depth+1, progCount, visited, stats, opts); err != nil {
				return err
			}
		} else if classifyErr == nil && classified.kind == entryDir {
			atomic.AddInt32(&stats.Dirs, 1)
			if err := Print(w, filepath.Join(path, entry.Name()), nextPrefix, depth+1, progCount, visited, stats, opts); err != nil {
				return err
			}
		} else {
			atomic.AddInt32(&stats.Files, 1)
		}
	}
	return nil
}
