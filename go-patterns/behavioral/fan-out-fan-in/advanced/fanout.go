package main

import (
	"context"
	"sync"
)

// FanOut distributes from in to n workers, each running fn. Results are unordered.
func FanOut[T, U any](ctx context.Context, in <-chan T, n int, fn func(context.Context, T) U) <-chan U {
	out := make(chan U, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range in {
				select {
				case out <- fn(ctx, item):
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// FanIn merges cs into a single channel.
func FanIn[T any](ctx context.Context, cs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
