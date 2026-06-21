//go:build windows

package platform

import "syscall"

func IsHiddenWindows(path string) bool {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
