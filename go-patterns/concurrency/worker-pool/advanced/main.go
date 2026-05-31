package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	p := New(ctx, 4, 16)

	for i := 1; i <= 12; i++ {
		_ = p.Submit(ctx, Job{ID: i, Payload: i * i})
	}
	p.Shutdown()

	count := 0
	for r := range p.Results() {
		count++
		_ = r
	}
	fmt.Printf("completed %d jobs\n", count)
}
