package main

import (
	"context"
	"testing"
)

func TestRequestID_roundtrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")
	id, ok := RequestID(ctx)
	if !ok {
		t.Fatal("expected ok")
	}
	if id != "abc-123" {
		t.Errorf("got %s, want abc-123", id)
	}
}

func TestRequestID_missing(t *testing.T) {
	if _, ok := RequestID(context.Background()); ok {
		t.Error("expected missing")
	}
}

func TestLayeredValues(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "r1")
	ctx = WithUserID(ctx, "u1")
	ctx = WithTenantID(ctx, "t1")

	if v, _ := RequestID(ctx); v != "r1" {
		t.Errorf("request id: got %s", v)
	}
	if v, _ := UserID(ctx); v != "u1" {
		t.Errorf("user id: got %s", v)
	}
	if v, _ := TenantID(ctx); v != "t1" {
		t.Errorf("tenant id: got %s", v)
	}
}
