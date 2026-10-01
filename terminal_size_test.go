//go:build !windows
// +build !windows

package uilive

import (
	"testing"
	"unsafe"
)

func TestWindowSizeMemoryLayout(t *testing.T) {
	var ws windowSize

	// POSIX/Darwin struct winsize has 4 uint16 fields: ws_row, ws_col, ws_xpixel, ws_ypixel.
	// Ensure windowSize matches the 8-byte layout expected by TIOCGWINSZ ioctl.
	const expectedSize = 8
	if actualSize := unsafe.Sizeof(ws); actualSize != expectedSize {
		t.Fatalf("expected windowSize size %d bytes, got %d", expectedSize, actualSize)
	}

	if offset := unsafe.Offsetof(ws.rows); offset != 0 {
		t.Errorf("expected offset of rows to be 0, got %d", offset)
	}
	if offset := unsafe.Offsetof(ws.cols); offset != 2 {
		t.Errorf("expected offset of cols to be 2, got %d", offset)
	}
	if offset := unsafe.Offsetof(ws.xpixel); offset != 4 {
		t.Errorf("expected offset of xpixel to be 4, got %d", offset)
	}
	if offset := unsafe.Offsetof(ws.ypixel); offset != 6 {
		t.Errorf("expected offset of ypixel to be 6, got %d", offset)
	}
}

func TestGetTermSize(t *testing.T) {
	w, h := getTermSize()
	if w < 0 || h < 0 {
		t.Fatalf("expected non-negative terminal dimensions, got width=%d, height=%d", w, h)
	}
}
