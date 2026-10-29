package metadata

import "testing"

func TestMetadataDefaults(t *testing.T) {
	if Version == "" {
		t.Fatal("Version is empty")
	}
	if Author == "" || RepoURL == "" {
		t.Fatalf("metadata constants are incomplete: author=%q repo=%q", Author, RepoURL)
	}
}

func TestMetadataVersionCanBeLinkerInjected(t *testing.T) {
	original := Version
	defer func() { Version = original }()
	Version = "v9.9.9"
	if Version != "v9.9.9" {
		t.Fatalf("Version = %q, want linker-compatible mutable value", Version)
	}
}
