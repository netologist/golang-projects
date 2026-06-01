package main

import "fmt"

func main() {
	p := New([]Shard{
		{ID: "shard-A", Start: "a", End: "h"},
		{ID: "shard-B", Start: "h", End: "p"},
		{ID: "shard-C", Start: "p", End: ""},
	})

	keys := []string{
		"alice", "bob", "charlie", "dave", "eve", "frank",
		"heidi", "ivan", "judy", "ken", "larry", "mallory",
		"nancy", "oscar", "peggy", "quinn", "roger", "sybil",
		"trent", "victor",
	}

	counts := map[string]int{}
	for _, k := range keys {
		shard, _ := p.Route(k)
		counts[shard]++
	}
	fmt.Println("Initial distribution:")
	for _, s := range p.Shards() {
		fmt.Printf("  %-10s [%s, %s) → %d keys\n", s.ID, s.Start, s.End, counts[s.ID])
	}

	// Split shard-C at 't'
	_ = p.RebalanceSplit("shard-C", "t", "shard-D")

	newCounts := map[string]int{}
	remapped := 0
	for _, k := range keys {
		shard, _ := p.Route(k)
		newCounts[shard]++
		if shard != (func() string { s, _ := p.Route(k); return s }()) {
			remapped++
		}
	}
	fmt.Println("\nAfter splitting shard-C at 't':")
	for _, s := range p.Shards() {
		fmt.Printf("  %-10s [%s, %s) → %d keys\n", s.ID, s.Start, s.End, newCounts[s.ID])
	}
}
