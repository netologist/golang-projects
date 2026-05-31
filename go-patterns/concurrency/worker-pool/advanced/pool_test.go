package main

import (
	"context"
	"testing"
)

func TestPool_processesAllJobs(t *testing.T) {
	ctx := context.Background()
	p := New(ctx, 3, 10)

	const n = 10
	for i := 0; i < n; i++ {
		if err := p.Submit(ctx, Job{ID: i, Payload: i}); err != nil {
			t.Fatal(err)
		}
	}
	p.Shutdown()

	var count int
	for range p.Results() {
		count++
	}
	if count != n {
		t.Errorf("got %d results, want %d", count, n)
	}
}

func TestPool_cancelledSubmit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := New(ctx, 1, 0)
	cancel()

	err := p.Submit(ctx, Job{ID: 1})
	if err == nil {
		t.Error("expected error on cancelled submit")
	}
}
