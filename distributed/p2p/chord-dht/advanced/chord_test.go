package main

import "testing"

func TestChord_LookupFindsCorrectNode(t *testing.T) {
	r := NewRing()
	for _, id := range []int{0, 8, 16, 24, 32, 40, 48, 56} {
		r.Join(id)
	}

	// key 10 should be on node 16 (first node ≥ 10)
	node, hops := r.Lookup(10)
	if node != 16 {
		t.Fatalf("want node 16, got %d", node)
	}
	if hops <= 0 {
		t.Fatalf("hops should be > 0, got %d", hops)
	}
}

func TestChord_FingerTableBuilt(t *testing.T) {
	r := NewRing()
	for _, id := range []int{0, 8, 16, 24, 32, 40, 48, 56} {
		r.Join(id)
	}
	n := r.nodes[0]
	for i, f := range n.FingerTable {
		if f < 0 || f >= RingSize {
			t.Errorf("finger[%d]=%d out of range", i, f)
		}
	}
}

func TestChord_LogNHops(t *testing.T) {
	// With 8 evenly spaced nodes on a 64-bit ring, max hops should be ≤ M
	r := NewRing()
	for _, id := range []int{0, 8, 16, 24, 32, 40, 48, 56} {
		r.Join(id)
	}
	maxHops := 0
	for key := 0; key < RingSize; key++ {
		_, hops := r.Lookup(key)
		if hops > maxHops {
			maxHops = hops
		}
	}
	if maxHops > M+2 {
		t.Fatalf("max hops %d exceeds O(log N) bound %d", maxHops, M)
	}
}
