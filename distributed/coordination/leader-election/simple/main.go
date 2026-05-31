package main

import "fmt"

// Lowest live node ID becomes the leader.
func electLeader(nodes []string, dead map[string]bool) string {
	leader := ""
	for _, n := range nodes {
		if dead[n] {
			continue
		}
		if leader == "" || n < leader {
			leader = n
		}
	}
	return leader
}

func main() {
	nodes := []string{"node-3", "node-1", "node-2"}
	fmt.Println("leader:", electLeader(nodes, nil))
	fmt.Println("after node-1 dies:", electLeader(nodes, map[string]bool{"node-1": true}))
}
