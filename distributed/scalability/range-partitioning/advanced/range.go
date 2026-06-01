package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// ErrNoShard is returned when no shard covers the key.
var ErrNoShard = errors.New("range-partition: no shard found")

// Shard describes a key-range partition.
type Shard struct {
	ID    string
	Start string // inclusive
	End   string // exclusive; empty means +∞
}

func (s Shard) contains(key string) bool {
	return key >= s.Start && (s.End == "" || key < s.End)
}

// RangePartitioner routes keys to shards by range.
type RangePartitioner struct {
	mu     sync.RWMutex
	shards []Shard // sorted by Start
}

// New creates a partitioner with initial shards.
func New(shards []Shard) *RangePartitioner {
	p := &RangePartitioner{}
	p.shards = append(p.shards, shards...)
	p.sort()
	return p
}

func (p *RangePartitioner) sort() {
	sort.Slice(p.shards, func(i, j int) bool {
		return p.shards[i].Start < p.shards[j].Start
	})
}

// Route returns the shard ID for key.
func (p *RangePartitioner) Route(key string) (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.shards {
		if s.contains(key) {
			return s.ID, nil
		}
	}
	return "", fmt.Errorf("%w for key %q", ErrNoShard, key)
}

// AddShard inserts a new shard and re-sorts.
func (p *RangePartitioner) AddShard(s Shard) {
	p.mu.Lock()
	p.shards = append(p.shards, s)
	p.sort()
	p.mu.Unlock()
}

// RemoveShard removes the shard with the given ID.
func (p *RangePartitioner) RemoveShard(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, s := range p.shards {
		if s.ID == id {
			p.shards = append(p.shards[:i], p.shards[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("range-partition: shard %q not found", id)
}

// Shards returns a copy of the current shard list.
func (p *RangePartitioner) Shards() []Shard {
	p.mu.RLock()
	out := append([]Shard{}, p.shards...)
	p.mu.RUnlock()
	return out
}

// RebalanceSplit splits shardID at splitAt, creating a new shard newID for [splitAt, oldEnd).
func (p *RangePartitioner) RebalanceSplit(shardID, splitAt, newID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, s := range p.shards {
		if s.ID == shardID {
			if splitAt <= s.Start || (s.End != "" && splitAt >= s.End) {
				return fmt.Errorf(
					"range-partition: split point %q out of range [%s,%s)",
					splitAt, s.Start, s.End,
				)
			}
			oldEnd := s.End
			p.shards[i].End = splitAt
			p.shards = append(p.shards, Shard{ID: newID, Start: splitAt, End: oldEnd})
			p.sort()
			return nil
		}
	}
	return fmt.Errorf("range-partition: shard %q not found", shardID)
}
