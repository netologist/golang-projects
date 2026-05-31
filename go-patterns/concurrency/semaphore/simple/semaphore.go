package main

import "context"

// Semaphore is a counting semaphore backed by a buffered channel.
type Semaphore chan struct{}

// New creates a semaphore allowing n concurrent holders.
func New(n int) Semaphore { return make(chan struct{}, n) }

// Acquire blocks until a slot is free or the context is cancelled.
func (s Semaphore) Acquire(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees one slot.
func (s Semaphore) Release() { <-s }

// TryAcquire grabs a slot without blocking; returns false if none free.
func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}
