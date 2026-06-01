package main

import "fmt"

const m = 3 // bits → ring size = 2^m = 8

var ring [1 << m]*int // ring[id] = &id if node exists

func successor(id int) int {
	size := 1 << m
	for i := 1; i <= size; i++ {
		if ring[(id+i)%size] != nil {
			return (id + i) % size
		}
	}
	return id
}

func lookup(key, start int) (node, hops int) {
	size := 1 << m
	cur := start
	for hops = 0; ; hops++ {
		next := successor(cur)
		if next >= key%size || next == cur {
			return next, hops + 1
		}
		cur = next
	}
}

func main() {
	size := 1 << m
	// Add 8 nodes at positions 0-7
	for i := 0; i < size; i++ {
		id := i
		ring[i] = &id
	}

	keys := []int{0, 1, 3, 5, 7, 10, 15}
	fmt.Printf("Chord ring: %d nodes (IDs 0..%d), m=%d\n", size, size-1, m)
	for _, key := range keys {
		node, hops := lookup(key, 0)
		fmt.Printf("  lookup(key=%2d) → node %d  (%d hops)\n", key, node, hops)
	}
}
