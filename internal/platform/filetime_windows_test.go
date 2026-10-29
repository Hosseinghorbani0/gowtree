//go:build windows

package platform

import (
	"os"
	"testing"
)

func TestFileTimeUsesWindowsMetadata(t *testing.T) {
	path := t.TempDir() + "\\file.txt"
	if err := os.WriteFile(path, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"mod", "change", "create", ""} {
		if got := FileTime(info, mode); got.IsZero() {
			t.Fatalf("FileTime(%q) = zero time", mode)
		}
	}
}
