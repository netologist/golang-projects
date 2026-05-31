package main

import (
	"errors"
	"net/http"
	"testing"
)

func TestStore_Get_notFound(t *testing.T) {
	s := NewStore()
	_, err := s.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if got := HTTPStatus(err); got != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", got)
	}
}

func TestStore_Create_conflict(t *testing.T) {
	s := NewStore()
	if err := s.Create("k", "v"); err != nil {
		t.Fatal(err)
	}
	err := s.Create("k", "v2")
	if !errors.Is(err, ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
	if got := HTTPStatus(err); got != http.StatusConflict {
		t.Errorf("status: got %d, want 409", got)
	}
}
