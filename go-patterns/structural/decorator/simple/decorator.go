package main

import (
	"fmt"
	"io"
)

// CountingWriter wraps an io.Writer and counts bytes written.
type CountingWriter struct {
	inner io.Writer
	count int
}

// NewCountingWriter wraps w.
func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{inner: w}
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.inner.Write(p)
	c.count += n
	return n, err
}

// BytesWritten reports the total bytes written so far.
func (c *CountingWriter) BytesWritten() int { return c.count }

// TimingWriter wraps an io.Writer and logs write calls.
type TimingWriter struct {
	inner io.Writer
	label string
}

// NewTimingWriter wraps w with a labelled log line per write.
func NewTimingWriter(w io.Writer, label string) *TimingWriter {
	return &TimingWriter{inner: w, label: label}
}

func (t *TimingWriter) Write(p []byte) (int, error) {
	fmt.Printf("[%s] writing %d bytes\n", t.label, len(p))
	return t.inner.Write(p)
}
