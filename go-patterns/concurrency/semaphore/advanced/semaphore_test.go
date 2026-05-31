package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter_boundsConcurrency(t *testing.T) {
	l := NewLimiter(3)
	var inFlight, maxSeen atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Do(context.Background(), 1, func() error {
				cur := inFlight.Add(1)
				for {
					m := maxSeen.Load()
					if cur <= m || maxSeen.CompareAndSwap(m, cur) {
						break
					}
				}
				time.Sleep(5 * time.Millisecond)
				inFlight.Add(-1)
				return nil
			})
		}()
	}
	wg.Wait()
	if maxSeen.Load() > 3 {
		t.Errorf("max concurrency %d exceeded limit 3", maxSeen.Load())
	}
}

func TestLimiter_weightTooLarge(t *testing.T) {
	l := NewLimiter(2)
	err := l.Do(context.Background(), 5, func() error { return nil })
	if err == nil {
		t.Error("expected error for oversized weight")
	}
}
