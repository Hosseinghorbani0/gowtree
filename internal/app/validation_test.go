package app

import "testing"

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name string
		opts runtimeOptions
		want bool
	}{
		{
			name: "valid defaults",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name"},
			want: true,
		},
		{
			name: "invalid color",
			opts: runtimeOptions{colorMode: "invalid", charset: "utf8", sortBy: "name"},
		},
		{
			name: "invalid charset",
			opts: runtimeOptions{colorMode: "auto", charset: "invalid", sortBy: "name"},
		},
		{
			name: "invalid sort",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "invalid"},
		},
		{
			name: "invalid time",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name", timeMode: "invalid"},
		},
		{
			name: "negative depth",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name", maxDepth: -1},
		},
		{
			name: "directory and file modes",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name", dirsOnly: true, filesOnly: true},
		},
		{
			name: "multiple formats",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name", jsonOut: true, htmlOut: true},
		},
		{
			name: "invalid ignore pattern",
			opts: runtimeOptions{colorMode: "auto", charset: "utf8", sortBy: "name", ignore: []string{"["}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(tt.opts)
			if (err == nil) != tt.want {
				t.Fatalf("validateOptions() error = %v, want valid = %v", err, tt.want)
			}
		})
	}
}
