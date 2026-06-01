package main

import (
	"errors"
	"testing"
)

func TestRoute_Basic(t *testing.T) {
	p := New([]Shard{
		{ID: "s1", Start: "a", End: "m"},
		{ID: "s2", Start: "m", End: ""},
	})
	tests := []struct {
		key  string
		want string
	}{
		{"apple", "s1"},
		{"mango", "s2"},
		{"zebra", "s2"},
	}
	for _, tt := range tests {
		got, err := p.Route(tt.key)
		if err != nil || got != tt.want {
			t.Errorf("Route(%q) = %q %v, want %q", tt.key, got, err, tt.want)
		}
	}
}

func TestRoute_NoShard(t *testing.T) {
	p := New([]Shard{{ID: "s1", Start: "m", End: "z"}})
	_, err := p.Route("apple")
	if !errors.Is(err, ErrNoShard) {
		t.Fatalf("want ErrNoShard, got %v", err)
	}
}

func TestRebalanceSplit(t *testing.T) {
	p := New([]Shard{{ID: "s1", Start: "a", End: ""}})
	if err := p.RebalanceSplit("s1", "m", "s2"); err != nil {
		t.Fatal(err)
	}
	s1, _ := p.Route("apple")
	s2, _ := p.Route("zebra")
	if s1 != "s1" || s2 != "s2" {
		t.Fatalf("split routing wrong: apple→%s zebra→%s", s1, s2)
	}
}
