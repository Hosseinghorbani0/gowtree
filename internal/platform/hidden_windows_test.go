//go:build windows

package platform

import "testing"

func TestIsHiddenWindowsReturnsFalseForMissingPath(t *testing.T) {
	if IsHiddenWindows(`C:\path\that\does\not\exist\gowtree`) {
		t.Fatal("missing path was classified as hidden")
	}
}
