// +build windows

package uilive

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetStdHandle               = kernel32.NewProc("GetStdHandle")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procGetConsoleCursorInfo       = kernel32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")
	procSetConsoleTextAttribute    = kernel32.NewProc("SetConsoleTextAttribute")
	procScrollConsoleScreenBuffer  = kernel32.NewProc("ScrollConsoleScreenBufferW")
)

type wchar uint16
type short int16
type dword uint32
type word uint16

type coord struct {
	x short
	y short
}

type smallRect struct {
	left   short
	top    short
	right  short
	bottom short
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}

type consoleCursorInfo struct {
	size    dword
	visible int32
}

func (w *Writer) clearLines() {
	f, ok := w.Out.(*os.File)
	if ok && !isatty.IsTerminal(f.Fd()) {
		ok = false
	}
	if !ok {
		for i := 0; i < w.lineCount; i++ {
			fmt.Fprintf(w.Out, "%c[%dA", ESC, 0) // move the cursor up
			fmt.Fprintf(w.Out, "%c[2K\r", ESC)   // clear the line
		}
		return
	}
	fd := f.Fd()
	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&csbi)))

	for i := 0; i < w.lineCount; i++ {
		// move the cursor up
		csbi.cursorPosition.y--
		procSetConsoleCursorPosition.Call(fd, uintptr(*(*int32)(unsafe.Pointer(&csbi.cursorPosition))))
		// clear the line
		cursor := coord{
			x: csbi.window.left,
			y: csbi.window.top + csbi.cursorPosition.y,
		}
		var count, w dword
		count = dword(csbi.size.x)
		procFillConsoleOutputCharacter.Call(fd, uintptr(' '), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&w)))
		procFillConsoleOutputAttribute.Call(fd, uintptr(csbi.attributes), uintptr(count), *(*uintptr)(unsafe.Pointer(&cursor)), uintptr(unsafe.Pointer(&w)))

	}
}
