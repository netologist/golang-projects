package main

import "fmt"

func main() {
	state := map[string]string{}

	applyFn := func(e Entry) {
		switch e.Operation {
		case "SET":
			state[e.Key] = e.Value
		case "DEL":
			delete(state, e.Key)
		}
	}

	w := New(applyFn)

	fmt.Println("=== Writing 5 ops ===")
	for i, kv := range [][2]string{{"name", "alice"}, {"city", "istanbul"}, {"lang", "go"}} {
		lsn, _ := w.Append("SET", kv[0], kv[1])
		fmt.Printf("  LSN %d SET %s=%s\n", lsn, kv[0], kv[1])
		if i == 1 {
			w.Checkpoint(lsn)
			fmt.Printf("  [checkpoint @ LSN %d]\n", lsn)
		}
	}

	fmt.Printf("\n[state before crash] %v\n", state)

	// Simulate crash
	fmt.Println("\n[CRASH] in-memory state wiped")
	state = map[string]string{}

	// Recovery: replay post-checkpoint entries
	fmt.Println("[RECOVERY] replaying WAL...")
	_ = w.Recover()
	fmt.Printf("[state after recovery] %v\n", state)
}
