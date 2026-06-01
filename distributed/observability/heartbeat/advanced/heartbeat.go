package main

import (
	"fmt"
	"sync"
	"time"
)

// NodeStatus represents the health state of a node.
type NodeStatus int

const (
	StatusUp      NodeStatus = iota
	StatusSuspect            // missed one interval
	StatusDown               // missed threshold intervals
)

func (s NodeStatus) String() string {
	switch s {
	case StatusUp:
		return "UP"
	case StatusSuspect:
		return "SUSPECT"
	default:
		return "DOWN"
	}
}

// NodeInfo holds liveness data for a single node.
type NodeInfo struct {
	LastBeat  time.Time
	Status    NodeStatus
	MissCount int
}

// Monitor tracks heartbeats and detects failures.
type Monitor struct {
	mu             sync.RWMutex
	nodes          map[string]*NodeInfo
	ttl            time.Duration // max silence before suspect
	downAfter      int           // miss count before marking down
	onStatusChange func(id string, old, new NodeStatus)
	quit           chan struct{}
}

// NewMonitor creates a monitor with given TTL and down-after miss count.
func NewMonitor(
	ttl time.Duration,
	downAfter int,
	onChange func(string, NodeStatus, NodeStatus),
) *Monitor {
	return &Monitor{
		nodes:          make(map[string]*NodeInfo),
		ttl:            ttl,
		downAfter:      downAfter,
		onStatusChange: onChange,
		quit:           make(chan struct{}),
	}
}

// Register adds a node to monitor.
func (m *Monitor) Register(id string) {
	m.mu.Lock()
	m.nodes[id] = &NodeInfo{LastBeat: time.Now(), Status: StatusUp}
	m.mu.Unlock()
}

// Beat records a heartbeat from a node.
func (m *Monitor) Beat(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ni, ok := m.nodes[id]
	if !ok {
		m.nodes[id] = &NodeInfo{LastBeat: time.Now(), Status: StatusUp}
		return
	}
	old := ni.Status
	ni.LastBeat = time.Now()
	ni.MissCount = 0
	ni.Status = StatusUp
	if old != StatusUp && m.onStatusChange != nil {
		m.onStatusChange(id, old, StatusUp)
	}
}

// IsAlive returns true if the node is UP or SUSPECT.
func (m *Monitor) IsAlive(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ni, ok := m.nodes[id]
	return ok && ni.Status != StatusDown
}

// Status returns the current status of a node.
func (m *Monitor) Status(id string) NodeStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if ni, ok := m.nodes[id]; ok {
		return ni.Status
	}
	return StatusDown
}

// Start launches the background liveness check loop.
func (m *Monitor) Start(checkInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.check()
			case <-m.quit:
				return
			}
		}
	}()
}

// Stop halts the background loop.
func (m *Monitor) Stop() { close(m.quit) }

func (m *Monitor) check() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, ni := range m.nodes {
		if time.Since(ni.LastBeat) < m.ttl {
			continue
		}
		ni.MissCount++
		old := ni.Status
		var next NodeStatus
		if ni.MissCount >= m.downAfter {
			next = StatusDown
		} else {
			next = StatusSuspect
		}
		ni.Status = next
		if old != next && m.onStatusChange != nil {
			m.onStatusChange(id, old, next)
		}
		_ = fmt.Sprintf("") // keep import
	}
}
