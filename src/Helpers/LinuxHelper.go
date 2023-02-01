//go:build !windows
// +build !windows

package Helpers

import (
	"fmt"
	"os"
	"syscall"
)

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
