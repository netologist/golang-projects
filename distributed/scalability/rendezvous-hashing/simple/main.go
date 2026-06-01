package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

var nodes = []string{"node-1", "node-2", "node-3", "node-4", "node-5"}

func weight(key, node string) uint64 {
	h := md5.Sum([]byte(key + ":" + node))
	return binary.BigEndian.Uint64(h[:8])
}

func assignNode(key string, pool []string) string {
	var best string
	var bestW uint64
	for _, n := range pool {
		if w := weight(key, n); w > bestW {
			bestW, best = w, n
		}
	}
	return best
}

func main() {
	keys := make([]string, 20)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%02d", i)
	}

	before := map[string]string{}
	for _, k := range keys {
		before[k] = assignNode(k, nodes)
	}

	// Remove node-3
	reduced := []string{"node-1", "node-2", "node-4", "node-5"}
	moved := 0
	for _, k := range keys {
		after := assignNode(k, reduced)
		if before[k] != after {
			moved++
			fmt.Printf("  %s: %s → %s (remapped)\n", k, before[k], after)
		}
	}
	fmt.Printf(
		"\nRemoved node-3: %d/%d keys remapped (ideal: only node-3's keys)\n",
		moved, len(keys),
	)
}
