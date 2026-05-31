package main

import "context"

type ctxKey string

const requestIDKey ctxKey = "request-id"

// WithRequestID attaches a request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID extracts the request ID, if present.
func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}
