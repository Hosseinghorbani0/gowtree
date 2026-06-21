package app

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hosseinghorbani0/gowtree/internal/config"
	"github.com/hosseinghorbani0/gowtree/internal/metadata"
	"github.com/hosseinghorbani0/gowtree/internal/output"
	"github.com/hosseinghorbani0/gowtree/internal/style"
	"github.com/hosseinghorbani0/gowtree/internal/tree"
)

type ignoreList []string

func (i *ignoreList) String() string { return strings.Join(*i, ", ") }
func (i *ignoreList) Set(v string) error {
	*i = append(*i, v)
	return nil
}

// Run is the main application entry point.
func Run(args []string) int {
	fs := flag.NewFlagSet("gowtree", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		showAll        bool
		maxDepth       int
		dirsOnly       bool
		filesOnly      bool
		showSize       bool
		ignorePats     ignoreList
		colorMode      string
		noColor        bool
		charset        string
		jsonOut        bool
		duFlag         bool
		verbose        bool
		progress       bool
		showVersion    bool
		showAuthor     bool
		noCredit       bool
		outFile        string
		clipFlag       bool
		followSymlinks bool
		sortBy         string
		sortRev        bool
		timeFlag       string
		filesFirst     bool
		regexPat       string
		icons          bool
		markdown       bool
		htmlOut        bool
		configFile     string
	)

	fs.BoolVar(&showAll, "a", false, "Show hidden files and directories")
	fs.IntVar(&maxDepth, "L", 0, "Maximum depth (0 = unlimited)")
	fs.BoolVar(&dirsOnly, "d", false, "Show directories only")
	fs.BoolVar(&filesOnly, "f", false, "Show files only")
	fs.BoolVar(&showSize, "s", false, "Show file sizes")
	fs.Var(&ignorePats, "I", "Ignore pattern (can be used multiple times)")
	fs.StringVar(&colorMode, "color", "auto", "Color mode: auto, always, never")
	fs.BoolVar(&noColor, "no-color", false, "Disable colors (same as --color=never)")
	fs.StringVar(&charset, "charset", "utf8", "Character set: utf8 or ascii")
	fs.BoolVar(&jsonOut, "J", false, "Output as JSON")
	fs.BoolVar(&duFlag, "du", false, "Show cumulative directory sizes")
	fs.BoolVar(&verbose, "v", false, "Verbose output")
	fs.BoolVar(&progress, "p", false, "Show progress bar (implies first scan)")
	fs.BoolVar(&showVersion, "version", false, "Print version and exit")
	fs.BoolVar(&showAuthor, "author", false, "Print author info and exit")
	fs.BoolVar(&noCredit, "no-credit", false, "Hide the signature line at the end")
	fs.StringVar(&outFile, "out", "", "Write output to a file")
	fs.BoolVar(&clipFlag, "clip", false, "Copy output to clipboard")
	fs.BoolVar(&followSymlinks, "l", false, "Follow symbolic links")
	fs.BoolVar(&followSymlinks, "follow-symlinks", false, "Follow symbolic links")
	fs.StringVar(&sortBy, "sort", "name", "Sort by: name, size, time, none")
	fs.BoolVar(&sortRev, "r", false, "Reverse sort order")
	fs.StringVar(&timeFlag, "time", "", "Show file time: mod, change, create")
	fs.BoolVar(&filesFirst, "files-first", false, "List files before directories")
	fs.StringVar(&regexPat, "R", "", "Filter by regex pattern")
	fs.BoolVar(&icons, "icons", false, "Show Nerd Font icons")
	fs.BoolVar(&markdown, "markdown", false, "Output in Markdown format")
	fs.BoolVar(&htmlOut, "html", false, "Output in HTML format")
	fs.StringVar(&configFile, "config", "", "Path to config file (default: ~/.gowtree.yaml)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "🌳 gowtree %s — A modern tree command for Windows\n", metadata.Version)
		fmt.Fprintf(fs.Output(), "Author: %s <%s>\n\n", metadata.Author, metadata.RepoURL)
		fmt.Fprintf(fs.Output(), "Usage:\n  gowtree [path] [flags]\n\n")
		fmt.Fprintf(fs.Output(), "Examples:\n")
		fmt.Fprintf(fs.Output(), "  gowtree\n")
		fmt.Fprintf(fs.Output(), "  gowtree -a -s -L 2 --icons\n")
		fmt.Fprintf(fs.Output(), "  gowtree -R \"\\.go$\" --sort size -r --markdown\n")
		fmt.Fprintf(fs.Output(), "  gowtree --html --out tree.html\n\n")
		fmt.Fprintf(fs.Output(), "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nReport bugs at: https://%s/issues\n", metadata.RepoURL)
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}

	userFlags := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		userFlags[f.Name] = true
	})

	applyConfig := func(name string, apply func()) {
		if !userFlags[name] {
			apply()
		}
	}

	cfgPath := configFile
	if cfgPath == "" {
		cfgPath = config.DefaultPath()
	}
	if cfg, err := config.Load(cfgPath); err == nil && cfg != nil {
		if cfg.Color != nil {
			applyConfig("color", func() { colorMode = *cfg.Color })
		}
		if cfg.NoColor != nil {
			applyConfig("no-color", func() { noColor = *cfg.NoColor })
		}
		if cfg.Charset != nil {
			applyConfig("charset", func() { charset = *cfg.Charset })
		}
		if cfg.ShowAll != nil {
			applyConfig("a", func() { showAll = *cfg.ShowAll })
		}
		if cfg.MaxDepth != nil {
			applyConfig("L", func() { maxDepth = *cfg.MaxDepth })
		}
		if cfg.DirsOnly != nil {
			applyConfig("d", func() { dirsOnly = *cfg.DirsOnly })
		}
		if cfg.FilesOnly != nil {
			applyConfig("f", func() { filesOnly = *cfg.FilesOnly })
		}
		if cfg.ShowSize != nil {
			applyConfig("s", func() { showSize = *cfg.ShowSize })
		}
		if cfg.Ignore != nil {
			applyConfig("I", func() { ignorePats = append(ignorePats, *cfg.Ignore...) })
		}
		if cfg.SortBy != nil {
			applyConfig("sort", func() { sortBy = *cfg.SortBy })
		}
		if cfg.SortRev != nil {
			applyConfig("r", func() { sortRev = *cfg.SortRev })
		}
		if cfg.TimeFlag != nil {
			applyConfig("time", func() { timeFlag = *cfg.TimeFlag })
		}
		if cfg.FilesFirst != nil {
			applyConfig("files-first", func() { filesFirst = *cfg.FilesFirst })
		}
		if cfg.Regex != nil {
			applyConfig("R", func() { regexPat = *cfg.Regex })
		}
		if cfg.Icons != nil {
			applyConfig("icons", func() { icons = *cfg.Icons })
		}
		if cfg.FollowLinks != nil {
			applyConfig("l", func() { followSymlinks = *cfg.FollowLinks })
		}
		if cfg.DU != nil {
			applyConfig("du", func() { duFlag = *cfg.DU })
		}
	}

	if showVersion {
		fmt.Println("gowtree", metadata.Version)
		return 0
	}
	if showAuthor {
		fmt.Print(`  🌳  gowtree
       │
   ────┼────
       │
  crafted by Hossein Ghorbani
  github.com/hosseinghorbani0/gowtree

  "In every tree of files, there is a hidden forest of ideas."
`)
		return 0
	}

	theme := style.NewTheme(colorMode, noColor, charset, timeFlag, icons)

	var compiledRegex *regexp.Regexp
	if regexPat != "" {
		var err error
		compiledRegex, err = regexp.Compile(regexPat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid regex pattern: %v\n", err)
			return 1
		}
	}

	opts := tree.Options{
		ShowAll:        showAll,
		MaxDepth:       maxDepth,
		DirsOnly:       dirsOnly,
		FilesOnly:      filesOnly,
		ShowSize:       showSize,
		Ignore:         ignorePats,
		SortBy:         sortBy,
		SortRev:        sortRev,
		TimeMode:       timeFlag,
		FilesFirst:     filesFirst,
		Regex:          compiledRegex,
		FollowSymlinks: followSymlinks,
		DU:             duFlag,
		Verbose:        verbose,
		Theme:          theme,
	}

	var clipBuffer *bytes.Buffer
	var writers []io.Writer
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating output file:", err)
			return 1
		}
		defer func() {
			_ = f.Sync()
			_ = f.Close()
		}()
		writers = append(writers, f)
	}
	if clipFlag {
		clipBuffer = &bytes.Buffer{}
		writers = append(writers, clipBuffer)
	}

	var outWriter io.Writer = os.Stdout
	if len(writers) > 0 {
		if outFile == "" && clipFlag {
			writers = append(writers, os.Stdout)
		}
		if len(writers) == 1 {
			outWriter = writers[0]
		} else {
			outWriter = io.MultiWriter(writers...)
		}
	}

	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	resolveRoot := followSymlinks
	if !resolveRoot {
		if fi, err := os.Lstat(root); err == nil && fi.Mode()&os.ModeSymlink != 0 {
			resolveRoot = true
		}
	}
	if resolveRoot {
		if resolved, err := filepath.EvalSymlinks(root); err == nil {
			root = resolved
		}
	}

	info, err := os.Stat(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening path:", err)
		return 1
	}
	if !info.IsDir() {
		fmt.Fprintln(outWriter, root)
		return 0
	}

	absRoot, _ := filepath.Abs(root)

	if jsonOut {
		node := tree.BuildJSON(root, 0, nil, opts)
		node.Name = strings.TrimRight(absRoot, "/\\") + "/"
		enc := json.NewEncoder(outWriter)
		enc.SetIndent("", "  ")
		_ = enc.Encode(node)
		return 0
	}

	if markdown || htmlOut {
		node := tree.BuildJSON(root, 0, nil, opts)
		node.Name = strings.TrimRight(absRoot, "/\\") + "/"
		if markdown {
			output.RenderMarkdown(outWriter, node, "", true, theme, duFlag)
		} else {
			output.RenderHTML(outWriter, node, "", true, theme, duFlag)
		}
		if !noCredit {
			fmt.Fprintf(outWriter, "\n🌳 gowtree — crafted by %s (%s)\n", metadata.Author, metadata.RepoURL)
		}
		if clipFlag && clipBuffer != nil {
			_ = output.CopyToClipboard(clipBuffer.String())
		}
		return 0
	}

	fmt.Fprintln(outWriter, strings.TrimRight(absRoot, "/\\")+"/")
	stats := &tree.Stats{}

	var progressTotal int64
	var progressCurrent int64
	if progress && style.IsTerminal(os.Stderr.Fd()) {
		progressTotal = tree.CountItems(root, opts)
		if progressTotal > 0 {
			go func() {
				for {
					cur := atomic.LoadInt64(&progressCurrent)
					pct := cur * 100 / progressTotal
					bar := strings.Repeat("#", int(pct)/2) + strings.Repeat(" ", 50-int(pct)/2)
					fmt.Fprintf(os.Stderr, "\r[%s] %d%% (%d/%d)", bar, pct, cur, progressTotal)
					if cur >= progressTotal {
						fmt.Fprint(os.Stderr, "\r"+strings.Repeat(" ", 80)+"\r")
						return
					}
					time.Sleep(100 * time.Millisecond)
				}
			}()
		}
	}

	visited := make(map[string]bool)
	tree.Print(outWriter, root, "", 0, &progressCurrent, visited, stats, opts)

	fmt.Fprintf(outWriter, "\n%d directories, %d files\n", stats.Dirs, stats.Files)
	if progress && style.IsTerminal(os.Stderr.Fd()) {
		fmt.Fprint(os.Stderr, "\r"+strings.Repeat(" ", 80)+"\r")
	}
	if !noCredit {
		fmt.Fprintf(outWriter, "\n🌳 gowtree — crafted by %s (%s)\n", metadata.Author, metadata.RepoURL)
	}
	if clipFlag && clipBuffer != nil {
		_ = output.CopyToClipboard(clipBuffer.String())
	}
	return 0
}
