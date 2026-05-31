package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSidecar_injectsHeader(t *testing.T) {
	var seen string
	upstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("X-Sidecar")
		w.WriteHeader(http.StatusOK)
	})
	sc := New(upstream, 0)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	sc.ServeHTTP(httptest.NewRecorder(), req)
	if seen != "v1" {
		t.Errorf("expected injected header v1, got %q", seen)
	}
}

func TestSidecar_retriesOn503(t *testing.T) {
	upstream := NewLoggingUpstream(2) // fail first 2 calls
	sc := New(upstream, 3)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	sc.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 after retries, got %d", rr.Code)
	}
	if sc.Retries() != 2 {
		t.Errorf("expected 2 retries, got %d", sc.Retries())
	}
}
