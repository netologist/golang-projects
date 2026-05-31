package main

import (
	"context"
	"fmt"
	"sync"
)

// Message wraps a payload with its topic.
type Message[T any] struct {
	Topic   string
	Payload T
}

type subscriber[T any] struct {
	ch  chan Message[T]
	ctx context.Context
}

// Topic is a typed pub/sub topic with async delivery.
type Topic[T any] struct {
	name string
	mu   sync.RWMutex
	subs []*subscriber[T]
}

// NewTopic creates a topic.
func NewTopic[T any](name string) *Topic[T] { return &Topic[T]{name: name} }

// Subscribe returns a receive channel and an unsubscribe func.
func (t *Topic[T]) Subscribe(ctx context.Context, bufSize int) (<-chan Message[T], func()) {
	subCtx, cancel := context.WithCancel(ctx)
	s := &subscriber[T]{ch: make(chan Message[T], bufSize), ctx: subCtx}

	t.mu.Lock()
	t.subs = append(t.subs, s)
	t.mu.Unlock()

	unsub := func() {
		cancel()
		t.mu.Lock()
		defer t.mu.Unlock()
		for i, sub := range t.subs {
			if sub == s {
				t.subs = append(t.subs[:i], t.subs[i+1:]...)
				break
			}
		}
	}
	return s.ch, unsub
}

// Publish delivers payload to all subscribers; full buffers are dropped.
func (t *Topic[T]) Publish(ctx context.Context, payload T) error {
	msg := Message[T]{Topic: t.name, Payload: payload}
	t.mu.RLock()
	subs := t.subs
	t.mu.RUnlock()

	for _, s := range subs {
		select {
		case s.ch <- msg:
		case <-s.ctx.Done():
		case <-ctx.Done():
			return fmt.Errorf("publish: %w", ctx.Err())
		default:
			fmt.Printf("warn: topic %s subscriber buffer full\n", t.name)
		}
	}
	return nil
}
