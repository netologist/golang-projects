package main

import "fmt"

func main() {
	qs := New(5, 3, 3)
	fmt.Printf("Quorum store: N=%d W=%d R=%d (W+R=%d > N=%d)\n\n",
		qs.N, qs.W, qs.R, qs.W+qs.R, qs.N)

	// Write
	if err := qs.Write("user", "alice", 1); err != nil {
		panic(err)
	}
	fmt.Println("[write] user=alice ver=1 to 5 nodes (W=3 acks required)")

	// Make node 4 stale
	qs.Node(4).write("user", VersionedValue{Value: "stale", Version: 0})
	fmt.Println("[inject] node-4 set stale (ver=0)")

	// Read quorum: should return fresh value
	vv, err := qs.Read("user")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[read] value=%q ver=%d (stale replica ignored)\n", vv.Value, vv.Version)

	// Demonstrate quorum failure
	for i := 0; i < 3; i++ {
		qs.Node(i).Offline = true
	}
	_, err = qs.Read("user")
	fmt.Printf("[read with 3 offline] err=%v\n", err)
}
