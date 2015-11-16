// Package uilive provides a writer that updates the UI
package uilive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
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

	buf       bytes.Buffer
	mtx       sync.Mutex
	lineCount int
}

// New returns a new writer with defaults
func New() *Writer {
	return &Writer{
		Out:             Out,
		RefreshInterval: RefreshInterval,
	}
}

// Flush writes to the out and resets the buffer. It should be called after the last call to Write to ensure that any data buffered in the Writer is written to output.
// Any incomplete escape sequence at the end is considered complete for formatting purposes.
func (w *Writer) Flush() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if len(w.buf.Bytes()) == 0 {
		return nil
	}
	for i := 0; i < w.lineCount; i++ {
		fmt.Fprintf(w.Out, "%c[%dA", ESC, 0) // move the cursor up
		fmt.Fprintf(w.Out, "%c[2K\r", ESC)   // clear the line
	}

	lines := 0
	for _, b := range w.buf.Bytes() {
		if b == '\n' {
			lines++
		}
	}
	w.lineCount = lines
	_, err := w.Out.Write(w.buf.Bytes())
	w.buf.Reset()
	return err
}

// Start starts the listener that updates the UI
func (w *Writer) Start() {
	go func() {
		for {
			w.Wait()
			w.Flush()
		}
	}()
}

// Wait waits for the writer to finish writing
func (w *Writer) Wait() {
	time.Sleep(w.RefreshInterval)
}

// Write writes buf to the writer b. The only errors returned are ones encountered while writing to the underlying output stream.
func (w *Writer) Write(buf []byte) (n int, err error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	return w.buf.Write(buf)
}
