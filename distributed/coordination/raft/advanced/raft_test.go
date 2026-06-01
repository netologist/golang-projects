package main

import (
	"sync"
	"testing"
	"time"
)

func buildCluster(n int) ([]*RaftNode, *sync.WaitGroup) {
	nodes := make([]*RaftNode, n)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	for i, node := range nodes {
		peers := make([]*RaftNode, 0, n-1)
		for j, p := range nodes {
			if i != j {
				peers = append(peers, p)
			}
		}
		node.SetPeers(peers)
	}
	wg := &sync.WaitGroup{}
	for _, node := range nodes {
		node.Start(wg)
	}
	return nodes, wg
}

func waitForLeader(nodes []*RaftNode, timeout time.Duration) *RaftNode {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, n := range nodes {
			if !n.isDead() && n.State() == Leader {
				return n
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func stopAll(nodes []*RaftNode, wg *sync.WaitGroup) {
	for _, n := range nodes {
		n.Stop()
	}
	wg.Wait()
}

func TestRaft_LeaderElected(t *testing.T) {
	nodes, wg := buildCluster(5)
	defer stopAll(nodes, wg)

	leader := waitForLeader(nodes, 5*time.Second)
	if leader == nil {
		t.Fatal("no leader elected")
	}

	// Only one leader
	leaderCount := 0
	for _, n := range nodes {
		if n.State() == Leader {
			leaderCount++
		}
	}
	if leaderCount != 1 {
		t.Fatalf("want 1 leader, got %d", leaderCount)
	}
}

func TestRaft_LogReplication(t *testing.T) {
	nodes, wg := buildCluster(5)
	defer stopAll(nodes, wg)

	leader := waitForLeader(nodes, 5*time.Second)
	if leader == nil {
		t.Fatal("no leader")
	}

	leader.Submit("cmd1")
	leader.Submit("cmd2")
	time.Sleep(300 * time.Millisecond)

	// Majority should have the entries
	count := 0
	for _, n := range nodes {
		if len(n.Log()) >= 2 {
			count++
		}
	}
	if count < 3 {
		t.Fatalf("want majority (>=3) replicated, got %d", count)
	}
}

func TestRaft_LeaderCrashReelection(t *testing.T) {
	nodes, wg := buildCluster(5)
	defer stopAll(nodes, wg)

	leader1 := waitForLeader(nodes, 5*time.Second)
	if leader1 == nil {
		t.Fatal("no initial leader")
	}
	term1 := leader1.Term()

	leader1.Stop()
	time.Sleep(100 * time.Millisecond)

	leader2 := waitForLeader(nodes, 5*time.Second)
	if leader2 == nil {
		t.Fatal("no new leader after crash")
	}
	if leader2.ID == leader1.ID {
		t.Fatal("same node re-elected after crash")
	}
	if leader2.Term() <= term1 {
		t.Fatalf("new term %d should be > old term %d", leader2.Term(), term1)
	}
}
