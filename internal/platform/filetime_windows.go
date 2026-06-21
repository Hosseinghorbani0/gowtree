//go:build windows

package platform

import (
	"os"
	"syscall"
	"time"
)

func FileTime(info os.FileInfo, timeMode string) time.Time {
	d := info.Sys().(*syscall.Win32FileAttributeData)
	switch timeMode {
	case "create":
		return time.Unix(0, d.CreationTime.Nanoseconds())
	case "change":
		return time.Unix(0, d.LastWriteTime.Nanoseconds())
	default:
		return time.Unix(0, d.LastWriteTime.Nanoseconds())
	}
}
