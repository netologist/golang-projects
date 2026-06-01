package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// VersionedValue is a key-value pair with a version.
type VersionedValue struct {
	Value   string
	Version uint64
}

// RepairRecord logs a repair event.
type RepairRecord struct {
	NodeID int
	Key    string
	OldVer uint64
	NewVer uint64
}

// RRNode is a node in the read-repair store.
type RRNode struct {
	ID  int
	mu  sync.RWMutex
	kvs map[string]VersionedValue
}

func newRRNode(id int) *RRNode {
	return &RRNode{ID: id, kvs: map[string]VersionedValue{}}
}

func (n *RRNode) get(key string) VersionedValue {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.kvs[key]
}

func (n *RRNode) set(key string, vv VersionedValue) {
	n.mu.Lock()
	n.kvs[key] = vv
	n.mu.Unlock()
}

// ReadRepairStore performs quorum reads and repairs stale replicas.
type ReadRepairStore struct {
	nodes       []*RRNode
	repairCount int64
	repairLog   []RepairRecord
	repairMu    sync.Mutex
}

// NewReadRepairStore creates a store backed by n nodes.
func NewReadRepairStore(n int) *ReadRepairStore {
	nodes := make([]*RRNode, n)
	for i := range nodes {
		nodes[i] = newRRNode(i)
	}
	return &ReadRepairStore{nodes: nodes}
}

// Node returns the i-th node (for test manipulation).
func (s *ReadRepairStore) Node(i int) *RRNode { return s.nodes[i] }

// Write writes key=value with version to all nodes.
func (s *ReadRepairStore) Write(key string, vv VersionedValue) {
	for _, n := range s.nodes {
		n.set(key, vv)
	}
}

// Read performs a read, returning the highest-version value, then
// asynchronously repairs any stale replicas.
func (s *ReadRepairStore) Read(key string) VersionedValue {
	responses := make([]VersionedValue, len(s.nodes))
	for i, n := range s.nodes {
		responses[i] = n.get(key)
	}

	// Find highest version
	var best VersionedValue
	for _, vv := range responses {
		if vv.Version > best.Version {
			best = vv
		}
	}

	// Repair stale nodes asynchronously
	go func() {
		for i, vv := range responses {
			if vv.Version < best.Version {
				old := vv.Version
				s.nodes[i].set(key, best)
				atomic.AddInt64(&s.repairCount, 1)
				s.repairMu.Lock()
				s.repairLog = append(s.repairLog, RepairRecord{
					NodeID: s.nodes[i].ID,
					Key:    key,
					OldVer: old,
					NewVer: best.Version,
				})
				s.repairMu.Unlock()
				fmt.Printf("[repair] node-%d key=%s ver %d→%d\n",
					s.nodes[i].ID, key, old, best.Version)
			}
		}
	}()

	return best
}

// RepairCount returns the total number of repairs performed.
func (s *ReadRepairStore) RepairCount() int64 {
	return atomic.LoadInt64(&s.repairCount)
}

// RepairLog returns all recorded repair events.
func (s *ReadRepairStore) RepairLog() []RepairRecord {
	s.repairMu.Lock()
	defer s.repairMu.Unlock()
	return append([]RepairRecord{}, s.repairLog...)
}
