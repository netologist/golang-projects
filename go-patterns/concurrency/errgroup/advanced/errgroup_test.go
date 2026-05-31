package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestBoundedGroup_limitsParallelism(t *testing.T) {
	var maxConcurrent, current atomic.Int64

	g, _ := NewBounded(context.Background(), 3)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			cur := current.Add(1)
			for {
				m := maxConcurrent.Load()
				if cur <= m || maxConcurrent.CompareAndSwap(m, cur) {
					break
				}
			}
			current.Add(-1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	if maxConcurrent.Load() > 3 {
		t.Errorf("max concurrent was %d, want <= 3", maxConcurrent.Load())
	}
}

func TestBoundedGroup_firstErrorCancels(t *testing.T) {
	errBoom := errors.New("boom")
	g, ctx := NewBounded(context.Background(), 2)
	g.Go(func() error { return errBoom })
	g.Go(func() error { <-ctx.Done(); return ctx.Err() })
	if err := g.Wait(); !errors.Is(err, errBoom) {
		t.Errorf("expected errBoom, got %v", err)
	}
}
