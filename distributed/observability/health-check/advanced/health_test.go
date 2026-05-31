package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_allHealthy(t *testing.T) {
	s := New()
	s.AddReadiness("db", func(_ context.Context) error { return nil })
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %d, want 200", rr.Code)
	}
}

func TestServer_degradedReturns503(t *testing.T) {
	s := New()
	s.AddReadiness("db", func(_ context.Context) error { return errors.New("connection refused") })
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("got %d, want 503", rr.Code)
	}
}

func TestServer_livenessIndependent(t *testing.T) {
	s := New()
	s.AddLiveness("self", func(_ context.Context) error { return nil })
	s.AddReadiness("db", func(_ context.Context) error { return errors.New("down") })

	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("liveness should be 200 even when readiness is down, got %d", rr.Code)
	}
}
