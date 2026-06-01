package main

import (
	"testing"
	"time"
)

func TestGCounter_Commutativity(t *testing.T) {
	a := NewGCounter()
	b := NewGCounter()
	a.Increment("n1")
	a.Increment("n1")
	b.Increment("n2")

	ab := a.Merge(b)
	ba := b.Merge(a)
	if ab.Value() != ba.Value() {
		t.Fatalf("not commutative: %d != %d", ab.Value(), ba.Value())
	}
}

func TestGCounter_Idempotency(t *testing.T) {
	a := NewGCounter()
	a.Increment("n1")
	ab := a.Merge(a)
	if ab.Value() != a.Value() {
		t.Fatal("not idempotent")
	}
}

func TestPNCounter(t *testing.T) {
	a := NewPNCounter()
	b := NewPNCounter()
	a.Increment("n1")
	a.Increment("n1")
	a.Decrement("n1")
	b.Increment("n2")

	merged := a.Merge(b)
	if merged.Value() != 2 {
		t.Fatalf("want 2, got %d", merged.Value())
	}
}

func TestLWWRegister(t *testing.T) {
	var r1, r2 LWWRegister
	t1 := time.Now()
	t2 := t1.Add(time.Millisecond)
	r1.Set("old", "n1", t1)
	r2.Set("new", "n2", t2)
	merged := r1.Merge(r2)
	if merged.Value != "new" {
		t.Fatalf("LWW want new, got %q", merged.Value)
	}
}

func TestTwoPhaseSet(t *testing.T) {
	s := NewTwoPhaseSet()
	s.Add("alice")
	s.Add("bob")
	s.Remove("alice")
	if s.Contains("alice") {
		t.Fatal("alice should be removed")
	}
	if !s.Contains("bob") {
		t.Fatal("bob should be present")
	}
}
