package main

import (
	"testing"
	"time"
)

func TestElection_singleLeader(t *testing.T) {
	cluster := NewCluster()
	n1 := NewNode("node-1", cluster, 100*time.Millisecond)
	n2 := NewNode("node-2", cluster, 100*time.Millisecond)

	if !n1.Campaign() {
		t.Fatal("node-1 should win the empty lease")
	}
	if n2.Campaign() {
		t.Fatal("node-2 should not win while node-1 holds the lease")
	}
	if cluster.Leader() != "node-1" {
		t.Errorf("leader: got %s, want node-1", cluster.Leader())
	}
}

func TestElection_failover(t *testing.T) {
	cluster := NewCluster()
	n1 := NewNode("node-1", cluster, 20*time.Millisecond)
	n2 := NewNode("node-2", cluster, 20*time.Millisecond)

	n1.Campaign()                     // n1 leads
	time.Sleep(30 * time.Millisecond) // n1 "crashes": lease expires without renewal

	if !n2.Campaign() {
		t.Fatal("node-2 should take over after lease expiry")
	}
	if cluster.Leader() != "node-2" {
		t.Errorf("leader after failover: got %s, want node-2", cluster.Leader())
	}
}
