//go:build !windows
// +build !windows

package Helpers

import (
	"fmt"
	"github.com/creack/pty"
	"image"
	"os"
	"os/exec"
	"syscall"
)

var OperatingSystem = "linux"

func IsBotRunning(pid int) bool {
	if !authenticated {
		fmt.Println("Bot tracking is deactivated until a valid license code is present.")
		return false
	}

	// Linux-specific implementation
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}

func KillBot(botId int, userId int) {
	pid := GetPid(botId)

	if pid > 0 {
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = process.Signal(syscall.SIGKILL)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = UpdateBotStatus(botId, 0, "")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}

func SetWindowTitle(title string) error {
	cmd := exec.Command("/bin/sh")
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "\033]0;%s\007", title)
	if err != nil {
		return err
	}
	return nil
}

func CaptureWindowScreenshot(pid uint32) (*image.RGBA, error) {
	return nil, fmt.Errorf("not currently supported on Linux")
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}
