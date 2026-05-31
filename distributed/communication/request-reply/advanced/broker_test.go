package main

import (
	"context"
	"testing"
	"time"
)

func TestBroker_requestReply(t *testing.T) {
	b := NewBroker()

	reply, err := b.Request(context.Background(), time.Second, func(id string) {
		// Simulate an async responder on another goroutine.
		go func() {
			time.Sleep(5 * time.Millisecond)
			b.HandleReply(id, []byte("pong"))
		}()
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(reply) != "pong" {
		t.Errorf("got %s, want pong", reply)
	}
}

func TestBroker_timeout(t *testing.T) {
	b := NewBroker()
	_, err := b.Request(context.Background(), 10*time.Millisecond, func(_ string) {
		// Never reply.
	})
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestBroker_unknownCorrelation(t *testing.T) {
	b := NewBroker()
	if b.HandleReply("nonexistent", []byte("x")) {
		t.Error("expected false for unknown correlation id")
	}
}
