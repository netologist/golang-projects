package main

import (
	"fmt"
	"time"
)

func main() {
	primary := NewPrimary("primary")
	r1 := NewReplica("replica-1", 30*time.Millisecond)
	r2 := NewReplica("replica-2", 80*time.Millisecond)
	replicas := []*Node{r1, r2}

	// Write 3 keys
	for _, kv := range [][2]string{{"name", "alice"}, {"city", "istanbul"}, {"lang", "go"}} {
		if err := primary.Write(kv[0], kv[1], replicas); err != nil {
			panic(err)
		}
		fmt.Printf("[primary] write %s=%s\n", kv[0], kv[1])
	}

	fmt.Println("\n[before lag] replica reads:")
	for _, r := range replicas {
		e, ok := r.Read("name")
		fmt.Printf("  %s name=%q ver=%d propagated=%v\n", r.ID, e.Value, e.Version, ok)
	}

	time.Sleep(150 * time.Millisecond)
	fmt.Println("\n[after lag] replica reads:")
	for _, r := range replicas {
		for _, k := range []string{"name", "city", "lang"} {
			e, _ := r.Read(k)
			fmt.Printf("  %s %s=%q ver=%d\n", r.ID, k, e.Value, e.Version)
		}
	}
}
