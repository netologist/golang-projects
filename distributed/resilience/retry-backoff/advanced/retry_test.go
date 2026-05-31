package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_returnsValueOnSuccess(t *testing.T) {
	n := 0
	v, err := Do(context.Background(), Config{
		MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1,
	}, func(_ context.Context) (string, error) {
		n++
		if n < 2 {
			return "", errors.New("transient")
		}
		return "ok", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if v != "ok" {
		t.Errorf("got %q, want ok", v)
	}
}

func TestDo_nonRetryable(t *testing.T) {
	fatal := errors.New("permanent")
	calls := 0
	_, err := Do(context.Background(), Config{
		MaxAttempts: 5, InitialDelay: time.Millisecond, Multiplier: 1,
		Retryable: func(e error) bool { return !errors.Is(e, fatal) },
	}, func(_ context.Context) (int, error) {
		calls++
		return 0, fatal
	})
	if !errors.Is(err, fatal) {
		t.Errorf("expected fatal, got %v", err)
	}
	if calls != 1 {
		t.Errorf("calls: got %d, want 1", calls)
	}
}
