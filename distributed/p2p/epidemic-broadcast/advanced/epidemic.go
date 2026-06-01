package main

import (
	"math/rand/v2"
	"sync"
)

// NodeStatus models SI epidemic states plus Dead.
type NodeStatus int

const (
	Susceptible NodeStatus = iota
	Infected
	Dead
)

// EpiNode is a single node in the epidemic model.
type EpiNode struct {
	ID     int
	Status NodeStatus
}

// EpidemicBroadcast runs SI epidemic spreading across a set of nodes.
type EpidemicBroadcast struct {
	mu     sync.Mutex
	nodes  []*EpiNode
	Fanout int
}

// New creates an epidemic broadcast with n nodes and given fanout.
func New(n, fanout int) *EpidemicBroadcast {
	nodes := make([]*EpiNode, n)
	for i := range nodes {
		nodes[i] = &EpiNode{ID: i, Status: Susceptible}
	}
	return &EpidemicBroadcast{nodes: nodes, Fanout: fanout}
}

// Infect seeds a node as infected.
func (e *EpidemicBroadcast) Infect(id int) {
	e.nodes[id].Status = Infected
}

// Kill marks a node as dead (unreachable).
func (e *EpidemicBroadcast) Kill(id int) {
	e.nodes[id].Status = Dead
}

// Spread performs one round of epidemic spreading.
func (e *EpidemicBroadcast) Spread() {
	e.mu.Lock()
	defer e.mu.Unlock()

	var spreaders []*EpiNode
	for _, n := range e.nodes {
		if n.Status == Infected {
			spreaders = append(spreaders, n)
		}
	}
	for range spreaders {
		for f := 0; f < e.Fanout; f++ {
			target := e.nodes[rand.IntN(len(e.nodes))]
			if target.Status == Susceptible {
				target.Status = Infected
			}
		}
	}
}

// Converged returns true when all live nodes are infected.
func (e *EpidemicBroadcast) Converged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, n := range e.nodes {
		if n.Status == Susceptible {
			return false
		}
	}
	return true
}

// Stats returns (infected, dead, total) counts.
func (e *EpidemicBroadcast) Stats() (infected, dead, total int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	total = len(e.nodes)
	for _, n := range e.nodes {
		switch n.Status {
		case Infected:
			infected++
		case Dead:
			dead++
		}
	}
	return
}

// StatusOf returns the status of node id.
func (e *EpidemicBroadcast) StatusOf(id int) NodeStatus {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.nodes[id].Status
}
