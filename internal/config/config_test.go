package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg != nil {
			t.Fatalf("Load() config = %#v, want nil", cfg)
		}
	})

	t.Run("malformed yaml", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bad.yaml")
		if err := os.WriteFile(path, []byte("max_depth: ["), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(path); err == nil {
			t.Fatal("Load() error = nil, want malformed YAML error")
		}
	})

	t.Run("valid config", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "valid.yaml")
		if err := os.WriteFile(path, []byte("color: never\nmax_depth: 2\n"), 0600); err != nil {
			t.Fatal(err)
		}
		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if cfg.Color == nil || *cfg.Color != "never" {
			t.Fatalf("Color = %v, want never", cfg.Color)
		}
		if cfg.MaxDepth == nil || *cfg.MaxDepth != 2 {
			t.Fatalf("MaxDepth = %v, want 2", cfg.MaxDepth)
		}
	})
}
