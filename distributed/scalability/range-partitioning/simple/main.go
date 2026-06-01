package main

import (
	"fmt"
	"sort"
	"strings"
)

type Shard struct {
	Name  string
	Start string // inclusive lower bound
	End   string // exclusive upper bound (empty = infinity)
}

func route(shards []Shard, key string) string {
	for _, s := range shards {
		if key >= s.Start && (s.End == "" || key < s.End) {
			return s.Name
		}
	}
	return "none"
}

func main() {
	shards := []Shard{
		{Name: "shard-1", Start: "a", End: "h"},
		{Name: "shard-2", Start: "h", End: "p"},
		{Name: "shard-3", Start: "p", End: ""},
	}

	keys := []string{"apple", "hello", "pear", "grape", "zebra", "moon", "ant"}
	fmt.Println("Initial routing:")
	for _, k := range keys {
		fmt.Printf("  %-10s → %s\n", k, route(shards, k))
	}

	// Add shard-4 by splitting shard-3 at 't'
	shards = append(shards[:2], Shard{Name: "shard-3", Start: "p", End: "t"})
	shards = append(shards, Shard{Name: "shard-4", Start: "t", End: ""})
	sort.Slice(shards, func(i, j int) bool { return shards[i].Start < shards[j].Start })

	fmt.Println("\nAfter adding shard-4 (split at 't'):")
	for _, k := range keys {
		shard := route(shards, k)
		remapped := ""
		if strings.HasPrefix(k, "t") || k >= "t" {
			remapped = " (remapped)"
		}
		fmt.Printf("  %-10s → %s%s\n", k, shard, remapped)
	}
}
