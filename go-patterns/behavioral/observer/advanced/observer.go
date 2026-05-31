package main

import (
	"context"
	"fmt"
	"sync"
)

// Event is a topic/payload pair.
type Event struct{ Topic, Payload string }

// Handler reacts to an Event.
type Handler func(Event)

type subscription struct {
	ch     chan Event
	cancel context.CancelFunc
}

// Bus delivers events asynchronously, one goroutine per subscriber.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]*subscription
}

// NewBus creates an empty async bus.
func NewBus() *Bus { return &Bus{subs: map[string][]*subscription{}} }

// Subscribe registers h for topic and returns a cancel func.
func (b *Bus) Subscribe(ctx context.Context, topic string, h Handler, bufSize int) func() {
	subCtx, cancel := context.WithCancel(ctx)
	s := &subscription{ch: make(chan Event, bufSize), cancel: cancel}

	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], s)
	b.mu.Unlock()

	go func() {
		for {
			select {
			case e := <-s.ch:
				h(e)
			case <-subCtx.Done():
				return
			}
		}
	}()

	return func() {
		cancel()
		b.mu.Lock()
		defer b.mu.Unlock()
		subs := b.subs[topic]
		for i, sub := range subs {
			if sub == s {
				b.subs[topic] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}

// Publish delivers e to current subscribers; full buffers are dropped with a warning.
func (b *Bus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	subs := b.subs[e.Topic]
	b.mu.RUnlock()
	for _, s := range subs {
		select {
		case s.ch <- e:
		case <-ctx.Done():
			return
		default:
			fmt.Printf("warn: subscriber buffer full for topic %s\n", e.Topic)
		}
	}
}
