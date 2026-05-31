package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// BoundedGroup wraps errgroup.Group with a semaphore to bound parallelism.
type BoundedGroup struct {
	g   *errgroup.Group
	ctx context.Context
	sem *semaphore.Weighted
}

// NewBounded creates a group running at most maxConcurrent tasks at once.
func NewBounded(ctx context.Context, maxConcurrent int64) (*BoundedGroup, context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	return &BoundedGroup{g: g, ctx: gCtx, sem: semaphore.NewWeighted(maxConcurrent)}, gCtx
}

// Go schedules fn, blocking if the concurrency limit is reached.
func (bg *BoundedGroup) Go(fn func() error) {
	if err := bg.sem.Acquire(bg.ctx, 1); err != nil {
		bg.g.Go(func() error { return fmt.Errorf("acquire semaphore: %w", err) })
		return
	}
	bg.g.Go(func() error {
		defer bg.sem.Release(1)
		return fn()
	})
}

// Wait blocks until all tasks finish and returns the first error.
func (bg *BoundedGroup) Wait() error { return bg.g.Wait() }
