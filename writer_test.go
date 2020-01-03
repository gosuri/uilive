package uilive

import (
	"bytes"
	"fmt"
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
