// Package uilive provides a writer that updates the UI
package uilive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	ESC = 27
)

var (
	// RefreshInterval is the default refresh interval to update the ui
	RefreshInterval = time.Millisecond
	// Out is the default out for the writer
	Out = os.Stdout
)

// Writer represent the writer that updates the UI
type Writer struct {
	// Out is the writer to write to
	Out io.Writer

	// RefreshInterval is the time the UI sould refresh
	RefreshInterval time.Duration

	writtenLines int
	buf          bytes.Buffer
	written      []byte
}

// New returns a new writer with defaults
func New() *Writer {
	return &Writer{
		Out:             Out,
		RefreshInterval: RefreshInterval,
	}
}

// Flush writes to the out and flushes the buffer
func (w *Writer) Flush() {
	cur := w.buf.Bytes()
	if len(cur) == 0 {
		return
	}

	lines := 0
	for _, b := range cur {
		if b == '\n' {
			lines++
		}
	}
	for i := 0; i < w.writtenLines; i++ {
		fmt.Fprintf(w.Out, "%c[%dA", ESC, i) // move the cursor up
		fmt.Fprintf(w.Out, "%c[2K\r", ESC)   // clear the line
	}

	w.writtenLines = lines
	w.Out.Write(cur)
	w.written = cur
	w.buf = bytes.Buffer{}
}

// Start starts the listener that updates the UI
func (w *Writer) Start() {
	go func() {
		for {
			w.Flush()
		}
	}()
}

// Wait waits for the writer to finish writing
func (w *Writer) Wait() {
	time.Sleep(w.RefreshInterval)
}

// Write writes the writers buffer
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}
