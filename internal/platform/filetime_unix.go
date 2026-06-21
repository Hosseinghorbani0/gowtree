//go:build !windows

package platform

import (
	"os"
	"time"
)

func FileTime(info os.FileInfo, timeMode string) time.Time {
	return info.ModTime()
}
