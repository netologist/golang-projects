package main

import (
	"context"
	"fmt"
	"sync/atomic"
)

func main() {
	g, _ := NewBounded(context.Background(), 4)
	var done atomic.Int64
	for i := 0; i < 20; i++ {
		g.Go(func() error {
			done.Add(1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("completed %d tasks with max 4 concurrent\n", done.Load())
}
