package main

import (
	"context"
	"testing"
	"time"
)

func TestLimiter_burstAllowed(t *testing.T) {
	l := New(1, 5) // rate 1/s, burst 5
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Errorf("allow %d: expected true", i)
		}
	}
	if l.Allow() {
		t.Error("6th allow: expected false (burst exhausted)")
	}
}

func TestLimiter_Wait_cancellation(t *testing.T) {
	l := New(0.5, 0) // no initial tokens, slow refill
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := l.Wait(ctx); err == nil {
		t.Error("expected context error")
	}
}

func TestLimiter_refillsOverTime(t *testing.T) {
	l := New(100, 1) // fast refill
	if !l.Allow() {
		t.Fatal("first token should be available")
	}
	if l.Allow() {
		t.Fatal("burst is 1, second immediate allow should fail")
	}
	time.Sleep(30 * time.Millisecond) // refill ~3 tokens
	if !l.Allow() {
		t.Error("expected token after refill")
	}
}
