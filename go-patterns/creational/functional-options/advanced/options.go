package main

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidConfig is returned when option validation fails.
var ErrInvalidConfig = errors.New("invalid config")

type config struct {
	addr    string
	timeout time.Duration
	maxConn int
}

// Server is the production-grade configured server.
type Server struct {
	cfg config
}

// Option mutates the config and may return a validation error.
type Option func(*config) error

// WithTimeout sets the read/write timeout. Must be > 0.
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return fmt.Errorf("%w: timeout must be > 0", ErrInvalidConfig)
		}
		c.timeout = d
		return nil
	}
}

// WithMaxConn sets the max concurrent connections. Must be > 0.
func WithMaxConn(n int) Option {
	return func(c *config) error {
		if n <= 0 {
			return fmt.Errorf("%w: maxConn must be > 0", ErrInvalidConfig)
		}
		c.maxConn = n
		return nil
	}
}

// New creates a Server, applies options, and validates.
func New(addr string, opts ...Option) (*Server, error) {
	if addr == "" {
		return nil, fmt.Errorf("%w: addr is required", ErrInvalidConfig)
	}
	cfg := config{
		addr:    addr,
		timeout: 30 * time.Second,
		maxConn: 100,
	}
	for _, o := range opts {
		if err := o(&cfg); err != nil {
			return nil, err
		}
	}
	return &Server{cfg: cfg}, nil
}

// Addr returns the configured address.
func (s *Server) Addr() string { return s.cfg.addr }

// Timeout returns the configured timeout.
func (s *Server) Timeout() time.Duration { return s.cfg.timeout }

// MaxConn returns the configured max connections.
func (s *Server) MaxConn() int { return s.cfg.maxConn }
