package GetTerminalSize

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	kernel32                       = windows.NewLazySystemDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type coord struct {
	X, Y int16
}

type smallRect struct {
	Left, Top, Right, Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

func GetTerminalWidth() (int, error) {
	stdoutHandle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return 0, err
	}

	var info consoleScreenBufferInfo
	ret, _, err := procGetConsoleScreenBufferInfo.Call(uintptr(stdoutHandle), uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return 0, err
	}

	width := int(info.Window.Right - info.Window.Left + 1)
	return width, nil
}
