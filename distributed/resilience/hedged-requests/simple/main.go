package main

import (
	"context"
	"fmt"
	"time"
)

func call(ctx context.Context, id int, latency time.Duration) (int, error) {
	select {
	case <-time.After(latency):
		return id, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan int, 2)
	// Primary is slow; hedge after 30ms with a fast backup.
	go func() { v, _ := call(ctx, 1, 100*time.Millisecond); results <- v }()
	go func() {
		time.Sleep(30 * time.Millisecond)
		v, _ := call(ctx, 2, 10*time.Millisecond)
		results <- v
	}()

	winner := <-results
	cancel()
	fmt.Printf("winner: replica %d\n", winner)
}
