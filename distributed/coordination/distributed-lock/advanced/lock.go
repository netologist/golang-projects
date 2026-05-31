package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrLockHeld is returned when acquisition fails due to an active holder.
var ErrLockHeld = errors.New("lock already held")

// Locker is the distributed lock interface.
type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (token string, err error)
	Unlock(ctx context.Context, key, token string) error
	Extend(ctx context.Context, key, token string, ttl time.Duration) error
}

type lockEntry struct {
	token   string
	expires time.Time
}

// MemoryLocker is an in-process implementation modelling a lock server.
type MemoryLocker struct {
	mu    sync.Mutex
	locks map[string]lockEntry
	seq   uint64
}

// NewMemoryLocker creates an empty locker.
func NewMemoryLocker() *MemoryLocker {
	return &MemoryLocker{locks: map[string]lockEntry{}}
}

// Lock blocks until the lock is acquired or the context ends.
func (l *MemoryLocker) Lock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	for {
		l.mu.Lock()
		entry, exists := l.locks[key]
		if !exists || time.Now().After(entry.expires) {
			l.seq++
			token := fmt.Sprintf("token-%d", l.seq)
			l.locks[key] = lockEntry{token: token, expires: time.Now().Add(ttl)}
			l.mu.Unlock()
			return token, nil
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("lock %s: %w", key, ctx.Err())
		case <-time.After(5 * time.Millisecond):
		}
	}
}

// Unlock releases the lock if the token matches.
func (l *MemoryLocker) Unlock(_ context.Context, key, token string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, exists := l.locks[key]
	if !exists || entry.token != token {
		return fmt.Errorf("unlock %s: invalid or expired token", key)
	}
	delete(l.locks, key)
	return nil
}

// Extend renews the lease if the token matches.
func (l *MemoryLocker) Extend(_ context.Context, key, token string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, exists := l.locks[key]
	if !exists || entry.token != token {
		return fmt.Errorf("extend %s: invalid or expired token", key)
	}
	entry.expires = time.Now().Add(ttl)
	l.locks[key] = entry
	return nil
}
