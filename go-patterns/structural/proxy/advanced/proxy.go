package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrUnauthorized is returned when an auth token is missing or invalid.
var ErrUnauthorized = errors.New("unauthorized")

// Loader is the subject interface.
type Loader interface {
	Load(ctx context.Context, url string) ([]byte, error)
}

type entry struct {
	data    []byte
	expires time.Time
}

// CachingProxy caches results with a TTL.
type CachingProxy struct {
	inner Loader
	mu    sync.RWMutex
	cache map[string]entry
	ttl   time.Duration
}

// NewCachingProxy wraps inner with a TTL cache.
func NewCachingProxy(inner Loader, ttl time.Duration) *CachingProxy {
	return &CachingProxy{inner: inner, cache: map[string]entry{}, ttl: ttl}
}

func (c *CachingProxy) Load(ctx context.Context, url string) ([]byte, error) {
	c.mu.RLock()
	if e, ok := c.cache[url]; ok && time.Now().Before(e.expires) {
		c.mu.RUnlock()
		return e.data, nil
	}
	c.mu.RUnlock()

	data, err := c.inner.Load(ctx, url)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.cache[url] = entry{data: data, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return data, nil
}

type ctxKey string

// TokenKey is the context key used to carry the auth token.
const TokenKey ctxKey = "auth-token"

// AuthProxy blocks requests without a valid token in the context.
type AuthProxy struct {
	inner      Loader
	validToken string
}

// NewAuthProxy wraps inner, requiring validToken in the context.
func NewAuthProxy(inner Loader, validToken string) *AuthProxy {
	return &AuthProxy{inner: inner, validToken: validToken}
}

func (a *AuthProxy) Load(ctx context.Context, url string) ([]byte, error) {
	tok, _ := ctx.Value(TokenKey).(string)
	if tok != a.validToken {
		return nil, fmt.Errorf("load %s: %w", url, ErrUnauthorized)
	}
	return a.inner.Load(ctx, url)
}
