package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockLoader struct{ calls int }

func (m *mockLoader) Load(_ context.Context, _ string) ([]byte, error) {
	m.calls++
	return []byte("data"), nil
}

func TestCachingProxy_cacheHit(t *testing.T) {
	inner := &mockLoader{}
	p := NewCachingProxy(inner, time.Minute)
	p.Load(context.Background(), "url")
	p.Load(context.Background(), "url")
	if inner.calls != 1 {
		t.Errorf("expected 1 inner call, got %d", inner.calls)
	}
}

func TestAuthProxy_blocked(t *testing.T) {
	inner := &mockLoader{}
	p := NewAuthProxy(inner, "secret")
	_, err := p.Load(context.Background(), "url")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthProxy_allowed(t *testing.T) {
	inner := &mockLoader{}
	p := NewAuthProxy(inner, "secret")
	ctx := context.WithValue(context.Background(), TokenKey, "secret")
	if _, err := p.Load(ctx, "url"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
