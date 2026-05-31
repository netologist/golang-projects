package main

import (
	"context"
	"time"
)

// Do issues up to n attempts of fn, staggered by delay. The first successful
// result wins and all other attempts are cancelled.
func Do[T any](ctx context.Context, delay time.Duration, n int, fn func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		val T
		err error
	}
	results := make(chan result, n)

	for i := 0; i < n; i++ {
		go func() {
			v, err := fn(ctx)
			results <- result{val: v, err: err}
		}()
		if i < n-1 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				var zero T
				return zero, ctx.Err()
			case r := <-results:
				if r.err == nil {
					return r.val, nil
				}
				// first attempt errored quickly; fall through to launch next
			}
		}
	}

	r := <-results
	return r.val, r.err
}
