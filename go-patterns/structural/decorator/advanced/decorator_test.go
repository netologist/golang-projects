package main

import (
	"context"
	"strings"
	"testing"
)

type mockStore struct {
	data     map[string]string
	getCalls int
}

func (m *mockStore) Get(_ context.Context, key string) (string, error) {
	m.getCalls++
	v, ok := m.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (m *mockStore) Set(_ context.Context, key, val string) error {
	m.data[key] = val
	return nil
}

func TestLoggingStore_logsGet(t *testing.T) {
	var buf strings.Builder
	inner := &mockStore{data: map[string]string{"k": "v"}}
	s := NewLoggingStore(inner, &buf)

	val, err := s.Get(context.Background(), "k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "v" {
		t.Errorf("got %s, want v", val)
	}
	if !strings.Contains(buf.String(), "k") {
		t.Errorf("expected log to contain key, got: %s", buf.String())
	}
}

func TestCachingStore_cachesResult(t *testing.T) {
	inner := &mockStore{data: map[string]string{"k": "v"}}
	s := NewCachingStore(inner)

	s.Get(context.Background(), "k")
	s.Get(context.Background(), "k")

	if inner.getCalls != 1 {
		t.Errorf("expected 1 inner call, got %d", inner.getCalls)
	}
}
