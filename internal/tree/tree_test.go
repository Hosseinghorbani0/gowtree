package tree

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hosseinghorbani0/gowtree/internal/style"
)

func defaultOpts() Options {
	return Options{
		ShowAll: true,
		SortBy:  "name",
		Theme:   style.NewTheme("never", true, "utf8", "", false),
	}
}

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"app.exe", true},
		{"script.sh", true},
		{"run.bat", true},
		{"hello.cmd", true},
		{"module.ps1", true},
		{"README.md", false},
	}
	for _, tt := range tests {
		if got := IsExecutable(tt.name); got != tt.want {
			t.Errorf("IsExecutable(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCountItems(t *testing.T) {
	dir, err := os.MkdirTemp("", "gowtree-count")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	os.Mkdir(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0644)
	os.WriteFile(filepath.Join(dir, ".hidden"), nil, 0644)

	opts := defaultOpts()
	opts.ShowAll = false
	if n := CountItems(dir, opts); n != 3 {
		t.Errorf("count without -a = %d, want 3", n)
	}
	opts.ShowAll = true
	if n := CountItems(dir, opts); n != 4 {
		t.Errorf("count with -a = %d, want 4", n)
	}
}

func TestIgnorePattern(t *testing.T) {
	dir, err := os.MkdirTemp("", "gowtree-ignore")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0644)
	os.WriteFile(filepath.Join(dir, "b.log"), nil, 0644)

	opts := defaultOpts()
	opts.Ignore = []string{"*.log"}
	if n := CountItems(dir, opts); n != 2 {
		t.Errorf("count with ignore = %d, want 2", n)
	}
}

func TestRegexFilter(t *testing.T) {
	re := regexp.MustCompile(`\.go$`)
	opts := defaultOpts()
	opts.Regex = re
	dir, err := os.MkdirTemp("", "gowtree-regex")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	os.WriteFile(filepath.Join(dir, "main.go"), nil, 0644)
	os.WriteFile(filepath.Join(dir, "readme.md"), nil, 0644)
	if n := CountItems(dir, opts); n != 1 {
		t.Errorf("regex filter count = %d, want 1", n)
	}
}
