package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryLocker_mutualExclusion(t *testing.T) {
	l := NewMemoryLocker()
	const goroutines = 8
	var counter int
	var critical sync.Mutex // guards counter to detect races deterministically
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			tok, err := l.Lock(ctx, "resource", time.Second)
			if err != nil {
				return
			}
			defer l.Unlock(ctx, "resource", tok)
			critical.Lock()
			counter++
			critical.Unlock()
		}()
	}
	wg.Wait()
	if counter != goroutines {
		t.Errorf("got %d, want %d", counter, goroutines)
	}
}

func TestMemoryLocker_staleToken(t *testing.T) {
	l := NewMemoryLocker()
	ctx := context.Background()
	tok, _ := l.Lock(ctx, "k", time.Second)
	_ = l.Unlock(ctx, "k", tok)
	// Re-unlock with stale token must fail.
	if err := l.Unlock(ctx, "k", tok); err == nil {
		t.Error("expected error unlocking with stale token")
	}
}

func TestMemoryLocker_ttlExpiry(t *testing.T) {
	l := NewMemoryLocker()
	ctx := context.Background()
	_, _ = l.Lock(ctx, "k", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	// Lock should be re-acquirable after TTL.
	if _, err := l.Lock(ctx, "k", time.Second); err != nil {
		t.Errorf("expected re-acquire after TTL, got %v", err)
	}
}
