//go:build windows
// +build windows

package Helpers

import (
	"fmt"
	"image"
	"os/exec"
	"strconv"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	kernel32             = syscall.MustLoadDLL("kernel32.dll")
	procSetConsoleTitleW = kernel32.MustFindProc("SetConsoleTitleW")
	OperatingSystem      = "windows"
)

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

func SetWindowTitle(title string) error {
	lpString := utf16.Encode([]rune(title))
	_, _, err := syscall.Syscall(
		procSetConsoleTitleW.Addr(),
		2,
		uintptr(unsafe.Pointer(&lpString[0])),
		uintptr(1),
		uintptr(0))
	if err != 0 {
		return err
	}
	return nil
}

func CaptureWindowScreenshot(pid uint32) (*image.RGBA, error) {
	return nil, fmt.Errorf("not currently supported")
}
