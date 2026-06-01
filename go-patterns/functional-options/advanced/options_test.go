package main

import (
	"errors"
	"testing"
	"time"
)

func TestNew_defaults(t *testing.T) {
	s, err := New(":8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Addr() != ":8080" {
		t.Errorf("addr: got %s, want :8080", s.Addr())
	}
	if s.Timeout() != 30*time.Second {
		t.Errorf("timeout: got %s, want 30s", s.Timeout())
	}
}

func TestNew_customOptions(t *testing.T) {
	s, err := New(":9090", WithTimeout(5*time.Second), WithMaxConn(50))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Timeout() != 5*time.Second {
		t.Errorf("timeout: got %s, want 5s", s.Timeout())
	}
	if s.MaxConn() != 50 {
		t.Errorf("maxConn: got %d, want 50", s.MaxConn())
	}
}

func TestNew_validation(t *testing.T) {
	_, err := New("", WithTimeout(-1*time.Second))
	if !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}
