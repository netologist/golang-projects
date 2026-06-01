package main

import (
	"errors"
	"fmt"
	"sync"
)

// ErrQuorumFailed is returned when not enough nodes acknowledge.
var ErrQuorumFailed = errors.New("quorum: insufficient acknowledgements")

// VersionedValue holds a value with a monotonic version number.
type VersionedValue struct {
	Value   string
	Version uint64
}

// QuorumNode is a single node in a quorum cluster.
type QuorumNode struct {
	ID      int
	mu      sync.RWMutex
	store   map[string]VersionedValue
	Offline bool // set true to simulate failure
}

func newNode(id int) *QuorumNode {
	return &QuorumNode{ID: id, store: map[string]VersionedValue{}}
}

func (n *QuorumNode) write(key string, vv VersionedValue) error {
	if n.Offline {
		return fmt.Errorf("node %d offline", n.ID)
	}
	n.mu.Lock()
	n.store[key] = vv
	n.mu.Unlock()
	return nil
}

func (n *QuorumNode) read(key string) (VersionedValue, error) {
	if n.Offline {
		return VersionedValue{}, fmt.Errorf("node %d offline", n.ID)
	}
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.store[key], nil
}

// QuorumStore orchestrates W-of-N writes and R-of-N reads.
type QuorumStore struct {
	N, W, R int
	nodes   []*QuorumNode
}

// New creates a QuorumStore with N nodes.
func New(n, w, r int) *QuorumStore {
	nodes := make([]*QuorumNode, n)
	for i := range nodes {
		nodes[i] = newNode(i)
	}
	return &QuorumStore{N: n, W: w, R: r, nodes: nodes}
}

// Node returns the i-th node (for test manipulation).
func (qs *QuorumStore) Node(i int) *QuorumNode { return qs.nodes[i] }

// Write writes key=value with the given version to W nodes.
func (qs *QuorumStore) Write(key, value string, ver uint64) error {
	vv := VersionedValue{Value: value, Version: ver}
	acks := 0
	for _, n := range qs.nodes {
		if err := n.write(key, vv); err == nil {
			acks++
		}
	}
	if acks < qs.W {
		return fmt.Errorf("%w: got %d, need %d", ErrQuorumFailed, acks, qs.W)
	}
	return nil
}

// Read collects R responses and returns the value with the highest version.
func (qs *QuorumStore) Read(key string) (VersionedValue, error) {
	var best VersionedValue
	acks := 0
	for _, n := range qs.nodes {
		vv, err := n.read(key)
		if err != nil {
			continue
		}
		if vv.Version > best.Version {
			best = vv
		}
		acks++
	}
	if acks < qs.R {
		return VersionedValue{}, fmt.Errorf("%w: got %d, need %d", ErrQuorumFailed, acks, qs.R)
	}
	return best, nil
}
