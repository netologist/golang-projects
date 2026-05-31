package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	l := NewLimiter(10) // 10 total units of capacity

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		weight := int64(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Do(context.Background(), weight, func() error {
				fmt.Printf("task using %d units\n", weight)
				return nil
			})
		}()
	}
	wg.Wait()
}
