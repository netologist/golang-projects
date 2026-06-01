package main

import "fmt"

func main() {
	r := NewRing()
	nodeIDs := []int{0, 8, 16, 24, 32, 40, 48, 56}
	for _, id := range nodeIDs {
		r.Join(id)
	}

	fmt.Printf("Chord DHT: %d nodes, M=%d, ring size=%d\n", len(nodeIDs), M, RingSize)
	fmt.Printf("Nodes: %v\n\n", r.Nodes())

	// Print finger table for node 0
	n0 := r.nodes[0]
	fmt.Println("Finger table for node 0:")
	for i, f := range n0.FingerTable {
		fmt.Printf("  finger[%d] start=%2d successor=%d\n", i, (0+(1<<i))%RingSize, f)
	}

	// Lookup 20 keys
	fmt.Println("\nKey lookups:")
	totalHops := 0
	keys := make([]int, 20)
	for i := range keys {
		keys[i] = i * 3
	}
	for _, key := range keys {
		node, hops := r.Lookup(key)
		totalHops += hops
		fmt.Printf("  key=%2d → node %2d (%d hops)\n", key, node, hops)
	}
	fmt.Printf("\nAverage hops: %.2f (O(log %d) = %.2f)\n",
		float64(totalHops)/float64(len(keys)),
		len(nodeIDs),
		log2(float64(len(nodeIDs))),
	)
}

func log2(x float64) float64 {
	if x <= 0 {
		return 0
	}
	n := 0.0
	for x > 1 {
		x /= 2
		n++
	}
	return n
}
