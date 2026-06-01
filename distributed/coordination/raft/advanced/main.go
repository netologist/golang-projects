package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const clusterSize = 5
	nodes := make([]*RaftNode, clusterSize)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	for i, node := range nodes {
		peers := make([]*RaftNode, 0, clusterSize-1)
		for j, p := range nodes {
			if i != j {
				peers = append(peers, p)
			}
		}
		node.SetPeers(peers)
	}

	var wg sync.WaitGroup
	for _, n := range nodes {
		n.Start(&wg)
	}

	// Wait for initial leader
	var leader *RaftNode
	for i := 0; i < 60; i++ {
		time.Sleep(100 * time.Millisecond)
		for _, n := range nodes {
			if !n.isDead() && n.State() == Leader {
				leader = n
				break
			}
		}
		if leader != nil {
			break
		}
	}
	if leader == nil {
		fmt.Println("ERROR: no leader elected")
		return
	}
	fmt.Printf("Initial leader: node-%d (term=%d)\n", leader.ID, leader.Term())

	// Submit entries
	for _, cmd := range []string{"set x=1", "set y=2", "set z=3"} {
		leader.Submit(cmd)
		fmt.Printf("[submit] %s\n", cmd)
	}
	time.Sleep(200 * time.Millisecond)

	// Crash leader
	fmt.Printf("\n[crash] leader node-%d\n", leader.ID)
	leader.Stop()
	time.Sleep(100 * time.Millisecond)

	// Wait for new leader
	var newLeader *RaftNode
	for i := 0; i < 60 && newLeader == nil; i++ {
		time.Sleep(100 * time.Millisecond)
		for _, n := range nodes {
			if !n.isDead() && n.State() == Leader {
				newLeader = n
				break
			}
		}
	}
	if newLeader == nil {
		fmt.Println("ERROR: no new leader")
		return
	}
	fmt.Printf("New leader: node-%d (term=%d)\n", newLeader.ID, newLeader.Term())

	// Submit more entries
	for _, cmd := range []string{"set a=1", "set b=2"} {
		newLeader.Submit(cmd)
		fmt.Printf("[submit] %s\n", cmd)
	}
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\nFinal state:")
	for _, n := range nodes {
		if n.isDead() {
			continue
		}
		fmt.Printf("  node-%d %-10s term=%d log=%v\n", n.ID, n.State(), n.Term(), func() []string {
			var cmds []string
			for _, e := range n.Log() {
				cmds = append(cmds, e.Command)
			}
			return cmds
		}())
	}

	for _, n := range nodes {
		n.Stop()
	}
	wg.Wait()
}
