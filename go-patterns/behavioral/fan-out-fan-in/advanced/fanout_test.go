package main

import (
	"context"
	"sort"
	"testing"
)

func TestFanOutFanIn_allItemsProcessed(t *testing.T) {
	in := make(chan int, 10)
	for i := 0; i < 10; i++ {
		in <- i
	}
	close(in)

	w1 := FanOut(context.Background(), in, 3, func(_ context.Context, n int) int { return n * 2 })
	merged := FanIn(context.Background(), w1)

	var results []int
	for v := range merged {
		results = append(results, v)
	}
	sort.Ints(results)
	if len(results) != 10 {
		t.Errorf("got %d items, want 10", len(results))
	}
}
