package main

import (
	"fmt"
	"time"
)

func main() {
	cluster := NewCluster()
	n1 := NewNode("node-1", cluster, 50*time.Millisecond)
	n2 := NewNode("node-2", cluster, 50*time.Millisecond)

	n1.Campaign()
	n2.Campaign()
	fmt.Printf("leader=%s (n1=%v n2=%v)\n", cluster.Leader(), n1.IsLeader(), n2.IsLeader())

	fmt.Println("node-1 stops renewing; waiting for lease to expire...")
	time.Sleep(60 * time.Millisecond)

	n2.Campaign()
	fmt.Printf("new leader=%s (n2=%v)\n", cluster.Leader(), n2.IsLeader())
}
