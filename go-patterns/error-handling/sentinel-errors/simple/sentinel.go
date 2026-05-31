package main

import (
	"errors"
	"fmt"
)

// ErrNotFound indicates a missing key.
var ErrNotFound = errors.New("not found")

type store struct{ data map[string]string }

func (s *store) get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("get %q: %w", key, ErrNotFound)
	}
	return v, nil
}
