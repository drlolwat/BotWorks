//go:build windows
// +build windows

package Helpers

import (
	"fmt"
	"syscall"
)

var OperatingSystem = "windows"

func IsBotRunning(pid int) bool {
	if !authenticated {
		fmt.Println("Bot tracking is deactivated until a valid license code is present.")
		return false
	}
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(handle)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false
	}

	return exitCode == 259
}

func KillBot(botId int, userId int) {
	pid := GetPid(botId)

	if pid > 0 {
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
		_ = cmd.Start()

		err := UpdateBotStatus(botId, 0, "")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}
