package main

import (
	"context"
	"fmt"
)

// SyncReader is a blocking legacy API.
type SyncReader interface {
	ReadSync(key string) (string, error)
}

// Result carries an async read outcome.
type Result struct {
	Value string
	Err   error
}

// AsyncReader is the modern interface our system uses.
type AsyncReader interface {
	Read(ctx context.Context, key string) <-chan Result
}

// AsyncAdapter wraps a SyncReader and exposes it as an AsyncReader.
type AsyncAdapter struct{ inner SyncReader }

// NewAsyncAdapter wraps r.
func NewAsyncAdapter(r SyncReader) *AsyncAdapter { return &AsyncAdapter{inner: r} }

// Read runs the blocking read in a goroutine and honours context cancellation.
func (a *AsyncAdapter) Read(ctx context.Context, key string) <-chan Result {
	ch := make(chan Result, 1)
	go func() {
		defer close(ch)
		val, err := a.inner.ReadSync(key)
		select {
		case ch <- Result{Value: val, Err: err}:
		case <-ctx.Done():
			ch <- Result{Err: fmt.Errorf("read %s: %w", key, ctx.Err())}
		}
	}()
	return ch
}
