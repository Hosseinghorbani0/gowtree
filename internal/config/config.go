package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// File holds optional defaults loaded from ~/.gowtree.yaml.
type File struct {
	Color       *string   `yaml:"color,omitempty"`
	NoColor     *bool     `yaml:"no_color,omitempty"`
	Charset     *string   `yaml:"charset,omitempty"`
	ShowAll     *bool     `yaml:"all,omitempty"`
	MaxDepth    *int      `yaml:"max_depth,omitempty"`
	DirsOnly    *bool     `yaml:"dirs_only,omitempty"`
	FilesOnly   *bool     `yaml:"files_only,omitempty"`
	ShowSize    *bool     `yaml:"show_size,omitempty"`
	Ignore      *[]string `yaml:"ignore,omitempty"`
	SortBy      *string   `yaml:"sort,omitempty"`
	SortRev     *bool     `yaml:"reverse,omitempty"`
	TimeFlag    *string   `yaml:"time,omitempty"`
	FilesFirst  *bool     `yaml:"files_first,omitempty"`
	Regex       *string   `yaml:"regex,omitempty"`
	Icons       *bool     `yaml:"icons,omitempty"`
	FollowLinks *bool     `yaml:"follow_symlinks,omitempty"`
	DU          *bool     `yaml:"du,omitempty"`
}

// DefaultPath returns the default config file location.
func DefaultPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir, _ = os.UserHomeDir()
	}
	return filepath.Join(dir, ".gowtree.yaml")
}

// Load reads YAML config from path. Missing files return nil, nil.
func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg File
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
