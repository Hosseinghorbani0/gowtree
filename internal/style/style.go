package style

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"

	"github.com/hosseinghorbani0/gowtree/internal/platform"
)

// Palette holds ANSI color codes (empty when colors disabled).
type Palette struct {
	Reset  string
	Blue   string
	Cyan   string
	Green  string
	Yellow string
}

// Branches holds tree-drawing characters.
type Branches struct {
	Middle string
	Last   string
	Pipe   string
	Space  string
}

// Theme combines colors, branch chars, and time format for one run.
type Theme struct {
	Colors     Palette
	Branches   Branches
	TimeFormat string
	UseIcons   bool
}

var iconMap = map[string]string{
	"dir":      "\uf07c",
	"dir_open": "\uf07c",
	".go":      "\ue627",
	".py":      "\ue606",
	".js":      "\ue74e",
	".ts":      "\ue628",
	".html":    "\uf13b",
	".css":     "\ue749",
	".json":    "\ue60b",
	".md":      "\uf48a",
	".txt":     "\uf15c",
	".pdf":     "\uf1c1",
	".zip":     "\uf410",
	".tar":     "\uf410",
	".gz":      "\uf410",
	".exe":     "\uf17a",
	".sh":      "\uf489",
	".bat":     "\uf17a",
	".ps1":     "\uf489",
	".jpg":     "\uf1c5",
	".png":     "\uf1c5",
	".gif":     "\uf1c5",
	".svg":     "\uf1c5",
	".mp4":     "\uf1c8",
	".mp3":     "\uf1c7",
}

const (
	defaultFileIcon = "\uf15b"
	symlinkIcon     = "\uf481"
)

// NewTheme builds the active theme from runtime options.
func NewTheme(colorMode string, noColor bool, charset string, timeMode string, useIcons bool) Theme {
	br := Branches{
		Middle: "├── ",
		Last:   "└── ",
		Pipe:   "│   ",
		Space:  "    ",
	}
	if charset == "ascii" {
		br = Branches{
			Middle: "|-- ",
			Last:   "\\-- ",
			Pipe:   "|   ",
			Space:  "    ",
		}
	}

	useColor := false
	switch colorMode {
	case "always":
		useColor = true
	case "never":
		useColor = false
	default:
		useColor = !noColor && IsTerminal(os.Stdout.Fd())
	}

	p := Palette{
		Reset:  "\033[0m",
		Blue:   "\033[34m",
		Cyan:   "\033[36m",
		Green:  "\033[32m",
		Yellow: "\033[33m",
	}
	if !useColor {
		p = Palette{}
	}

	tf := "2006-01-02 15:04"
	switch timeMode {
	case "mod", "modified", "change", "create", "":
		tf = "2006-01-02 15:04"
	}

	return Theme{
		Colors:     p,
		Branches:   br,
		TimeFormat: tf,
		UseIcons:   useIcons,
	}
}

func IsTerminal(fd uintptr) bool {
	if term.IsTerminal(int(fd)) {
		return true
	}
	if os.Getenv("WT_SESSION") != "" || os.Getenv("ConEmuANSI") == "ON" || os.Getenv("ANSICON") != "" {
		return true
	}
	return false
}

func IsHidden(name, fullPath string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}
	return platform.IsHiddenWindows(fullPath)
}

func IconFor(name string, isDir bool, isSymlink bool) string {
	if isSymlink {
		return symlinkIcon + " "
	}
	if isDir {
		return iconMap["dir"] + " "
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ic, ok := iconMap[ext]; ok {
		return ic + " "
	}
	return defaultFileIcon + " "
}

func DirIcon() string  { return iconMap["dir"] }
func FileIcon(ext string) string {
	if ic, ok := iconMap[strings.ToLower(ext)]; ok {
		return ic
	}
	return defaultFileIcon
}
