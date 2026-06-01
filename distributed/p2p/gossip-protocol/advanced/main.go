package main

import "fmt"

func main() {
	c := NewCluster(10, 2)
	c.Seed(0)

	fmt.Printf("Gossip: %d nodes, fanout=%d, seed=node-0\n", c.TotalNodes(), c.Fanout)

	for round := 1; ; round++ {
		c.Tick()
		count := c.InfectedCount()
		fmt.Printf("  round %2d: %d/%d infected\n", round, count, c.TotalNodes())
		if count == c.TotalNodes() {
			fmt.Printf("Converged in %d rounds\n", round)
			break
		}
		if round > 100 {
			fmt.Println("Did not converge")
			break
		}
	}
}
