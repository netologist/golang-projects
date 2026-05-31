package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	b := New(map[string]Config{
		"critical": {MaxConcurrent: 3},
		"batch":    {MaxConcurrent: 1},
	})

	var accepted, rejected atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := b.Execute(context.Background(), "batch", func() error {
				time.Sleep(15 * time.Millisecond)
				return nil
			})
			if err != nil {
				rejected.Add(1)
			} else {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("batch partition: accepted=%d rejected=%d\n", accepted.Load(), rejected.Load())
}
