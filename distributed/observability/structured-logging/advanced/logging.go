package main

import (
	"context"
	"log/slog"
)

type ctxKey string

const logAttrsKey ctxKey = "log-attrs"

// WithAttrs returns a context carrying extra log attributes.
func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing, _ := ctx.Value(logAttrsKey).([]slog.Attr)
	merged := make([]slog.Attr, 0, len(existing)+len(attrs))
	merged = append(merged, existing...)
	merged = append(merged, attrs...)
	return context.WithValue(ctx, logAttrsKey, merged)
}

// contextHandler injects context attributes into every record.
type contextHandler struct {
	slog.Handler
}

// Handle adds context attributes before delegating.
func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, r)
}

// NewLogger wraps a base handler with context enrichment.
func NewLogger(base slog.Handler) *slog.Logger {
	return slog.New(contextHandler{Handler: base})
}
