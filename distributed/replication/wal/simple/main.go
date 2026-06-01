package main

import (
	"fmt"
	"strings"
)

type entry struct {
	lsn int
	op  string
}

func main() {
	var wal []entry
	state := map[string]string{}
	lsn := 0

	apply := func(op string) {
		lsn++
		wal = append(wal, entry{lsn, op})
		parts := strings.SplitN(op, "=", 2)
		state[parts[0]] = parts[1]
		fmt.Printf("[wal] LSN %d: %s\n", lsn, op)
	}

	apply("name=alice")
	apply("city=istanbul")
	apply("lang=go")

	fmt.Printf("\n[state] %v\n", state)

	// Simulate crash: wipe in-memory state
	fmt.Println("\n[crash] state lost!")
	state = map[string]string{}

	// Recover by replaying WAL
	fmt.Println("[recovery] replaying WAL...")
	for _, e := range wal {
		parts := strings.SplitN(e.op, "=", 2)
		state[parts[0]] = parts[1]
		fmt.Printf("  replayed LSN %d: %s\n", e.lsn, e.op)
	}
	fmt.Printf("[recovered] %v\n", state)
}
