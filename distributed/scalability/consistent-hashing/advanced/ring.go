package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"
)

// Ring is a thread-safe consistent hash ring with virtual nodes.
type Ring struct {
	mu     sync.RWMutex
	vnodes int
	ring   []uint32
	nodes  map[uint32]string
}

// New creates a ring with the given number of virtual nodes per physical node.
func New(vnodes int) *Ring {
	return &Ring{vnodes: vnodes, nodes: map[uint32]string{}}
}

func (r *Ring) hash(key string) uint32 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint32(h[:4])
}

// Add inserts a physical node and its virtual points.
func (r *Ring) Add(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < r.vnodes; i++ {
		h := r.hash(fmt.Sprintf("%s-%d", node, i))
		r.ring = append(r.ring, h)
		r.nodes[h] = node
	}
	sort.Slice(r.ring, func(i, j int) bool { return r.ring[i] < r.ring[j] })
}

// Remove deletes a node and all its virtual points.
func (r *Ring) Remove(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newRing := r.ring[:0]
	for _, h := range r.ring {
		if r.nodes[h] != node {
			newRing = append(newRing, h)
		} else {
			delete(r.nodes, h)
		}
	}
	r.ring = newRing
}

// Get returns the node responsible for key.
func (r *Ring) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.ring) == 0 {
		return "", false
	}
	h := r.hash(key)
	idx := sort.Search(len(r.ring), func(i int) bool { return r.ring[i] >= h })
	if idx == len(r.ring) {
		idx = 0
	}
	return r.nodes[r.ring[idx]], true
}
