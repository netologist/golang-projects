package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter implements a token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	maxBurst float64
	rate     float64 // tokens per second
	lastTick time.Time
}

// New creates a limiter that refills `rate` tokens/sec up to `burst` capacity.
func New(rate float64, burst int) *Limiter {
	return &Limiter{
		tokens:   float64(burst),
		maxBurst: float64(burst),
		rate:     rate,
		lastTick: time.Now(),
	}
}

func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens = min(l.maxBurst, l.tokens+elapsed*l.rate)
	l.lastTick = now
}

// Allow reports whether a token is available, consuming one if so.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Wait blocks until a token is available or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.Allow() {
			return nil
		}
		wait := time.Duration(float64(time.Second) / l.rate / 2)
		select {
		case <-ctx.Done():
			return fmt.Errorf("rate limiter: %w", ctx.Err())
		case <-time.After(wait):
		}
	}
}
