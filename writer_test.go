package uilive

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	w := New()
	b := &bytes.Buffer{}
	w.Out = b
	w.Start()
	for i := 0; i < 2; i++ {
		_, _ = fmt.Fprintln(w, "foo")
	}
	w.Stop()
	_, _ = fmt.Fprintln(b, "bar")

	want := "foo\nfoo\nbar\n"
	if b.String() != want {
		t.Fatalf("want %q, got %q", want, b.String())
	}
}

func TestStartCalledTwice(t *testing.T) {
	w := New()
	b := &bytes.Buffer{}
	w.Out = b

	w.Start()
	w.Stop()
	w.Start()
	w.Stop()
}

func TestWriterErase(t *testing.T) {
	w := New()
	b := &bytes.Buffer{}
	w.Out = b

	_, _ = fmt.Fprintln(w, "line 1")
	_, _ = fmt.Fprintln(w, "line 2")
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	if w.lineCount != 2 {
		t.Fatalf("want lineCount 2, got %d", w.lineCount)
	}

	_, _ = fmt.Fprintln(w, "buffered line")
	if w.buf.Len() == 0 {
		t.Fatal("expected buffer to contain data before erase")
	}

	w.Erase()

	if w.lineCount != 0 {
		t.Errorf("want lineCount 0, got %d", w.lineCount)
	}
	if w.buf.Len() != 0 {
		t.Errorf("want empty buf, got %d bytes", w.buf.Len())
	}

	wantClear := strings.Repeat(clear, 2)
	if !strings.HasSuffix(b.String(), wantClear) {
		t.Errorf("expected out to end with clear sequence %q, got %q", wantClear, b.String())
	}
}

