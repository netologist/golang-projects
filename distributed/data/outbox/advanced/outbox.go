package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// OutboxEvent is a pending message awaiting publication.
type OutboxEvent struct {
	ID          string
	Aggregate   string
	Payload     []byte
	CreatedAt   time.Time
	PublishedAt *time.Time
}

// Store models a transactional outbox over an in-memory dataset.
type Store struct {
	mu      sync.Mutex
	records map[string][]byte
	outbox  map[string]*OutboxEvent
	seq     int
}

// NewStore creates an empty store.
func NewStore() *Store {
	return &Store{records: map[string][]byte{}, outbox: map[string]*OutboxEvent{}}
}

// SaveWithOutbox writes a record and its event atomically.
func (s *Store) SaveWithOutbox(_ context.Context, key string, record []byte, aggregate string, payload []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[key] = record
	s.seq++
	id := fmt.Sprintf("evt-%d", s.seq)
	s.outbox[id] = &OutboxEvent{
		ID:        id,
		Aggregate: aggregate,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
	return id, nil
}

// PendingEvents returns up to limit unpublished events.
func (s *Store) PendingEvents(_ context.Context, limit int) ([]*OutboxEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []*OutboxEvent
	for _, e := range s.outbox {
		if e.PublishedAt == nil {
			out = append(out, e)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

// MarkPublished records that an event was delivered.
func (s *Store) MarkPublished(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.outbox[id]
	if !ok {
		return fmt.Errorf("outbox: event %q not found", id)
	}
	now := time.Now()
	e.PublishedAt = &now
	return nil
}

// Relay publishes pending events via publish and marks them published.
// It returns the number of events delivered.
func (s *Store) Relay(ctx context.Context, batch int, publish func(*OutboxEvent) error) (int, error) {
	events, err := s.PendingEvents(ctx, batch)
	if err != nil {
		return 0, err
	}
	delivered := 0
	for _, e := range events {
		if err := publish(e); err != nil {
			return delivered, fmt.Errorf("relay publish %s: %w", e.ID, err)
		}
		if err := s.MarkPublished(ctx, e.ID); err != nil {
			return delivered, err
		}
		delivered++
	}
	return delivered, nil
}
