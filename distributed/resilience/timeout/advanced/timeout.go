package main

import (
	"context"
	"fmt"
	"time"
)

// Do runs fn with a deadline of d. On timeout it returns a wrapped
// context.DeadlineExceeded.
func Do[T any](ctx context.Context, d time.Duration, fn func(context.Context) (T, error)) (T, error) {
	tCtx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	type result struct {
		val T
		err error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := fn(tCtx)
		ch <- result{v, err}
	}()

	select {
	case r := <-ch:
		return r.val, r.err
	case <-tCtx.Done():
		var zero T
		return zero, fmt.Errorf("timeout after %s: %w", d, tCtx.Err())
	}
}
