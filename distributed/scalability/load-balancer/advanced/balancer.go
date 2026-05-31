package main

import (
	"errors"
	"sync/atomic"
)

// ErrNoHealthyBackend is returned when no backend can serve a request.
var ErrNoHealthyBackend = errors.New("no healthy backend available")

// Backend is a target with health and connection tracking.
type Backend struct {
	Addr        string
	Weight      int
	ActiveConns atomic.Int64
	healthy     atomic.Bool
}

// NewBackend creates a healthy backend.
func NewBackend(addr string) *Backend {
	b := &Backend{Addr: addr, Weight: 1}
	b.healthy.Store(true)
	return b
}

// SetHealthy updates the backend's health.
func (b *Backend) SetHealthy(v bool) { b.healthy.Store(v) }

// Healthy reports the backend's health.
func (b *Backend) Healthy() bool { return b.healthy.Load() }

// Balancer selects the next backend.
type Balancer interface {
	Next(backends []*Backend) (*Backend, error)
}

// RoundRobin rotates through healthy backends.
type RoundRobin struct{ counter atomic.Uint64 }

// Next returns the next healthy backend in rotation.
func (rr *RoundRobin) Next(backends []*Backend) (*Backend, error) {
	healthy := make([]*Backend, 0, len(backends))
	for _, b := range backends {
		if b.Healthy() {
			healthy = append(healthy, b)
		}
	}
	if len(healthy) == 0 {
		return nil, ErrNoHealthyBackend
	}
	idx := rr.counter.Add(1) - 1
	return healthy[idx%uint64(len(healthy))], nil
}

// LeastConns picks the healthy backend with the fewest active connections.
type LeastConns struct{}

// Next returns the least-loaded healthy backend.
func (lc *LeastConns) Next(backends []*Backend) (*Backend, error) {
	var best *Backend
	for _, b := range backends {
		if !b.Healthy() {
			continue
		}
		if best == nil || b.ActiveConns.Load() < best.ActiveConns.Load() {
			best = b
		}
	}
	if best == nil {
		return nil, ErrNoHealthyBackend
	}
	return best, nil
}
