package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type pending struct{ ch chan []byte }

// Broker correlates asynchronous requests with their replies.
type Broker struct {
	waiters sync.Map // correlationID -> *pending
	counter atomic.Uint64
}

// NewBroker creates an empty broker.
func NewBroker() *Broker { return &Broker{} }

// Request registers a pending reply, invokes send with the new correlation ID,
// and waits for the reply or a timeout.
func (b *Broker) Request(ctx context.Context, timeout time.Duration, send func(correlationID string)) ([]byte, error) {
	id := fmt.Sprintf("corr-%d", b.counter.Add(1))
	p := &pending{ch: make(chan []byte, 1)}
	b.waiters.Store(id, p)
	defer b.waiters.Delete(id)

	send(id)

	select {
	case reply := <-p.ch:
		return reply, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("request %s: timed out after %s", id, timeout)
	case <-ctx.Done():
		return nil, fmt.Errorf("request %s: %w", id, ctx.Err())
	}
}

// HandleReply delivers a reply by correlation ID; returns false if no waiter.
func (b *Broker) HandleReply(correlationID string, reply []byte) bool {
	v, ok := b.waiters.Load(correlationID)
	if !ok {
		return false
	}
	v.(*pending).ch <- reply
	return true
}
