package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBus_deliverToSubscribers(t *testing.T) {
	bus := NewBus()
	var mu sync.Mutex
	var received []string
	cancel := bus.Subscribe(context.Background(), "greet", func(e Event) {
		mu.Lock()
		received = append(received, e.Payload)
		mu.Unlock()
	}, 10)
	defer cancel()

	bus.Publish(context.Background(), Event{Topic: "greet", Payload: "hello"})
	bus.Publish(context.Background(), Event{Topic: "greet", Payload: "world"})

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Errorf("got %d events, want 2", len(received))
	}
}

func TestBus_cancelUnsubscribes(t *testing.T) {
	bus := NewBus()
	var mu sync.Mutex
	count := 0
	cancel := bus.Subscribe(context.Background(), "ping", func(_ Event) {
		mu.Lock()
		count++
		mu.Unlock()
	}, 10)
	cancel()

	bus.Publish(context.Background(), Event{Topic: "ping"})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if count > 0 {
		t.Error("received event after unsubscribe")
	}
}
