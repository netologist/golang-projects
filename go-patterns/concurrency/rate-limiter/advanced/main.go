package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	l := New(20, 5) // 20/sec, burst 5
	ctx := context.Background()

	start := time.Now()
	for i := 1; i <= 10; i++ {
		_ = l.Wait(ctx)
		fmt.Printf("request %d at %v\n", i, time.Since(start).Round(time.Millisecond))
	}
}
