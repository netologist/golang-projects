package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// ErrNotFound is returned when a key is absent.
var ErrNotFound = errors.New("not found")

// Store is the interface being decorated.
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string) error
}

// LoggingStore logs every Get and Set call.
type LoggingStore struct {
	inner  Store
	output io.Writer
}

// NewLoggingStore wraps inner, writing logs to w.
func NewLoggingStore(inner Store, w io.Writer) *LoggingStore {
	return &LoggingStore{inner: inner, output: w}
}

func (l *LoggingStore) Get(ctx context.Context, key string) (string, error) {
	val, err := l.inner.Get(ctx, key)
	fmt.Fprintf(l.output, "Get key=%s val=%s err=%v\n", key, val, err)
	return val, err
}

func (l *LoggingStore) Set(ctx context.Context, key, val string) error {
	err := l.inner.Set(ctx, key, val)
	fmt.Fprintf(l.output, "Set key=%s val=%s err=%v\n", key, val, err)
	return err
}

// CachingStore caches Get results in memory.
type CachingStore struct {
	inner Store
	mu    sync.RWMutex
	cache map[string]string
}

// NewCachingStore wraps inner with an in-memory read cache.
func NewCachingStore(inner Store) *CachingStore {
	return &CachingStore{inner: inner, cache: map[string]string{}}
}

func (c *CachingStore) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	if v, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	val, err := c.inner.Get(ctx, key)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.cache[key] = val
	c.mu.Unlock()
	return val, nil
}

func (c *CachingStore) Set(ctx context.Context, key, val string) error {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
	return c.inner.Set(ctx, key, val)
}
