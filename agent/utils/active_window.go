//go:build windows

package utils

import (
	"syscall"
	"unsafe"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	getForeground = user32.NewProc("GetForegroundWindow")
	getWindowText = user32.NewProc("GetWindowTextW")
)

func GetActiveWindowTitle() string {
	hwnd, _, _ := getForeground.Call()
	if hwnd == 0 {
		return ""
	}

	b := make([]uint16, 256)
	getWindowText.Call(hwnd, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	return syscall.UTF16ToString(b)
}
