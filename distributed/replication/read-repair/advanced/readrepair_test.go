package main

import (
	"testing"
	"time"
)

func TestReadRepair_StaleReplica(t *testing.T) {
	s := NewReadRepairStore(3)
	// Write fresh to nodes 0 and 1, stale to node 2
	s.Node(0).set("k", VersionedValue{"fresh", 2})
	s.Node(1).set("k", VersionedValue{"fresh", 2})
	s.Node(2).set("k", VersionedValue{"stale", 1})

	vv := s.Read("k")
	if vv.Value != "fresh" {
		t.Fatalf("want fresh, got %q", vv.Value)
	}
	time.Sleep(50 * time.Millisecond) // let repair run
	repaired := s.Node(2).get("k")
	if repaired.Value != "fresh" {
		t.Fatalf("node-2 should be repaired, got %q", repaired.Value)
	}
	if s.RepairCount() < 1 {
		t.Fatal("expected at least 1 repair")
	}
}

func TestReadRepair_NoRepairNeeded(t *testing.T) {
	s := NewReadRepairStore(3)
	vv := VersionedValue{"alice", 1}
	s.Write("k", vv)
	res := s.Read("k")
	if res.Value != "alice" {
		t.Fatalf("want alice, got %q", res.Value)
	}
	time.Sleep(30 * time.Millisecond)
	if s.RepairCount() != 0 {
		t.Fatalf("expected 0 repairs, got %d", s.RepairCount())
	}
}
