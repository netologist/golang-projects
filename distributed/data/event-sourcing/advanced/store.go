package main

import (
	"fmt"
	"sync"
	"time"
)

// Event is a domain event.
type Event interface {
	EventType() string
	OccurredAt() time.Time
}

// BaseEvent provides common event fields.
type BaseEvent struct {
	Type string
	Time time.Time
}

// EventType returns the event type name.
func (b BaseEvent) EventType() string { return b.Type }

// OccurredAt returns when the event happened.
func (b BaseEvent) OccurredAt() time.Time { return b.Time }

// Store is an append-only event store with optimistic concurrency.
type Store struct {
	mu      sync.RWMutex
	streams map[string][]Event
}

// NewStore creates an empty store.
func NewStore() *Store { return &Store{streams: map[string][]Event{}} }

// Append adds events to a stream. expectedVersion of -1 skips the check.
func (s *Store) Append(streamID string, events []Event, expectedVersion int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.streams[streamID]
	if expectedVersion >= 0 && len(current) != expectedVersion {
		return fmt.Errorf("optimistic concurrency: stream %q at version %d, expected %d",
			streamID, len(current), expectedVersion)
	}
	s.streams[streamID] = append(current, events...)
	return nil
}

// Load returns a copy of a stream's events.
func (s *Store) Load(streamID string) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}
	out := make([]Event, len(events))
	copy(out, events)
	return out, nil
}

// --- Bank account aggregate ---

// Deposited is emitted on a deposit.
type Deposited struct {
	BaseEvent
	Amount int
}

// Withdrawn is emitted on a withdrawal.
type Withdrawn struct {
	BaseEvent
	Amount int
}

// Account folds events into a balance and version.
type Account struct {
	Balance int
	Version int
}

// Apply mutates the account for one event.
func (a *Account) Apply(e Event) {
	switch ev := e.(type) {
	case Deposited:
		a.Balance += ev.Amount
	case Withdrawn:
		a.Balance -= ev.Amount
	}
	a.Version++
}

// Rebuild folds a slice of events into an Account.
func Rebuild(events []Event) *Account {
	a := &Account{}
	for _, e := range events {
		a.Apply(e)
	}
	return a
}
