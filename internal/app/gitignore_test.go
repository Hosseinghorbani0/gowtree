package app

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadGitignore(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("# generated\n\n/bin/\n*.tmp\n"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := loadGitignore(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"bin", "*.tmp"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("patterns = %#v, want %#v", got, want)
	}
}

func TestLoadGitignoreMissingFile(t *testing.T) {
	got, err := loadGitignore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("patterns = %#v, want nil", got)
	}
}

func TestLoadGitignoreRejectsNegation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("!keep.txt\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := loadGitignore(root); err == nil {
		t.Fatal("loadGitignore() error = nil, want unsupported negation error")
	}
}
