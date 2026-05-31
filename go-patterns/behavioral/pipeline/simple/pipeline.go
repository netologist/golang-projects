package main

import "context"

func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func filter(ctx context.Context, in <-chan int, pred func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if pred(n) {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
