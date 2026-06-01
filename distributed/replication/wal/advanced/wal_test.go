package main

import (
	"testing"
)

func TestWAL_AppendAndApply(t *testing.T) {
	var applied []Entry
	w := New(func(e Entry) { applied = append(applied, e) })

	w.Append("SET", "k", "v")
	w.Append("SET", "k2", "v2")

	if len(applied) != 2 {
		t.Fatalf("want 2 applied, got %d", len(applied))
	}
	if applied[0].Key != "k" || applied[0].Value != "v" {
		t.Errorf("wrong first entry: %+v", applied[0])
	}
}

func TestWAL_Recovery(t *testing.T) {
	var applied []Entry
	w := New(func(e Entry) { applied = append(applied, e) })

	w.Append("SET", "a", "1")
	w.Append("SET", "b", "2")
	w.Checkpoint(2) // mark first 2 as checkpointed
	w.Append("SET", "c", "3")
	w.Append("SET", "d", "4")

	// Simulate crash: reset applied
	applied = nil

	// Recover: only entries after checkpoint should replay
	if err := w.Recover(); err != nil {
		t.Fatal(err)
	}
	if len(applied) != 2 {
		t.Fatalf("want 2 replayed, got %d: %v", len(applied), applied)
	}
}

func TestWAL_Closed(t *testing.T) {
	w := New(func(Entry) {})
	w.Close()
	_, err := w.Append("SET", "k", "v")
	if err == nil {
		t.Fatal("expected error on closed WAL")
	}
}
