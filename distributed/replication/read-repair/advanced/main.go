package main

import (
	"fmt"
	"time"
)

func main() {
	s := NewReadRepairStore(5)

	// Initial write to all nodes
	s.Write("user", VersionedValue{"alice-v1", 1})
	fmt.Println("[write] user=alice-v1 ver=1 to all 5 nodes")

	// Simulate stale replicas (nodes 3 and 4 missed an update)
	s.Node(3).set("user", VersionedValue{"alice-v0", 0})
	s.Node(4).set("user", VersionedValue{"alice-v0", 0})
	fmt.Println("[inject] nodes 3 and 4 set stale (ver=0)")

	// Read: should return fresh value and trigger background repair
	result := s.Read("user")
	fmt.Printf("[read] value=%q ver=%d\n", result.Value, result.Version)

	// Wait for async repair to complete
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n[repair summary] %d repair(s) applied:\n", s.RepairCount())
	for _, r := range s.RepairLog() {
		fmt.Printf("  node-%d key=%s ver %d→%d\n", r.NodeID, r.Key, r.OldVer, r.NewVer)
	}

	fmt.Println("\n[post-repair reads]")
	for i := 0; i < 5; i++ {
		vv := s.Node(i).get("user")
		fmt.Printf("  node-%d: value=%q ver=%d\n", i, vv.Value, vv.Version)
	}
}
