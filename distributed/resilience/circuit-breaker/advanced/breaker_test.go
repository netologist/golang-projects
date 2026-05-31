package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func newBreaker(t *testing.T) *Breaker {
	t.Helper()
	return New(Config{FailureThreshold: 3, RecoveryTimeout: 50 * time.Millisecond, HalfOpenProbes: 1},
		prometheus.NewRegistry())
}

var errBoom = errors.New("boom")

func TestBreaker_tripsAfterThreshold(t *testing.T) {
	b := newBreaker(t)
	for i := 0; i < 3; i++ {
		_ = b.Execute(context.Background(), func() error { return errBoom })
	}
	err := b.Execute(context.Background(), func() error { return nil })
	if !errors.Is(err, ErrOpen) {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_recoversAfterCooldown(t *testing.T) {
	b := newBreaker(t)
	for i := 0; i < 3; i++ {
		_ = b.Execute(context.Background(), func() error { return errBoom })
	}
	time.Sleep(60 * time.Millisecond)
	if err := b.Execute(context.Background(), func() error { return nil }); err != nil {
		t.Fatalf("expected recovery, got: %v", err)
	}
	if b.State() != StateClosed {
		t.Errorf("expected Closed, got %s", b.State())
	}
}
