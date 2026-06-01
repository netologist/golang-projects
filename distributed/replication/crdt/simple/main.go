package main

import "fmt"

// GCounter: each node can only increment its own counter.
type GCounter struct {
	counters map[string]int
}

func (g *GCounter) Increment(nodeID string) {
	if g.counters == nil {
		g.counters = map[string]int{}
	}
	g.counters[nodeID]++
}

func (g GCounter) Value() int {
	n := 0
	for _, v := range g.counters {
		n += v
	}
	return n
}

func (g GCounter) Merge(other GCounter) GCounter {
	out := GCounter{counters: map[string]int{}}
	for k, v := range g.counters {
		out.counters[k] = v
	}
	for k, v := range other.counters {
		if v > out.counters[k] {
			out.counters[k] = v
		}
	}
	return out
}

func main() {
	node1 := GCounter{counters: map[string]int{}}
	node2 := GCounter{counters: map[string]int{}}
	node3 := GCounter{counters: map[string]int{}}

	node1.Increment("n1")
	node1.Increment("n1")
	node2.Increment("n2")
	node3.Increment("n3")
	node3.Increment("n3")
	node3.Increment("n3")

	fmt.Printf("node1 local=%d\n", node1.Value())
	fmt.Printf("node2 local=%d\n", node2.Value())
	fmt.Printf("node3 local=%d\n", node3.Value())

	// Merge all three
	merged := node1.Merge(node2).Merge(node3)
	fmt.Printf("\nmerged (converged) = %d\n", merged.Value())
}
