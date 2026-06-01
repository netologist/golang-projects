package main

import (
	"testing"
	"time"
)

func TestLWWResolver_RemoteWins(t *testing.T) {
	l1 := NewLeader("l1", LWWResolver)
	l2 := NewLeader("l2", LWWResolver)

	l1.Write("k", "old")
	time.Sleep(5 * time.Millisecond)
	l2.Write("k", "new")

	conflicts := l1.Sync(l2)
	if len(conflicts) != 1 {
		t.Fatalf("want 1 conflict, got %d", len(conflicts))
	}
	r, _ := l1.Read("k")
	if r.Value != "new" {
		t.Fatalf("LWW should pick newer value, got %q", r.Value)
	}
}

func TestSync_NoConflict(t *testing.T) {
	l1 := NewLeader("l1", LWWResolver)
	l2 := NewLeader("l2", LWWResolver)

	l1.Write("k1", "v1")
	l2.Write("k2", "v2")

	conflicts := l1.Sync(l2)
	if len(conflicts) != 0 {
		t.Fatalf("want 0 conflicts, got %d", len(conflicts))
	}
	_, ok := l1.Read("k2")
	if !ok {
		t.Fatal("k2 not synced")
	}
}
