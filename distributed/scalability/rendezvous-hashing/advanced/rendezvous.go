package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"sync"
)

// ErrNoNodes is returned when the node pool is empty.
var ErrNoNodes = errors.New("rendezvous: no nodes available")

// Rendezvous implements Highest Random Weight (HRW) hashing.
type Rendezvous struct {
	mu    sync.RWMutex
	nodes []string // sorted for determinism
}

// New creates a Rendezvous hasher with the given nodes.
func New(nodes []string) *Rendezvous {
	r := &Rendezvous{}
	r.nodes = append(r.nodes, nodes...)
	sort.Strings(r.nodes)
	return r
}

func weight(key, node string) uint64 {
	h := sha256.Sum256([]byte(key + "\x00" + node))
	return binary.BigEndian.Uint64(h[:8])
}

// Get returns the node with the highest weight for key.
func (r *Rendezvous) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return "", ErrNoNodes
	}
	var best string
	var bestW uint64
	for _, n := range r.nodes {
		if w := weight(key, n); w > bestW {
			bestW, best = w, n
		}
	}
	return best, nil
}

// GetN returns the top-n nodes for key (for replication).
func (r *Rendezvous) GetN(key string, n int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return nil, ErrNoNodes
	}
	type weighted struct {
		node string
		w    uint64
	}
	ws := make([]weighted, len(r.nodes))
	for i, nd := range r.nodes {
		ws[i] = weighted{nd, weight(key, nd)}
	}
	sort.Slice(ws, func(i, j int) bool { return ws[i].w > ws[j].w })
	if n > len(ws) {
		n = len(ws)
	}
	out := make([]string, n)
	for i := range out {
		out[i] = ws[i].node
	}
	return out, nil
}

// AddNode adds a node to the pool.
func (r *Rendezvous) AddNode(node string) {
	r.mu.Lock()
	r.nodes = append(r.nodes, node)
	sort.Strings(r.nodes)
	r.mu.Unlock()
}

// RemoveNode removes a node from the pool.
func (r *Rendezvous) RemoveNode(node string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, n := range r.nodes {
		if n == node {
			r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("rendezvous: node %q not found", node)
}

// Nodes returns the current node list.
func (r *Rendezvous) Nodes() []string {
	r.mu.RLock()
	out := append([]string{}, r.nodes...)
	r.mu.RUnlock()
	return out
}
