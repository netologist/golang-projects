package main

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupRecorder() (*OrderService, *tracetest.SpanRecorder) {
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	svc := NewOrderService(tp.Tracer("test"))
	return svc, sr
}

func TestPlaceOrder_emitsSpans(t *testing.T) {
	svc, sr := setupRecorder()
	if err := svc.PlaceOrder(context.Background(), "order-1", 100); err != nil {
		t.Fatal(err)
	}
	spans := sr.Ended()
	// Expect at least PlaceOrder + reserveStock + charge.
	if len(spans) < 3 {
		t.Errorf("expected >= 3 spans, got %d", len(spans))
	}

	names := map[string]bool{}
	for _, s := range spans {
		names[s.Name()] = true
	}
	for _, want := range []string{"PlaceOrder", "reserveStock", "charge"} {
		if !names[want] {
			t.Errorf("missing span %q", want)
		}
	}
}

func TestPlaceOrder_recordsError(t *testing.T) {
	svc, sr := setupRecorder()
	err := svc.PlaceOrder(context.Background(), "order-2", -5)
	if err == nil {
		t.Fatal("expected error for invalid amount")
	}
	// The PlaceOrder span should record the error status.
	for _, s := range sr.Ended() {
		if s.Name() == "PlaceOrder" && s.Status().Code.String() != "Error" {
			t.Errorf("PlaceOrder span status: got %s, want Error", s.Status().Code)
		}
	}
}
