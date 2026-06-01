package main

import (
	"fmt"
	"time"
)

func main() {
	l1 := NewLeader("leader-1", LWWResolver)
	l2 := NewLeader("leader-2", LWWResolver)

	// Both leaders write the same key concurrently
	l1.Write("user:1", "alice-from-dc1")
	time.Sleep(10 * time.Millisecond)
	l2.Write("user:1", "alice-from-dc2") // newer timestamp

	// Unique keys (no conflict)
	l1.Write("user:2", "bob")
	l2.Write("user:3", "carol")

	fmt.Println("Before sync:")
	fmt.Printf("  leader-1 user:1 = %q\n", func() string { r, _ := l1.Read("user:1"); return r.Value }())
	fmt.Printf("  leader-2 user:1 = %q\n", func() string { r, _ := l2.Read("user:1"); return r.Value }())

	conflicts := l1.Sync(l2)
	fmt.Printf("\nAfter sync: %d conflict(s) resolved\n", len(conflicts))
	for _, c := range conflicts {
		r, _ := l1.Read(c.Key)
		fmt.Printf("  key=%s local=%q remote=%q winner=%q\n",
			c.Key, c.Local.Value, c.Remote.Value, r.Value)
	}

	fmt.Println("\nFinal leader-1 store:")
	for k, v := range l1.Snapshot() {
		fmt.Printf("  %-10s = %q (origin=%s)\n", k, v.Value, v.Origin)
	}
}
