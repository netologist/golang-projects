package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestHedge_returnsResult(t *testing.T) {
	var calls atomic.Int32
	val, err := Do(context.Background(), 5*time.Millisecond, 3, func(_ context.Context) (int, error) {
		calls.Add(1)
		time.Sleep(8 * time.Millisecond)
		return 42, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}
}

func TestHedge_fastFirstWinsBeforeHedge(t *testing.T) {
	var calls atomic.Int32
	_, err := Do(context.Background(), 50*time.Millisecond, 3, func(_ context.Context) (int, error) {
		calls.Add(1)
		return 1, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Errorf("expected 1 call (fast first), got %d", calls.Load())
	}
}
