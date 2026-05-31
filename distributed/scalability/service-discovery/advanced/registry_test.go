package main

import (
	"testing"
	"time"
)

func TestRegistry_registerDiscover(t *testing.T) {
	r := New()
	defer r.Close()

	dereg := r.Register(Instance{ID: "i1", Service: "api", Addr: "10.0.0.1:80", TTL: time.Second})
	defer dereg()

	insts, err := r.Discover("api")
	if err != nil {
		t.Fatal(err)
	}
	if len(insts) != 1 || insts[0].Addr != "10.0.0.1:80" {
		t.Errorf("unexpected instances: %+v", insts)
	}
}

func TestRegistry_ttlExpiry(t *testing.T) {
	r := New()
	defer r.Close()

	r.Register(Instance{ID: "i1", Service: "api", Addr: "x", TTL: 30 * time.Millisecond})
	time.Sleep(120 * time.Millisecond) // past TTL + a reaper tick

	insts, _ := r.Discover("api")
	if len(insts) != 0 {
		t.Errorf("expected instance to expire, got %d", len(insts))
	}
}

func TestRegistry_deregister(t *testing.T) {
	r := New()
	defer r.Close()

	dereg := r.Register(Instance{ID: "i1", Service: "api", Addr: "x", TTL: time.Second})
	dereg()

	insts, _ := r.Discover("api")
	if len(insts) != 0 {
		t.Errorf("expected 0 after deregister, got %d", len(insts))
	}
}
