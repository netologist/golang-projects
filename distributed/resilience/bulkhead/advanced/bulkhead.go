package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrRejected is returned when a partition is at capacity.
var ErrRejected = errors.New("bulkhead: request rejected")

// Config sizes a partition.
type Config struct {
	MaxConcurrent int
}

// Partition bounds concurrency for one workload.
type Partition struct {
	sem chan struct{}
}

func newPartition(cfg Config) *Partition {
	return &Partition{sem: make(chan struct{}, cfg.MaxConcurrent)}
}

// Execute runs fn if a slot is free, else returns ErrRejected.
func (p *Partition) Execute(ctx context.Context, fn func() error) error {
	select {
	case p.sem <- struct{}{}:
		defer func() { <-p.sem }()
		return fn()
	case <-ctx.Done():
		return fmt.Errorf("bulkhead: %w", ctx.Err())
	default:
		return fmt.Errorf("execute: %w", ErrRejected)
	}
}

// Bulkhead is a set of named partitions.
type Bulkhead struct {
	mu         sync.RWMutex
	partitions map[string]*Partition
}

// New builds a bulkhead from a partition->config map.
func New(partitions map[string]Config) *Bulkhead {
	b := &Bulkhead{partitions: map[string]*Partition{}}
	for name, cfg := range partitions {
		b.partitions[name] = newPartition(cfg)
	}
	return b
}

// Execute runs fn in the named partition.
func (b *Bulkhead) Execute(ctx context.Context, partition string, fn func() error) error {
	b.mu.RLock()
	p, ok := b.partitions[partition]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("bulkhead: partition %q not found", partition)
	}
	return p.Execute(ctx, fn)
}
