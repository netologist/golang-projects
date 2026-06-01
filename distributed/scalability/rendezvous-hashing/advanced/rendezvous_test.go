package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestGet_Consistent(t *testing.T) {
	r := New([]string{"a", "b", "c"})
	for i := 0; i < 100; i++ {
		n1, _ := r.Get("mykey")
		n2, _ := r.Get("mykey")
		if n1 != n2 {
			t.Fatal("inconsistent routing")
		}
	}
}

func TestGet_MinimalDisruption(t *testing.T) {
	nodes := []string{"n1", "n2", "n3", "n4", "n5"}
	r := New(nodes)

	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	before := map[string]string{}
	for _, k := range keys {
		before[k], _ = r.Get(k)
	}

	_ = r.RemoveNode("n3")

	moved := 0
	for _, k := range keys {
		after, _ := r.Get(k)
		if before[k] != after && before[k] != "n3" {
			moved++ // moved but wasn't on removed node
		}
	}
	if moved > 0 {
		t.Fatalf("minimal disruption violated: %d non-n3 keys remapped", moved)
	}
}

func TestGet_EmptyError(t *testing.T) {
	r := New(nil)
	_, err := r.Get("key")
	if !errors.Is(err, ErrNoNodes) {
		t.Fatalf("want ErrNoNodes, got %v", err)
	}
}

func TestGetN(t *testing.T) {
	r := New([]string{"a", "b", "c", "d", "e"})
	nodes, err := r.GetN("mykey", 3)
	if err != nil || len(nodes) != 3 {
		t.Fatalf("GetN failed: %v %v", nodes, err)
	}
}

var _ = fmt.Sprintf // keep fmt import
