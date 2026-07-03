//go:build windows

package collector

import (
	"agent/internal/model"
	"syscall"
	"unsafe"
)

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	getForeground = user32.NewProc("GetForegroundWindow")
	getWindowText = user32.NewProc("GetWindowTextW")
)

func GetFocusedWindowTitle() (string, error) {
	hwnd, _, _ := getForeground.Call()

	var title string
	if hwnd == 0 {
		title = model.EmptyFocusedWindow
		return title, nil
	}

	b := make([]uint16, 256)
	_, _, _ = getWindowText.Call(hwnd, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))

	title = syscall.UTF16ToString(b)
	if title == "" {
		title = model.EmptyFocusedWindow
	}
	return title, nil
}
