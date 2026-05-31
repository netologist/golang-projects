package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_completesInTime(t *testing.T) {
	val, err := Do(context.Background(), time.Second, func(_ context.Context) (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}
}

func TestDo_timesOut(t *testing.T) {
	_, err := Do(context.Background(), 10*time.Millisecond, func(ctx context.Context) (int, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
