package main

import (
	"context"
	"sync"
)

// FanOut distributes items from in to n goroutines, each running fn.
func FanOut(ctx context.Context, in <-chan int, n int, fn func(int) int) []<-chan int {
	outs := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		out := make(chan int)
		outs[i] = out
		go func(out chan int) {
			defer close(out)
			for v := range in {
				select {
				case out <- fn(v):
				case <-ctx.Done():
					return
				}
			}
		}(out)
	}
	return outs
}

// FanIn merges multiple channels into one.
func FanIn(ctx context.Context, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(ch <-chan int) {
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
