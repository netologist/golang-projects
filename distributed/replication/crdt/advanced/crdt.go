package main

import (
	"time"
)

// ---- G-Counter (grow-only counter) ----

// GCounter allows only increment operations; merge takes per-node max.
type GCounter struct {
	counters map[string]uint64
}

// NewGCounter creates an empty GCounter.
func NewGCounter() *GCounter { return &GCounter{counters: map[string]uint64{}} }

// Increment increments the counter for nodeID.
func (g *GCounter) Increment(nodeID string) { g.counters[nodeID]++ }

// Value returns the total across all nodes.
func (g *GCounter) Value() uint64 {
	var n uint64
	for _, v := range g.counters {
		n += v
	}
	return n
}

// Merge merges other into a new GCounter (taking per-node max).
func (g *GCounter) Merge(other *GCounter) *GCounter {
	out := NewGCounter()
	for k, v := range g.counters {
		out.counters[k] = v
	}
	for k, v := range other.counters {
		if v > out.counters[k] {
			out.counters[k] = v
		}
	}
	return out
}

// ---- PN-Counter (increment + decrement) ----

// PNCounter supports increment and decrement via two G-Counters.
type PNCounter struct {
	P *GCounter // positive
	N *GCounter // negative
}

// NewPNCounter creates an empty PNCounter.
func NewPNCounter() *PNCounter { return &PNCounter{P: NewGCounter(), N: NewGCounter()} }

// Increment increments the positive counter.
func (c *PNCounter) Increment(nodeID string) { c.P.Increment(nodeID) }

// Decrement increments the negative counter.
func (c *PNCounter) Decrement(nodeID string) { c.N.Increment(nodeID) }

// Value returns P - N.
func (c *PNCounter) Value() int64 { return int64(c.P.Value()) - int64(c.N.Value()) }

// Merge merges other into a new PNCounter.
func (c *PNCounter) Merge(other *PNCounter) *PNCounter {
	return &PNCounter{P: c.P.Merge(other.P), N: c.N.Merge(other.N)}
}

// ---- LWW-Register (last-write-wins register) ----

// LWWRegister holds one value, resolved by wall-clock timestamp.
type LWWRegister struct {
	Value string
	TS    time.Time
	Node  string
}

// Set updates the register if ts is newer.
func (r *LWWRegister) Set(value, node string, ts time.Time) {
	if ts.After(r.TS) {
		r.Value, r.TS, r.Node = value, ts, node
	}
}

// Merge returns the register with the later timestamp.
func (r LWWRegister) Merge(other LWWRegister) LWWRegister {
	if !r.TS.Before(other.TS) {
		return r
	}
	return other
}

// ---- 2P-Set (two-phase set: add and remove) ----

// TwoPhaseSet is a grow-only add set + grow-only remove set.
// An element is a member iff it is in add and NOT in remove.
type TwoPhaseSet struct {
	Added   map[string]struct{}
	Removed map[string]struct{}
}

// NewTwoPhaseSet creates an empty 2P-Set.
func NewTwoPhaseSet() *TwoPhaseSet {
	return &TwoPhaseSet{
		Added:   map[string]struct{}{},
		Removed: map[string]struct{}{},
	}
}

// Add adds an element.
func (s *TwoPhaseSet) Add(elem string) { s.Added[elem] = struct{}{} }

// Remove marks an element as removed (must have been added first).
func (s *TwoPhaseSet) Remove(elem string) {
	if _, ok := s.Added[elem]; ok {
		s.Removed[elem] = struct{}{}
	}
}

// Contains returns true if elem is currently a member.
func (s *TwoPhaseSet) Contains(elem string) bool {
	_, added := s.Added[elem]
	_, removed := s.Removed[elem]
	return added && !removed
}

// Merge merges two 2P-Sets by union of both add and remove sets.
func (s *TwoPhaseSet) Merge(other *TwoPhaseSet) *TwoPhaseSet {
	out := NewTwoPhaseSet()
	for k := range s.Added {
		out.Added[k] = struct{}{}
	}
	for k := range other.Added {
		out.Added[k] = struct{}{}
	}
	for k := range s.Removed {
		out.Removed[k] = struct{}{}
	}
	for k := range other.Removed {
		out.Removed[k] = struct{}{}
	}
	return out
}
