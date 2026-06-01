package main

import (
	"fmt"
	"sync"
	"time"
)

// Record is a versioned value with a wall-clock timestamp.
type Record struct {
	Value     string
	Timestamp time.Time
	Origin    string // leader that wrote it
}

// Conflict describes a detected divergence between two records.
type Conflict struct {
	Key    string
	Local  Record
	Remote Record
}

// ConflictResolver chooses the winner among two conflicting records.
type ConflictResolver func(Conflict) Record

// LWWResolver implements Last-Write-Wins conflict resolution.
func LWWResolver(c Conflict) Record {
	if !c.Local.Timestamp.Before(c.Remote.Timestamp) {
		return c.Local
	}
	return c.Remote
}

// Leader is a multi-leader replication node.
type Leader struct {
	ID       string
	mu       sync.RWMutex
	store    map[string]Record
	resolver ConflictResolver
}

// NewLeader creates a leader with the given conflict resolver.
func NewLeader(id string, resolver ConflictResolver) *Leader {
	return &Leader{
		ID:       id,
		store:    map[string]Record{},
		resolver: resolver,
	}
}

// Write stores a key with the current time.
func (l *Leader) Write(key, value string) {
	l.mu.Lock()
	l.store[key] = Record{Value: value, Timestamp: time.Now(), Origin: l.ID}
	l.mu.Unlock()
}

// Read returns the record for key.
func (l *Leader) Read(key string) (Record, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	r, ok := l.store[key]
	return r, ok
}

// Sync merges remote's store into this leader, resolving conflicts.
func (l *Leader) Sync(remote *Leader) []Conflict {
	remote.mu.RLock()
	snap := make(map[string]Record, len(remote.store))
	for k, v := range remote.store {
		snap[k] = v
	}
	remote.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	var conflicts []Conflict
	for k, remoteRec := range snap {
		if localRec, ok := l.store[k]; ok && localRec.Origin != remoteRec.Origin {
			c := Conflict{Key: k, Local: localRec, Remote: remoteRec}
			conflicts = append(conflicts, c)
			l.store[k] = l.resolver(c)
		} else if !ok {
			l.store[k] = remoteRec
		}
	}
	return conflicts
}

// Snapshot returns a copy of the store.
func (l *Leader) Snapshot() map[string]Record {
	l.mu.RLock()
	out := make(map[string]Record, len(l.store))
	for k, v := range l.store {
		out[k] = v
	}
	l.mu.RUnlock()
	return out
}

var _ = fmt.Sprintf // keep import
