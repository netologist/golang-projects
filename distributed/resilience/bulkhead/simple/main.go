package main

import (
	"fmt"
	"sync"
	"time"
)

// partition is a counting semaphore.
type partition chan struct{}

func (p partition) tryRun(fn func()) bool {
	select {
	case p <- struct{}{}:
		defer func() { <-p }()
		fn()
		return true
	default:
		return false
	}
}

func main() {
	critical := make(partition, 2)
	var wg sync.WaitGroup
	accepted, rejected := 0, 0
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := critical.tryRun(func() { time.Sleep(20 * time.Millisecond) })
			mu.Lock()
			if ok {
				accepted++
			} else {
				rejected++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("accepted=%d rejected=%d (capacity 2)\n", accepted, rejected)
}
