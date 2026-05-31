package main

import (
	"context"
	"sync"
)

// Stage processes one item, possibly returning an error.
type Stage[T, U any] func(ctx context.Context, in T) (U, error)

// ErrItem carries a processing error with its source value.
type ErrItem[T any] struct {
	Input T
	Err   error
}

// RunStage runs fn on every item from in using up to workers goroutines.
// Results go to out; errors to errs. Both channels close when input drains.
func RunStage[T, U any](
	ctx context.Context,
	in <-chan T,
	workers int,
	fn Stage[T, U],
) (<-chan U, <-chan ErrItem[T]) {
	out := make(chan U, workers)
	errs := make(chan ErrItem[T], workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range in {
				result, err := fn(ctx, item)
				if err != nil {
					select {
					case errs <- ErrItem[T]{Input: item, Err: err}:
					case <-ctx.Done():
						return
					}
					continue
				}
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()
	return out, errs
}
