package main

import (
	"math/rand/v2"
	"sync"
)

// State represents whether a node has received the rumor.
type State int

const (
	Susceptible State = iota
	Infected
)

// Node is a single cluster member.
type Node struct {
	ID    int
	State State
}

// Cluster is a set of nodes running gossip.
type Cluster struct {
	mu     sync.Mutex
	nodes  []*Node
	Fanout int
}

// NewCluster creates a cluster of n nodes with the given fanout.
func NewCluster(n, fanout int) *Cluster {
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = &Node{ID: i}
	}
	return &Cluster{nodes: nodes, Fanout: fanout}
}

// Seed marks the given node as infected (the rumor source).
func (c *Cluster) Seed(id int) {
	c.nodes[id].State = Infected
}

// Tick performs one round of gossip: every infected node pushes to Fanout random peers.
func (c *Cluster) Tick() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Snapshot currently infected
	var spreaders []*Node
	for _, n := range c.nodes {
		if n.State == Infected {
			spreaders = append(spreaders, n)
		}
	}

	newInfections := 0
	for _, s := range spreaders {
		for f := 0; f < c.Fanout; f++ {
			target := c.nodes[rand.IntN(len(c.nodes))]
			if target.ID == s.ID {
				continue
			}
			if target.State != Infected {
				target.State = Infected
				newInfections++
			}
		}
	}
	return newInfections
}

// InfectedCount returns the number of infected nodes.
func (c *Cluster) InfectedCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	count := 0
	for _, n := range c.nodes {
		if n.State == Infected {
			count++
		}
	}
	return count
}

// TotalNodes returns the total number of nodes.
func (c *Cluster) TotalNodes() int { return len(c.nodes) }

// ConvergedIn runs gossip until all nodes are infected and returns the number of rounds.
// Returns -1 if maxRounds reached without convergence.
func (c *Cluster) ConvergedIn(maxRounds int) int {
	for r := 1; r <= maxRounds; r++ {
		c.Tick()
		if c.InfectedCount() == c.TotalNodes() {
			return r
		}
	}
	return -1
}
