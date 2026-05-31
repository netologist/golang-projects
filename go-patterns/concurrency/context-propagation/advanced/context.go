package main

import "context"

type ctxKey string

const (
	requestIDKey ctxKey = "request-id"
	userIDKey    ctxKey = "user-id"
	tenantIDKey  ctxKey = "tenant-id"
)

// WithRequestID attaches a request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID extracts the request ID.
func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

// WithUserID attaches a user ID.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserID extracts the user ID.
func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

// WithTenantID attaches a tenant ID.
func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tenantIDKey, id)
}

// TenantID extracts the tenant ID.
func TenantID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(tenantIDKey).(string)
	return v, ok
}
