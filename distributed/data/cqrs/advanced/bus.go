package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrHandlerNotFound is returned when no handler is registered.
var ErrHandlerNotFound = errors.New("handler not found")

// Command is a write intent.
type Command interface{ CommandName() string }

// Query is a read intent.
type Query interface{ QueryName() string }

// CommandHandler handles a command.
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// QueryHandler handles a query.
type QueryHandler interface {
	Handle(ctx context.Context, q Query) (any, error)
}

// CommandBus dispatches commands to handlers by name.
type CommandBus struct {
	mu       sync.RWMutex
	handlers map[string]CommandHandler
}

// NewCommandBus creates an empty command bus.
func NewCommandBus() *CommandBus { return &CommandBus{handlers: map[string]CommandHandler{}} }

// Register binds a handler to a command name.
func (b *CommandBus) Register(name string, h CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
}

// Dispatch routes cmd to its handler.
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	b.mu.RLock()
	h, ok := b.handlers[cmd.CommandName()]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("command %q: %w", cmd.CommandName(), ErrHandlerNotFound)
	}
	return h.Handle(ctx, cmd)
}

// QueryBus dispatches queries to handlers by name.
type QueryBus struct {
	mu       sync.RWMutex
	handlers map[string]QueryHandler
}

// NewQueryBus creates an empty query bus.
func NewQueryBus() *QueryBus { return &QueryBus{handlers: map[string]QueryHandler{}} }

// Register binds a handler to a query name.
func (b *QueryBus) Register(name string, h QueryHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
}

// Ask routes q to its handler and returns the result.
func (b *QueryBus) Ask(ctx context.Context, q Query) (any, error) {
	b.mu.RLock()
	h, ok := b.handlers[q.QueryName()]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("query %q: %w", q.QueryName(), ErrHandlerNotFound)
	}
	return h.Handle(ctx, q)
}
