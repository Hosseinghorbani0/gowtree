package app_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binaryPath := filepath.Join(os.TempDir(), "gowtree-test.exe")
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/gowtree")
	cmd.Dir = filepath.Join("..", "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping integration test: cannot build binary: %v\n%s", err, out)
	}
	defer os.Remove(binaryPath)

	testDir, err := os.MkdirTemp("", "gowtree-int")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)
	os.Mkdir(filepath.Join(testDir, "sub"), 0755)
	os.WriteFile(filepath.Join(testDir, "readme.md"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(testDir, "sub", "app.exe"), []byte{}, 0755)

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"basic tree", []string{"-a", testDir}, []string{"sub", "readme.md", "app.exe", "directories", "files"}},
		{"JSON output", []string{"-J", testDir}, []string{`"type": "directory"`, `"name": "sub"`}},
		{"Markdown output", []string{"--markdown", testDir}, []string{"# ", "📂", "📄"}},
		{"HTML output", []string{"--html", testDir}, []string{"<h1>", "</ul>", "<li>"}},
		{"Icons", []string{"--icons", "-a", testDir}, []string{"\uf17a", "\uf07c"}},
		{"Sort by size", []string{"--sort", "size", "-a", testDir}, []string{"readme.md"}},
		{"No credit", []string{"--no-credit", testDir}, []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(binaryPath, tt.args...).Output()
			if err != nil {
				t.Fatalf("running gowtree: %v", err)
			}
			strOut := string(out)
			for _, w := range tt.want {
				if w == "" {
					if strings.Contains(strOut, "crafted by") {
						t.Error("expected no credit line")
					}
				} else if !strings.Contains(strOut, w) {
					t.Errorf("output missing %q", w)
				}
			}
		})
	}
}
