package app_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hosseinghorbani0/gowtree/internal/app"
)

func captureRun(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	outRead, outWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errRead, errWrite, err := os.Pipe()
	if err != nil {
		_ = outRead.Close()
		_ = outWrite.Close()
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = outWrite, errWrite
	t.Cleanup(func() {
		os.Stdout, os.Stderr = oldOut, oldErr
		_ = outRead.Close()
		_ = outWrite.Close()
		_ = errRead.Close()
		_ = errWrite.Close()
	})

	code := app.Run(args)
	if err := outWrite.Close(); err != nil {
		t.Fatal(err)
	}
	if err := errWrite.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(outRead)
	if err != nil {
		t.Fatal(err)
	}
	errOutput, err := io.ReadAll(errRead)
	if err != nil {
		t.Fatal(err)
	}
	return code, string(out), string(errOutput)
}

func writeEmptyConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunCLIContracts(t *testing.T) {
	configPath := writeEmptyConfig(t)
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file.txt"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "version uses stdout",
			args:       []string{"--config", configPath, "--version"},
			wantCode:   0,
			wantStdout: "gowtree ",
		},
		{
			name:       "invalid sort uses stderr",
			args:       []string{"--config", configPath, "--sort", "invalid", root},
			wantCode:   2,
			wantStderr: "Invalid options:",
		},
		{
			name:       "missing path uses stderr",
			args:       []string{"--config", configPath, filepath.Join(root, "missing")},
			wantCode:   1,
			wantStderr: "Error opening path:",
		},
		{
			name:       "json is stdout",
			args:       []string{"--config", configPath, "--no-credit", "--J", root},
			wantCode:   0,
			wantStdout: `"type": "directory"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := captureRun(t, tt.args...)
			if code != tt.wantCode {
				t.Fatalf("Run() code = %d, want %d; stdout=%q stderr=%q", code, tt.wantCode, stdout, stderr)
			}
			if tt.wantStdout != "" && !strings.Contains(stdout, tt.wantStdout) {
				t.Errorf("stdout %q does not contain %q", stdout, tt.wantStdout)
			}
			if tt.name == "json is stdout" {
				var document map[string]any
				if err := json.Unmarshal([]byte(stdout), &document); err != nil {
					t.Fatalf("JSON stdout is invalid: %v; output=%q", err, stdout)
				}
				if document["type"] != "directory" {
					t.Fatalf("JSON root type = %v, want directory", document["type"])
				}
			}
			if tt.wantStderr != "" && !strings.Contains(stderr, tt.wantStderr) {
				t.Errorf("stderr %q does not contain %q", stderr, tt.wantStderr)
			}
			if tt.wantStdout != "" && stderr != "" {
				t.Errorf("unexpected stderr for successful output: %q", stderr)
			}
			if tt.wantStderr != "" && stdout != "" {
				t.Errorf("unexpected stdout for error: %q", stdout)
			}
		})
	}
}

func TestRunReportsOutputWriterFailure(t *testing.T) {
	configPath := writeEmptyConfig(t)
	oldOut, oldErr := os.Stdout, os.Stderr
	outRead, outWrite, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errRead, errWrite, err := os.Pipe()
	if err != nil {
		_ = outRead.Close()
		_ = outWrite.Close()
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = outWrite, errWrite
	t.Cleanup(func() {
		os.Stdout, os.Stderr = oldOut, oldErr
		_ = outRead.Close()
		_ = outWrite.Close()
		_ = errRead.Close()
		_ = errWrite.Close()
	})

	if err := outWrite.Close(); err != nil {
		t.Fatal(err)
	}
	code := app.Run([]string{"--config", configPath, t.TempDir()})
	if code == 0 {
		t.Fatal("Run() code = 0, want output writer failure")
	}
	_ = errWrite.Close()
	_, _ = io.ReadAll(outRead)
	_, _ = io.ReadAll(errRead)
}
