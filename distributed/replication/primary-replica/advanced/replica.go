package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Role of a replication node.
type Role int

const (
	Primary Role = iota
	Replica
)

// ErrReadOnPrimary is returned when a write is attempted on a replica.
var ErrReadOnPrimary = errors.New("replica: writes not allowed on replica")

// Entry is a versioned key-value record.
type Entry struct {
	Value   string
	Version uint64
}

// Node is a single replication participant.
type Node struct {
	ID    string
	role  Role
	mu    sync.RWMutex
	store map[string]Entry
	lag   time.Duration
}

// NewPrimary creates a primary node.
func NewPrimary(id string) *Node {
	return &Node{ID: id, role: Primary, store: map[string]Entry{}}
}

// NewReplica creates a replica node with optional replication lag.
func NewReplica(id string, lag time.Duration) *Node {
	return &Node{ID: id, role: Replica, store: map[string]Entry{}, lag: lag}
}

// Write stores a key on the primary and replicates to all replicas.
func (n *Node) Write(key, value string, replicas []*Node) error {
	if n.role != Primary {
		return fmt.Errorf("node %s: %w", n.ID, ErrReadOnPrimary)
	}
	n.mu.Lock()
	var ver uint64
	if e, ok := n.store[key]; ok {
		ver = e.Version + 1
	}
	n.store[key] = Entry{Value: value, Version: ver}
	n.mu.Unlock()

	entry := Entry{Value: value, Version: ver}
	for _, r := range replicas {
		r := r
		go func() {
			time.Sleep(r.lag)
			r.mu.Lock()
			r.store[key] = entry
			r.mu.Unlock()
		}()
	}
	return nil
}

// Read returns the value for key.
func (n *Node) Read(key string) (Entry, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	e, ok := n.store[key]
	return e, ok
}
