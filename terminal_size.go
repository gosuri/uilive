package uilive

import (
	"github.com/nsf/termbox-go"
)

func getTermSize() (int, int) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()
	return w, h
}
