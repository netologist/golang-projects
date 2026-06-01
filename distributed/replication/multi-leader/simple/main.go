package main

import "fmt"

type record struct {
	value string
	ts    int64
}

func lww(a, b record) record {
	if a.ts >= b.ts {
		return a
	}
	return b
}

func main() {
	// Two leaders write the same key concurrently (different timestamps)
	leader1 := map[string]record{"user:1": {"alice-v1", 100}}
	leader2 := map[string]record{"user:1": {"alice-v2", 200}}

	fmt.Printf("leader1 user:1 = %+v\n", leader1["user:1"])
	fmt.Printf("leader2 user:1 = %+v\n", leader2["user:1"])

	// Conflict detected on sync
	resolved := lww(leader1["user:1"], leader2["user:1"])
	fmt.Printf("LWW resolved  = %+v (leader2 wins)\n", resolved)

	// No conflict key
	leader1["user:2"] = record{"bob", 50}
	leader2["user:3"] = record{"carol", 60}
	merged := map[string]record{}
	for k, v := range leader1 {
		merged[k] = v
	}
	for k, v := range leader2 {
		if existing, ok := merged[k]; ok {
			merged[k] = lww(existing, v)
		} else {
			merged[k] = v
		}
	}
	fmt.Printf("merged store: %v\n", merged)
}
