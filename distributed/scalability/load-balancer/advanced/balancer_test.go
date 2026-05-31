package main

import (
	"errors"
	"testing"
)

func backends() []*Backend {
	return []*Backend{NewBackend("a:80"), NewBackend("b:80"), NewBackend("c:80")}
}

func TestRoundRobin_distribution(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	seen := map[string]int{}
	for i := 0; i < 9; i++ {
		b, err := rr.Next(bs)
		if err != nil {
			t.Fatal(err)
		}
		seen[b.Addr]++
	}
	for _, b := range bs {
		if seen[b.Addr] != 3 {
			t.Errorf("%s: got %d calls, want 3", b.Addr, seen[b.Addr])
		}
	}
}

func TestRoundRobin_skipsUnhealthy(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	bs[1].SetHealthy(false)
	for i := 0; i < 10; i++ {
		b, _ := rr.Next(bs)
		if b.Addr == "b:80" {
			t.Error("selected unhealthy backend b")
		}
	}
}

func TestLeastConns_selectsLowest(t *testing.T) {
	lc := &LeastConns{}
	bs := backends()
	bs[0].ActiveConns.Store(10)
	bs[1].ActiveConns.Store(2)
	bs[2].ActiveConns.Store(5)
	b, _ := lc.Next(bs)
	if b.Addr != "b:80" {
		t.Errorf("got %s, want b:80", b.Addr)
	}
}

func TestBalancer_allUnhealthy(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	for _, b := range bs {
		b.SetHealthy(false)
	}
	if _, err := rr.Next(bs); !errors.Is(err, ErrNoHealthyBackend) {
		t.Error("expected ErrNoHealthyBackend")
	}
}
