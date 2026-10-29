package tree

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/hosseinghorbani0/gowtree/internal/style"
)

type metadataErrorEntry struct {
	name string
}

func (e metadataErrorEntry) Name() string      { return e.name }
func (e metadataErrorEntry) IsDir() bool       { return false }
func (e metadataErrorEntry) Type() os.FileMode { return 0 }
func (e metadataErrorEntry) Info() (os.FileInfo, error) {
	return nil, errors.New("metadata unavailable")
}

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
	if n := CountItems(dir, opts); n != 2 {
		t.Errorf("count without -a = %d, want 2", n)
	}
	opts.ShowAll = true
	if n := CountItems(dir, opts); n != 3 {
		t.Errorf("count with -a = %d, want 3", n)
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
	if n := CountItems(dir, opts); n != 1 {
		t.Errorf("count with ignore = %d, want 1", n)
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

func TestBuildJSONDoesNotFollowSymlinkWithoutFollowOption(t *testing.T) {
	dir, err := os.MkdirTemp("", "gowtree-symlink")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	link := filepath.Join(dir, "parent")
	if err := os.Symlink(dir, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	opts := defaultOpts()
	opts.MaxDepth = 1
	node := BuildJSON(dir, 0, nil, opts)

	if len(node.Children) != 1 {
		t.Fatalf("got %d children, want 1", len(node.Children))
	}
	if got := node.Children[0].Type; got != "symlink" {
		t.Fatalf("symlink type = %q, want %q", got, "symlink")
	}
}

func TestBuildJSONDetectsSymlinkCycleWhenFollowingLinks(t *testing.T) {
	dir, err := os.MkdirTemp("", "gowtree-symlink-cycle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	link := filepath.Join(dir, "parent")
	if err := os.Symlink(dir, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	opts := defaultOpts()
	opts.FollowSymlinks = true
	node := BuildJSON(dir, 0, nil, opts)

	if len(node.Children) != 1 {
		t.Fatalf("got %d children, want 1", len(node.Children))
	}
	if got := node.Children[0].Type; got != "symlink-loop" {
		t.Fatalf("symlink cycle type = %q, want %q", got, "symlink-loop")
	}
}

func TestBuildJSONReportsBrokenSymlinkWhenFollowingLinks(t *testing.T) {
	dir, err := os.MkdirTemp("", "gowtree-broken-symlink")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	link := filepath.Join(dir, "missing")
	if err := os.Symlink(filepath.Join(dir, "does-not-exist"), link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	opts := defaultOpts()
	opts.FollowSymlinks = true
	node := BuildJSON(dir, 0, nil, opts)

	if len(node.Children) != 1 {
		t.Fatalf("got %d children, want 1", len(node.Children))
	}
	if got := node.Children[0].Type; got != "broken-symlink" {
		t.Fatalf("broken symlink type = %q, want %q", got, "broken-symlink")
	}
}

func BenchmarkDirectorySizeRepeatedScans(b *testing.B) {
	root := b.TempDir()
	for dir := 0; dir < 8; dir++ {
		sub := filepath.Join(root, "dir"+strconv.Itoa(dir))
		if err := os.Mkdir(sub, 0755); err != nil {
			b.Fatal(err)
		}

		for file := 0; file < 32; file++ {
			path := filepath.Join(sub, "file"+strconv.Itoa(file))
			if err := os.WriteFile(path, make([]byte, 256), 0644); err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("uncached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = DirSize(root, false)
		}
	})
	b.Run("cached", func(b *testing.B) {
		cache := make(map[string]int64)
		for i := 0; i < b.N; i++ {
			_ = dirSize(root, cache)
		}
	})
}

func TestDirSizeWithErrorReportsMissingPath(t *testing.T) {
	_, err := DirSizeWithError(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("DirSizeWithError() error = nil, want missing path error")
	}
}

func TestSortEntriesStrictReportsMetadataError(t *testing.T) {
	_, err := SortEntriesStrict([]os.DirEntry{metadataErrorEntry{name: "gone"}}, defaultOpts())
	if err == nil {
		t.Fatal("SortEntriesStrict() error = nil, want metadata error")
	}
}
