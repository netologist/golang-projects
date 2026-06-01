package main

import (
	"fmt"
	"time"
)

type node struct{ data map[string]string }

func quorumRead(nodes []*node, key string) string {
	// Collect responses
	values := map[string]int{}
	for _, n := range nodes {
		v := n.data[key]
		values[v]++
	}
	// Find majority
	var best string
	bestCount := 0
	for v, c := range values {
		if c > bestCount {
			bestCount, best = c, v
		}
	}
	// Repair stale replicas in background
	go func() {
		for _, n := range nodes {
			if n.data[key] != best {
				fmt.Printf("[repair] repairing stale replica\n")
				n.data[key] = best
			}
		}
	}()
	return best
}

func main() {
	n1 := &node{data: map[string]string{"name": "alice"}}
	n2 := &node{data: map[string]string{"name": "alice"}}
	n3 := &node{data: map[string]string{"name": "stale"}} // stale

	nodes := []*node{n1, n2, n3}
	v := quorumRead(nodes, "name")
	fmt.Printf("quorum read: %q\n", v)

	time.Sleep(50 * time.Millisecond) // let repair finish
	fmt.Printf("n3 after repair: %q\n", n3.data["name"])
}
