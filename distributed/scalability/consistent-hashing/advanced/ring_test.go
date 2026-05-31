package main

import (
	"fmt"
	"testing"
)

func TestRing_basicGet(t *testing.T) {
	r := New(100)
	r.Add("node1")
	r.Add("node2")
	r.Add("node3")

	for i := 0; i < 10; i++ {
		node, ok := r.Get(fmt.Sprintf("key-%d", i))
		if !ok {
			t.Errorf("key-%d: not found", i)
		}
		if node == "" {
			t.Errorf("key-%d: empty node", i)
		}
	}
}

func TestRing_removeNode(t *testing.T) {
	r := New(100)
	r.Add("a")
	r.Add("b")
	r.Add("c")
	r.Remove("b")

	for i := 0; i < 20; i++ {
		node, _ := r.Get(fmt.Sprintf("k%d", i))
		if node == "b" {
			t.Errorf("key routed to removed node b")
		}
	}
}

func TestRing_distribution(t *testing.T) {
	r := New(150)
	nodes := []string{"n1", "n2", "n3", "n4"}
	for _, n := range nodes {
		r.Add(n)
	}

	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		n, _ := r.Get(fmt.Sprintf("key-%d", i))
		counts[n]++
	}
	for _, n := range nodes {
		pct := float64(counts[n]) / 10.0
		if pct < 10 || pct > 40 {
			t.Errorf("node %s got %.1f%% of keys, expected ~25%%", n, pct)
		}
	}
}

func TestRing_minimalRemapOnAdd(t *testing.T) {
	r := New(150)
	for _, n := range []string{"a", "b", "c"} {
		r.Add(n)
	}
	before := map[string]string{}
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key-%d", i)
		before[k], _ = r.Get(k)
	}
	r.Add("d")
	moved := 0
	for k, oldNode := range before {
		newNode, _ := r.Get(k)
		if newNode != oldNode {
			moved++
		}
	}
	// Adding a 4th node should move roughly 1/4 of keys, well under half.
	if moved > 500 {
		t.Errorf("too many keys moved: %d/1000", moved)
	}
}
