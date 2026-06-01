package main

import (
	"fmt"
	"sort"
)

// M is the number of identifier bits (ring size = 2^M).
const M = 6 // 64-node address space

// RingSize is the total address space.
const RingSize = 1 << M

// ChordNode is a single node in the Chord ring.
type ChordNode struct {
	ID          int
	FingerTable [M]int // finger[i] = successor of (ID + 2^i) mod RingSize
	successor   int
	predecessor int
}

// Ring holds all active nodes.
type Ring struct {
	nodes  map[int]*ChordNode // nodeID → node
	sorted []int              // sorted node IDs
}

// NewRing creates an empty ring.
func NewRing() *Ring { return &Ring{nodes: map[int]*ChordNode{}} }

// Join adds a node with the given ID.
func (r *Ring) Join(id int) {
	r.nodes[id] = &ChordNode{ID: id}
	r.sorted = append(r.sorted, id)
	sort.Ints(r.sorted)
	r.stabilize()
}

// stabilize rebuilds finger tables and successor/predecessor for all nodes.
func (r *Ring) stabilize() {
	for _, id := range r.sorted {
		n := r.nodes[id]
		n.successor = r.findSuccessor(id + 1)
		n.predecessor = r.findPredecessor(id)
		for i := 0; i < M; i++ {
			start := (id + (1 << i)) % RingSize
			n.FingerTable[i] = r.findSuccessor(start)
		}
	}
}

// findSuccessor returns the first node ID ≥ id (wrapping around).
func (r *Ring) findSuccessor(id int) int {
	id = id % RingSize
	for _, nid := range r.sorted {
		if nid >= id {
			return nid
		}
	}
	return r.sorted[0] // wrap around
}

// findPredecessor returns the last node ID < id (wrapping around).
func (r *Ring) findPredecessor(id int) int {
	for i := len(r.sorted) - 1; i >= 0; i-- {
		if r.sorted[i] < id {
			return r.sorted[i]
		}
	}
	return r.sorted[len(r.sorted)-1] // wrap
}

// Lookup finds the responsible node for key, returning (nodeID, hops).
func (r *Ring) Lookup(key int) (int, int) {
	key = key % RingSize
	// Start at an arbitrary node (first sorted)
	startID := r.sorted[0]
	return r.lookupFrom(r.nodes[startID], key, 0)
}

func (r *Ring) lookupFrom(n *ChordNode, key, hops int) (int, int) {
	if hops > M+2 {
		return n.ID, hops // safety
	}
	// Node is responsible for its own ID.
	if key == n.ID {
		return n.ID, hops + 1
	}
	// Is key in (n.ID, n.successor]?
	if r.inRange(key, n.ID, n.successor) {
		return n.successor, hops + 1
	}
	// Forward to closest preceding finger
	next := r.closestPrecedingFinger(n, key)
	return r.lookupFrom(r.nodes[next], key, hops+1)
}

func (r *Ring) closestPrecedingFinger(n *ChordNode, key int) int {
	for i := M - 1; i >= 0; i-- {
		f := n.FingerTable[i]
		if r.inRange(f, n.ID, key) {
			return f
		}
	}
	return n.successor
}

// inRange checks if id is in (start, end] on the ring (with wrapping).
func (r *Ring) inRange(id, start, end int) bool {
	if start < end {
		return id > start && id <= end
	}
	// wrapping
	return id > start || id <= end
}

// Nodes returns the sorted node IDs.
func (r *Ring) Nodes() []int { return append([]int{}, r.sorted...) }

var _ = fmt.Sprintf // keep import
