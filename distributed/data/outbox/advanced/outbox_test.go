package main

import (
	"context"
	"testing"
)

func TestOutbox_atLeastOnce(t *testing.T) {
	s := NewStore()
	ctx := context.Background()

	_, err := s.SaveWithOutbox(ctx, "order:1", []byte("record"), "order", []byte("OrderCreated"))
	if err != nil {
		t.Fatal(err)
	}

	var published []string
	n, err := s.Relay(ctx, 10, func(e *OutboxEvent) error {
		published = append(published, e.ID)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("delivered %d, want 1", n)
	}

	// Second relay delivers nothing (already published).
	n2, _ := s.Relay(ctx, 10, func(_ *OutboxEvent) error { return nil })
	if n2 != 0 {
		t.Errorf("second relay delivered %d, want 0", n2)
	}
}

func TestOutbox_pendingCount(t *testing.T) {
	s := NewStore()
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		_, _ = s.SaveWithOutbox(ctx, "k", []byte("r"), "agg", []byte("p"))
	}
	pending, _ := s.PendingEvents(ctx, 10)
	if len(pending) != 3 {
		t.Errorf("got %d pending, want 3", len(pending))
	}
}
