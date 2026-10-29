package app

import (
	"fmt"
	"path/filepath"
)

type runtimeOptions struct {
	colorMode string
	charset   string
	sortBy    string
	timeMode  string
	maxDepth  int
	dirsOnly  bool
	filesOnly bool
	jsonOut   bool
	markdown  bool
	htmlOut   bool
	ignore    []string
}

func validateOptions(opts runtimeOptions) error {
	switch opts.colorMode {
	case "auto", "always", "never":
	default:
		return fmt.Errorf("invalid color mode %q: want auto, always, or never", opts.colorMode)
	}
	switch opts.charset {
	case "utf8", "ascii":
	default:
		return fmt.Errorf("invalid charset %q: want utf8 or ascii", opts.charset)
	}
	switch opts.sortBy {
	case "name", "size", "time", "none":
	default:
		return fmt.Errorf("invalid sort mode %q: want name, size, time, or none", opts.sortBy)
	}
	switch opts.timeMode {
	case "", "mod", "change", "create":
	default:
		return fmt.Errorf("invalid time mode %q: want mod, change, or create", opts.timeMode)
	}
	if opts.maxDepth < 0 {
		return fmt.Errorf("invalid maximum depth %d: want a non-negative value", opts.maxDepth)
	}
	if opts.dirsOnly && opts.filesOnly {
		return fmt.Errorf("directories-only and files-only modes cannot be used together")
	}
	for _, pattern := range opts.ignore {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid ignore pattern %q: %w", pattern, err)
		}
	}
	formats := 0
	if opts.jsonOut {
		formats++
	}
	if opts.markdown {
		formats++
	}
	if opts.htmlOut {
		formats++
	}
	if formats > 1 {
		return fmt.Errorf("JSON, Markdown, and HTML output modes are mutually exclusive")
	}
	return nil
}
