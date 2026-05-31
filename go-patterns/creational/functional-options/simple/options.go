package main

import "time"

// Server is the type being configured.
type Server struct {
	addr    string
	timeout time.Duration
	maxConn int
}

// Option configures a Server.
type Option func(*Server)

// WithTimeout sets the read/write timeout.
func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

// WithMaxConn sets the maximum number of concurrent connections.
func WithMaxConn(n int) Option {
	return func(s *Server) { s.maxConn = n }
}

// New creates a Server with defaults, then applies opts.
func New(addr string, opts ...Option) *Server {
	s := &Server{
		addr:    addr,
		timeout: 30 * time.Second,
		maxConn: 100,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}
