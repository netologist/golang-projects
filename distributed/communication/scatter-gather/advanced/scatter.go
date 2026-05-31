package main

import (
	"context"
	"time"
)

// Result carries one backend's outcome.
type Result[T any] struct {
	Source  string
	Value   T
	Err     error
	Latency time.Duration
}

// Gather fans out to all fns and returns the results that arrive before the
// context deadline (partial results on timeout).
func Gather[T any](ctx context.Context, fns map[string]func(context.Context) (T, error)) []Result[T] {
	ch := make(chan Result[T], len(fns))

	for name, fn := range fns {
		go func(name string, fn func(context.Context) (T, error)) {
			start := time.Now()
			v, err := fn(ctx)
			ch <- Result[T]{Source: name, Value: v, Err: err, Latency: time.Since(start)}
		}(name, fn)
	}

	results := make([]Result[T], 0, len(fns))
	for i := 0; i < len(fns); i++ {
		select {
		case r := <-ch:
			results = append(results, r)
		case <-ctx.Done():
			return results
		}
	}
	return results
}
