package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 9; i++ {
			in <- i
		}
	}()

	workers := FanOut(ctx, in, 3, func(n int) int { return n * n })
	merged := FanIn(ctx, workers...)

	sum := 0
	for v := range merged {
		sum += v
	}
	fmt.Println("sum of squares 1..9 =", sum)
}
