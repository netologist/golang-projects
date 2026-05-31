package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAppError_unwrap(t *testing.T) {
	cause := errors.New("db connection refused")
	e := New(http.StatusInternalServerError, "DB_ERR", "database error", cause)
	if !errors.Is(e, cause) {
		t.Error("expected errors.Is to find cause via Unwrap")
	}
}

func TestAppError_json(t *testing.T) {
	e := New(http.StatusNotFound, "NOT_FOUND", "resource not found", nil)
	data, err := e.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "NOT_FOUND") || !strings.Contains(s, "404") {
		t.Errorf("unexpected JSON: %s", s)
	}
}

func TestAppError_as(t *testing.T) {
	wrapped := New(http.StatusConflict, "CONFLICT", "already exists", nil)
	var target *AppError
	if !errors.As(error(wrapped), &target) {
		t.Fatal("expected errors.As to match *AppError")
	}
	if target.Status != http.StatusConflict {
		t.Errorf("got status %d, want 409", target.Status)
	}
}
