package main

import (
	"testing"
	"time"
)

func TestWrite_ReplicatesToReplicas(t *testing.T) {
	primary := NewPrimary("p1")
	r1 := NewReplica("r1", 20*time.Millisecond)
	r2 := NewReplica("r2", 20*time.Millisecond)

	if err := primary.Write("k", "v", []*Node{r1, r2}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	e, ok := r1.Read("k")
	if !ok || e.Value != "v" {
		t.Fatalf("r1 not replicated: %v %v", e, ok)
	}
	e, ok = r2.Read("k")
	if !ok || e.Value != "v" {
		t.Fatalf("r2 not replicated: %v %v", e, ok)
	}
}

func TestWrite_VersionIncrement(t *testing.T) {
	p := NewPrimary("p")
	_ = p.Write("k", "v1", nil)
	_ = p.Write("k", "v2", nil)
	e, _ := p.Read("k")
	if e.Version != 1 {
		t.Fatalf("want version 1, got %d", e.Version)
	}
}
