package style

import (
	"strings"
	"testing"
)

func TestNewThemeASCIIWithoutColor(t *testing.T) {
	theme := NewTheme("never", false, "ascii", "mod", false)
	if theme.Colors != (Palette{}) {
		t.Fatalf("colors = %#v, want empty palette", theme.Colors)
	}
	if theme.Branches.Middle != "|-- " || theme.Branches.Last != "\\-- " ||
		theme.Branches.Pipe != "|   " || theme.Branches.Space != "    " {
		t.Fatalf("unexpected ASCII branches: %#v", theme.Branches)
	}
	if theme.TimeFormat == "" {
		t.Fatal("time format is empty")
	}
}

func TestNewThemeAlwaysUsesColorsAndIcons(t *testing.T) {
	theme := NewTheme("always", false, "utf8", "create", true)
	if theme.Colors.Blue == "" || theme.Colors.Reset == "" {
		t.Fatalf("colors = %#v, want ANSI palette", theme.Colors)
	}
	if !theme.UseIcons {
		t.Fatal("UseIcons = false, want true")
	}
	if theme.Branches.Middle != "├── " || theme.Branches.Last != "└── " {
		t.Fatalf("unexpected UTF-8 branches: %#v", theme.Branches)
	}
}

func TestIconForAndFileIcon(t *testing.T) {
	if got := IconFor("link", false, true); !strings.Contains(got, symlinkIcon) {
		t.Fatalf("symlink icon = %q, want %q", got, symlinkIcon)
	}
	if got := IconFor("src", true, false); !strings.Contains(got, iconMap["dir"]) {
		t.Fatalf("directory icon = %q, want %q", got, iconMap["dir"])
	}
	if got := FileIcon(".GO"); got != iconMap[".go"] {
		t.Fatalf("FileIcon(.GO) = %q, want %q", got, iconMap[".go"])
	}
	if got := FileIcon(".unknown"); got != defaultFileIcon {
		t.Fatalf("FileIcon(.unknown) = %q, want default %q", got, defaultFileIcon)
	}
}

func TestIsHiddenDotfile(t *testing.T) {
	if !IsHidden(".env", "C:\\temporary\\.env") {
		t.Fatal("dotfile was not classified as hidden")
	}
	if IsHidden("visible", "C:\\temporary\\visible") {
		t.Fatal("visible file was classified as hidden")
	}
}
