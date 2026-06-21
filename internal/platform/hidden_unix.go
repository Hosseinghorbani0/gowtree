//go:build !windows

package platform

func IsHiddenWindows(path string) bool { return false }
