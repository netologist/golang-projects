package main

import (
	"errors"
	"testing"
)

func TestQuorum_WriteRead(t *testing.T) {
	qs := New(5, 3, 3)
	if err := qs.Write("k", "v", 1); err != nil {
		t.Fatal(err)
	}
	vv, err := qs.Read("k")
	if err != nil || vv.Value != "v" {
		t.Fatalf("want v, got %v %v", vv, err)
	}
}

func TestQuorum_StaleReplica(t *testing.T) {
	qs := New(5, 3, 3)
	// Stale node 4 gets old version
	qs.Node(4).Offline = true
	_ = qs.Write("k", "fresh", 2)
	qs.Node(4).Offline = false
	// Manually inject stale data on node 4
	qs.Node(4).write("k", VersionedValue{Value: "stale", Version: 1})

	vv, err := qs.Read("k")
	if err != nil || vv.Value != "fresh" {
		t.Fatalf("want fresh, got %v %v", vv, err)
	}
}

func TestQuorum_InsufficientNodes(t *testing.T) {
	qs := New(5, 3, 3)
	// Take 3 nodes offline
	qs.Node(0).Offline = true
	qs.Node(1).Offline = true
	qs.Node(2).Offline = true
	err := qs.Write("k", "v", 1)
	if !errors.Is(err, ErrQuorumFailed) {
		t.Fatalf("want ErrQuorumFailed, got %v", err)
	}
}
