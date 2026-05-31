package main

import (
	"context"
	"testing"
	"time"
)

func TestTopic_delivery(t *testing.T) {
	topic := NewTopic[string]("greet")
	ch, unsub := topic.Subscribe(context.Background(), 10)
	defer unsub()

	_ = topic.Publish(context.Background(), "hello")
	_ = topic.Publish(context.Background(), "world")

	time.Sleep(10 * time.Millisecond)
	if len(ch) != 2 {
		t.Errorf("got %d messages, want 2", len(ch))
	}
}

func TestTopic_unsubscribe(t *testing.T) {
	topic := NewTopic[int]("nums")
	_, unsub := topic.Subscribe(context.Background(), 4)
	unsub()

	// After unsubscribe there are no subscribers.
	if err := topic.Publish(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	topic.mu.RLock()
	n := len(topic.subs)
	topic.mu.RUnlock()
	if n != 0 {
		t.Errorf("expected 0 subscribers, got %d", n)
	}
}
