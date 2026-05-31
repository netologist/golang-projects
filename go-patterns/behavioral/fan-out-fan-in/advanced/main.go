package main

import (
	"context"
	"fmt"
)

func main() {
	in := make(chan int, 8)
	for i := 1; i <= 8; i++ {
		in <- i
	}
	close(in)

	out := FanOut(context.Background(), in, 4, func(_ context.Context, n int) int { return n * n })
	merged := FanIn(context.Background(), out)

	count := 0
	for range merged {
		count++
	}
	fmt.Printf("processed %d items across 4 workers\n", count)
}
