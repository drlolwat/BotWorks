//go:build !windows
// +build !windows

package Helpers

import (
	"fmt"
	"os"
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
