// Package uilive provides a writer that updates the UI
package uilive

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

// ESC is the ASCII code for escape character
const ESC = 27

// RefreshInterval is the default refresh interval to update the ui
var RefreshInterval = time.Millisecond

// Out is the default out for the writer
var Out = os.Stdout

// ErrClosedPipe is the error returned when trying to writer is not listening
var ErrClosedPipe = errors.New("uilive: read/write on closed pipe")

// Writer represent the writer that updates the UI
type Writer struct {
	// Out is the writer to write to
	Out io.Writer

	// RefreshInterval is the time the UI sould refresh
	RefreshInterval time.Duration

	// stopChan is buffered channel for stopping the listener
	stopChan chan struct{}
	// running is flag for determining if the listerner is running
	running bool

	buf       bytes.Buffer
	mtx       sync.Mutex
	lineCount int
}

// New returns a new writer with defaults
func New() *Writer {
	return &Writer{
		Out:             Out,
		RefreshInterval: RefreshInterval,

		stopChan: make(chan struct{}, 1),
	}
}

// Flush writes to the out and resets the buffer. It should be called after the last call to Write to ensure that any data buffered in the Writer is written to output.
// Any incomplete escape sequence at the end is considered complete for formatting purposes.
func (w *Writer) Flush() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	// do nothing is  buffer is empty
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
	w.Flush()
	w.stopChan <- struct{}{}
}

// Listen listens for updates to the writers buffer and flushes to the out. It blocks the runtime.
func (w *Writer) Listen() {
	if w.running {
		return
	}
	go func() {
		w.running = true
		for {
			w.Wait()
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
