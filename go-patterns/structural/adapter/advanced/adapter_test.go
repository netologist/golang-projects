package main

import (
	"context"
	"fmt"
	"testing"
)

type fakeSyncReader struct{ data map[string]string }

func (f *fakeSyncReader) ReadSync(key string) (string, error) {
	v, ok := f.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

func TestAsyncAdapter_happyPath(t *testing.T) {
	r := NewAsyncAdapter(&fakeSyncReader{data: map[string]string{"x": "42"}})
	res := <-r.Read(context.Background(), "x")
	if res.Err != nil {
		t.Fatal(res.Err)
	}
	if res.Value != "42" {
		t.Errorf("got %s, want 42", res.Value)
	}
}

func TestAsyncAdapter_missingKey(t *testing.T) {
	r := NewAsyncAdapter(&fakeSyncReader{data: map[string]string{}})
	res := <-r.Read(context.Background(), "x")
	if res.Err == nil {
		t.Error("expected error for missing key")
	}
}
