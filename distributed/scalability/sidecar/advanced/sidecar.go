package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Sidecar is an in-process reverse proxy adding cross-cutting behaviour to an
// upstream handler: header injection, request counting, and retry on 503.
type Sidecar struct {
	upstream   http.Handler
	maxRetries int
	requests   atomic.Int64
	retries    atomic.Int64
}

// New wraps upstream with up to maxRetries retries on 503.
func New(upstream http.Handler, maxRetries int) *Sidecar {
	return &Sidecar{upstream: upstream, maxRetries: maxRetries}
}

// Requests returns the total requests handled.
func (s *Sidecar) Requests() int64 { return s.requests.Load() }

// Retries returns the total retries performed.
func (s *Sidecar) Retries() int64 { return s.retries.Load() }

type bufferingWriter struct {
	header http.Header
	status int
	body   []byte
}

func (b *bufferingWriter) Header() http.Header { return b.header }
func (b *bufferingWriter) WriteHeader(c int)   { b.status = c }
func (b *bufferingWriter) Write(p []byte) (int, error) {
	b.body = append(b.body, p...)
	return len(p), nil
}

// ServeHTTP injects a header, forwards to upstream, and retries on 503.
func (s *Sidecar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.requests.Add(1)
	r.Header.Set("X-Sidecar", "v1")

	var last *bufferingWriter
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		bw := &bufferingWriter{header: http.Header{}, status: http.StatusOK}
		s.upstream.ServeHTTP(bw, r)
		last = bw
		if bw.status != http.StatusServiceUnavailable {
			break
		}
		if attempt < s.maxRetries {
			s.retries.Add(1)
		}
	}

	for k, vals := range last.header {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(last.status)
	_, _ = w.Write(last.body)
}

// LoggingUpstream is a tiny upstream that fails the first n calls with 503.
type LoggingUpstream struct {
	failFirst int
	calls     atomic.Int64
}

// NewLoggingUpstream fails the first failFirst calls with 503, then succeeds.
func NewLoggingUpstream(failFirst int) *LoggingUpstream {
	return &LoggingUpstream{failFirst: failFirst}
}

func (u *LoggingUpstream) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	n := u.calls.Add(1)
	if int(n) <= u.failFirst {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "unavailable")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
