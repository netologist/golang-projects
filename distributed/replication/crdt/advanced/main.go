package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== G-Counter (3 nodes) ===")
	n1, n2, n3 := NewGCounter(), NewGCounter(), NewGCounter()
	n1.Increment("n1")
	n1.Increment("n1")
	n2.Increment("n2")
	n3.Increment("n3")
	n3.Increment("n3")
	n3.Increment("n3")
	merged := n1.Merge(n2).Merge(n3)
	fmt.Printf("  n1=%d n2=%d n3=%d merged=%d\n", n1.Value(), n2.Value(), n3.Value(), merged.Value())

	fmt.Println("\n=== PN-Counter ===")
	pn := NewPNCounter()
	pn.Increment("n1")
	pn.Increment("n1")
	pn.Increment("n1")
	pn.Decrement("n1")
	fmt.Printf("  value=%d\n", pn.Value())

	fmt.Println("\n=== LWW-Register (3 nodes) ===")
	var r1, r2, r3 LWWRegister
	now := time.Now()
	r1.Set("alice-v1", "n1", now)
	r2.Set("alice-v2", "n2", now.Add(5*time.Millisecond))
	r3.Set("alice-v3", "n3", now.Add(2*time.Millisecond))
	result := r1.Merge(r2).Merge(r3)
	fmt.Printf("  n1=%q n2=%q n3=%q merged=%q (latest wins)\n",
		r1.Value, r2.Value, r3.Value, result.Value)

	fmt.Println("\n=== 2P-Set (3 nodes) ===")
	s1, s2, s3 := NewTwoPhaseSet(), NewTwoPhaseSet(), NewTwoPhaseSet()
	s1.Add("alice")
	s1.Add("bob")
	s2.Add("carol")
	s2.Add("alice")
	s2.Remove("alice")
	s3.Add("dave")
	final := s1.Merge(s2).Merge(s3)
	fmt.Printf("  alice=%v bob=%v carol=%v dave=%v\n",
		final.Contains("alice"), final.Contains("bob"),
		final.Contains("carol"), final.Contains("dave"))
}
