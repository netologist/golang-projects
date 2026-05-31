package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_succeedsOnThirdAttempt(t *testing.T) {
	attempt := 0
	err := Do(context.Background(), Config{
		MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1, Jitter: 0,
	}, func() error {
		attempt++
		if attempt < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempt != 3 {
		t.Errorf("attempt: got %d, want 3", attempt)
	}
}

func TestDo_nonRetryableStops(t *testing.T) {
	fatal := errors.New("fatal")
	called := 0
	err := Do(context.Background(), Config{
		MaxAttempts: 5, InitialDelay: time.Millisecond, Multiplier: 1, Jitter: 0,
		Retryable: func(e error) bool { return !errors.Is(e, fatal) },
	}, func() error {
		called++
		return fatal
	})
	if !errors.Is(err, fatal) {
		t.Errorf("expected fatal, got %v", err)
	}
	if called != 1 {
		t.Errorf("called %d times, want 1", called)
	}
}

func TestDo_respectsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Do(ctx, DefaultConfig, func() error { return errors.New("x") })
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
