package main

import "fmt"

func main() {
	nodeList := []string{"node-1", "node-2", "node-3", "node-4", "node-5"}
	r := New(nodeList)

	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%03d", i)
	}

	dist := map[string]int{}
	for _, k := range keys {
		n, _ := r.Get(k)
		dist[n]++
	}
	fmt.Println("Initial distribution (100 keys, 5 nodes):")
	for _, n := range r.Nodes() {
		fmt.Printf("  %-10s %d keys\n", n, dist[n])
	}

	// Remove node-3
	_ = r.RemoveNode("node-3")

	moved := 0
	for _, k := range keys {
		before, _ := func() (string, error) {
			tmp := New(nodeList)
			return tmp.Get(k)
		}()
		after, _ := r.Get(k)
		if before != after {
			moved++
		}
	}
	fmt.Printf("\nRemoved node-3: %d/%d keys remapped\n", moved, len(keys))

	newDist := map[string]int{}
	for _, k := range keys {
		n, _ := r.Get(k)
		newDist[n]++
	}
	fmt.Println("New distribution (4 nodes):")
	for _, n := range r.Nodes() {
		fmt.Printf("  %-10s %d keys\n", n, newDist[n])
	}
}
