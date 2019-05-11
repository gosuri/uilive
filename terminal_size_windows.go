// +build windows

package uilive

import (
	"math"
	"syscall"
	"unsafe"
)

type consoleFontInfo struct {
	font     uint32
	fontSize coord
}

const (
	SmCxMin = 28
	SmCyMin = 29
)

var (
	tmpConsoleFontInfo        consoleFontInfo
	moduleUser32              = syscall.NewLazyDLL("user32.dll")
	procGetCurrentConsoleFont = kernel32.NewProc("GetCurrentConsoleFont")
	getSystemMetrics          = moduleUser32.NewProc("GetSystemMetrics")
)

func getCurrentConsoleFont(h syscall.Handle, info *consoleFontInfo) (err error) {
	r0, _, e1 := syscall.Syscall(
		procGetCurrentConsoleFont.Addr(), 3, uintptr(h), 0, uintptr(unsafe.Pointer(info)),
	)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getTermSize() (int, int) {
	out, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		return 0, 0
	}

	x, _, err := getSystemMetrics.Call(SmCxMin)
	y, _, err := getSystemMetrics.Call(SmCyMin)

	if x == 0 || y == 0 {
		if err != nil {
			panic(err)
		}
	}

	err = getCurrentConsoleFont(out, &tmpConsoleFontInfo)
	if err != nil {
		panic(err)
	}

	return int(math.Ceil(float64(x) / float64(tmpConsoleFontInfo.fontSize.x))), int(math.Ceil(float64(y) / float64(tmpConsoleFontInfo.fontSize.y)))
}
