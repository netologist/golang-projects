package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrOpen is returned when the breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

type state int

const (
	closed state = iota
	open
	halfOpen
)

// Breaker is a minimal circuit breaker.
type Breaker struct {
	mu              sync.Mutex
	state           state
	failures        int
	threshold       int
	halfOpenAt      time.Time
	recoveryTimeout time.Duration
}

// New creates a breaker that trips after `threshold` failures.
func New(threshold int, recoveryTimeout time.Duration) *Breaker {
	return &Breaker{threshold: threshold, recoveryTimeout: recoveryTimeout}
}

// Execute runs fn unless the breaker is open.
func (b *Breaker) Execute(fn func() error) error {
	b.mu.Lock()
	if b.state == open {
		if time.Now().Before(b.halfOpenAt) {
			b.mu.Unlock()
			return ErrOpen
		}
		b.state = halfOpen
	}
	b.mu.Unlock()

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.failures++
		if b.state == halfOpen || b.failures >= b.threshold {
			b.state = open
			b.halfOpenAt = time.Now().Add(b.recoveryTimeout)
		}
		return fmt.Errorf("breaker: %w", err)
	}
	b.failures = 0
	b.state = closed
	return nil
}
