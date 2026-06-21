package output

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CopyToClipboard(text string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "clip")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	default:
		for _, prog := range []string{"xclip", "xsel"} {
			args := []string{}
			if prog == "xclip" {
				args = []string{"-selection", "clipboard"}
			} else if prog == "xsel" {
				args = []string{"-ib"}
			}
			if _, err := exec.LookPath(prog); err == nil {
				cmd := exec.Command(prog, args...)
				cmd.Stdin = strings.NewReader(text)
				return cmd.Run()
			}
		}
		return fmt.Errorf("no clipboard tool found (install xclip or xsel)")
	}
}
