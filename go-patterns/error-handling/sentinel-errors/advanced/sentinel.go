package main

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for the store API.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
)

// HTTPStatus maps a sentinel error to an HTTP status code.
func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Store is an in-memory store returning wrapped sentinel errors.
type Store struct{ data map[string]string }

// NewStore creates an empty store.
func NewStore() *Store { return &Store{data: map[string]string{}} }

// Get returns the value or a wrapped ErrNotFound.
func (s *Store) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("store.Get %q: %w", key, ErrNotFound)
	}
	return v, nil
}

// Create inserts a value or returns a wrapped ErrConflict.
func (s *Store) Create(key, val string) error {
	if _, exists := s.data[key]; exists {
		return fmt.Errorf("store.Create %q: %w", key, ErrConflict)
	}
	s.data[key] = val
	return nil
}
