// Package uilive provides a writer that updates the UI
package uilive

import (
	"bytes"
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
	stopChan  chan struct{}
	running   bool
	lineCount int
}

// New returns a new writer with defaults
func New() *Writer {
	return &Writer{
		Out:             Out,
		RefreshInterval: RefreshInterval,
		stopChan:        make(chan struct{}, 1),
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
	w.clearLines()

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

// Start starts the listener in a non blocking manner
func (w *Writer) Start() {
	go w.Listen()
}

// Stop stops the listener that updates the UI
func (w *Writer) Stop() {
	w.stopChan <- struct{}{}
	return
}

// Listen listens for updates to the writers buffer and flushes to the out. It blocks the runtime.
func (w *Writer) Listen() {
	if w.running {
		return
	}
	go func() {
		w.running = true
		for {
			time.Sleep(w.RefreshInterval)
			w.Flush()
		}
	}()
	<-w.stopChan
	w.running = false
}

// Wait waits for the writer to finish writing
func (w *Writer) Wait() {
	time.Sleep(w.RefreshInterval)
	w.Flush()
}

// Write writes buf to the writer b. The only errors returned are ones encountered while writing to the underlying output stream.
func (w *Writer) Write(buf []byte) (n int, err error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	return w.buf.Write(buf)
}
