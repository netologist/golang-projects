package main

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestBulkhead_isolatesPartitions(t *testing.T) {
	b := New(map[string]Config{
		"critical": {MaxConcurrent: 2},
		"batch":    {MaxConcurrent: 1},
	})

	var wg sync.WaitGroup
	block := make(chan struct{})
	started := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = b.Execute(context.Background(), "batch", func() error {
			close(started)
			<-block
			return nil
		})
	}()
	<-started // batch partition now saturated

	// Critical partition is independent and should accept.
	err := b.Execute(context.Background(), "critical", func() error { return nil })
	close(block)
	wg.Wait()

	if err != nil {
		t.Errorf("critical partition blocked by batch: %v", err)
	}
}

func TestBulkhead_rejectsWhenFull(t *testing.T) {
	b := New(map[string]Config{"p": {MaxConcurrent: 1}})
	block := make(chan struct{})
	started := make(chan struct{})
	go func() {
		_ = b.Execute(context.Background(), "p", func() error {
			close(started)
			<-block
			return nil
		})
	}()
	<-started

	err := b.Execute(context.Background(), "p", func() error { return nil })
	close(block)
	if !errors.Is(err, ErrRejected) {
		t.Errorf("expected ErrRejected, got %v", err)
	}
}
