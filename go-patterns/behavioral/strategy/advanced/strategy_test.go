package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSelectStrategy_bySize(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"tiny", 10},
		{"medium", 1000},
		{"large", 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := SelectStrategy(tt.size)
			if c == nil {
				t.Fatal("nil strategy")
			}
		})
	}
}

func TestGzipCompress_roundtrip(t *testing.T) {
	original := []byte(strings.Repeat("distributed systems ", 500))
	compressed, err := GzipCompress(original)
	if err != nil {
		t.Fatal(err)
	}
	if len(compressed) >= len(original) {
		t.Errorf("expected compression, got %d >= %d", len(compressed), len(original))
	}
	back, err := gunzip(compressed)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(back, original) {
		t.Error("round-trip mismatch")
	}
}

func TestNoCompress_tiny(t *testing.T) {
	data := []byte("hi")
	c := SelectStrategy(len(data))
	out, _ := c(data)
	if !bytes.Equal(out, data) {
		t.Error("tiny payload should be unchanged")
	}
}
