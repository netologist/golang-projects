package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

// Limiter wraps a weighted semaphore to bound total resource units in flight.
type Limiter struct {
	sem   *semaphore.Weighted
	total int64
}

// NewLimiter allows up to `total` weighted units concurrently.
func NewLimiter(total int64) *Limiter {
	return &Limiter{sem: semaphore.NewWeighted(total), total: total}
}

// Do acquires `weight` units, runs fn, and releases them.
func (l *Limiter) Do(ctx context.Context, weight int64, fn func() error) error {
	if weight > l.total {
		return fmt.Errorf("weight %d exceeds capacity %d", weight, l.total)
	}
	if err := l.sem.Acquire(ctx, weight); err != nil {
		return fmt.Errorf("acquire %d units: %w", weight, err)
	}
	defer l.sem.Release(weight)
	return fn()
}
