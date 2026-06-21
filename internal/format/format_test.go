package format

import "testing"

func TestSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}
	for _, tt := range tests {
		if got := Size(tt.bytes); got != tt.want {
			t.Errorf("Size(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}
