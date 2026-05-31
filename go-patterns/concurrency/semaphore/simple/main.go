package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	sem := New(2) // at most 2 concurrent
	var concurrent, maxConcurrent atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sem.Acquire(context.Background()); err != nil {
				return
			}
			defer sem.Release()

			cur := concurrent.Add(1)
			for {
				m := maxConcurrent.Load()
				if cur <= m || maxConcurrent.CompareAndSwap(m, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			concurrent.Add(-1)
		}()
	}
	wg.Wait()
	fmt.Printf("max concurrent observed: %d (limit 2)\n", maxConcurrent.Load())
}
