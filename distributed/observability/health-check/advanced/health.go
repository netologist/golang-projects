package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status is an overall health verdict.
type Status string

// Health statuses.
const (
	StatusOK   Status = "ok"
	StatusDown Status = "down"
)

// CheckFn reports the health of one dependency.
type CheckFn func(ctx context.Context) error

// Server aggregates liveness and readiness checks.
type Server struct {
	mu    sync.RWMutex
	live  map[string]CheckFn
	ready map[string]CheckFn
}

// New creates an empty health server.
func New() *Server {
	return &Server{live: map[string]CheckFn{}, ready: map[string]CheckFn{}}
}

// AddLiveness registers a liveness check.
func (s *Server) AddLiveness(name string, fn CheckFn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.live[name] = fn
}

// AddReadiness registers a readiness check.
func (s *Server) AddReadiness(name string, fn CheckFn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready[name] = fn
}

type response struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// Handler returns an http.Handler serving /healthz/live and /healthz/ready.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/live", s.handleChecks(s.live))
	mux.HandleFunc("/healthz/ready", s.handleChecks(s.ready))
	return mux
}

func (s *Server) handleChecks(checks map[string]CheckFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		s.mu.RLock()
		fns := make(map[string]CheckFn, len(checks))
		for k, v := range checks {
			fns[k] = v
		}
		s.mu.RUnlock()

		resp := response{Status: string(StatusOK), Checks: map[string]string{}}
		for name, fn := range fns {
			if err := fn(ctx); err != nil {
				resp.Checks[name] = err.Error()
				resp.Status = string(StatusDown)
			} else {
				resp.Checks[name] = "ok"
			}
		}

		code := http.StatusOK
		if resp.Status != string(StatusOK) {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
