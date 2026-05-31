package main

import (
	"sync"
	"time"
)

// Cluster holds the shared lease that nodes compete for.
type Cluster struct {
	mu       sync.Mutex
	leaderID string
	expires  time.Time
}

// NewCluster creates an empty cluster lease.
func NewCluster() *Cluster { return &Cluster{} }

// tryAcquire grants or renews the lease for nodeID if it is free or already held
// by this node. Returns true if nodeID holds leadership afterward.
func (c *Cluster) tryAcquire(nodeID string, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	free := c.leaderID == "" || now.After(c.expires)
	if free || c.leaderID == nodeID {
		c.leaderID = nodeID
		c.expires = now.Add(ttl)
		return true
	}
	return false
}

// Leader returns the current leader ID if the lease is still valid.
func (c *Cluster) Leader() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Now().After(c.expires) {
		return ""
	}
	return c.leaderID
}

// Node campaigns for leadership against a shared Cluster.
type Node struct {
	ID       string
	cluster  *Cluster
	ttl      time.Duration
	mu       sync.Mutex
	isLeader bool
}

// NewNode creates a node bound to a cluster.
func NewNode(id string, cluster *Cluster, ttl time.Duration) *Node {
	return &Node{ID: id, cluster: cluster, ttl: ttl}
}

// Campaign attempts to acquire/renew leadership once and reports the result.
func (n *Node) Campaign() bool {
	won := n.cluster.tryAcquire(n.ID, n.ttl)
	n.mu.Lock()
	n.isLeader = won
	n.mu.Unlock()
	return won
}

// IsLeader reports whether this node currently believes it is the leader.
func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.isLeader
}
