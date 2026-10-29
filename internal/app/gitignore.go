package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func loadGitignore(root string) ([]string, error) {
	return loadIgnoreFile(filepath.Join(root, ".gitignore"))
}

func loadIgnoreFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			return nil, fmt.Errorf("negated patterns are not supported: %q", line)
		}
		if line != "" {
			patterns = append(patterns, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}
