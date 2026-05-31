package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// InitTracer sets up an OTel tracer provider with an in-memory span recorder for
// inspection in tests. It returns the tracer, the recorder, and a shutdown func.
func InitTracer(service string) (oteltrace.Tracer, *trace.TracerProvider, func(context.Context) error) {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	return tp.Tracer(service), tp, tp.Shutdown
}

// OrderService is a small example service that emits spans.
type OrderService struct {
	tracer oteltrace.Tracer
}

// NewOrderService creates a service using the given tracer.
func NewOrderService(tracer oteltrace.Tracer) *OrderService {
	return &OrderService{tracer: tracer}
}

// PlaceOrder creates a parent span and child spans for sub-operations.
func (s *OrderService) PlaceOrder(ctx context.Context, orderID string, amount int) error {
	ctx, span := s.tracer.Start(ctx, "PlaceOrder")
	defer span.End()
	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.Int("order.amount", amount),
	)

	if err := s.reserveStock(ctx, orderID); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	if err := s.charge(ctx, amount); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	span.AddEvent("order placed")
	return nil
}

func (s *OrderService) reserveStock(ctx context.Context, orderID string) error {
	_, span := s.tracer.Start(ctx, "reserveStock")
	defer span.End()
	span.SetAttributes(attribute.String("order.id", orderID))
	time.Sleep(2 * time.Millisecond)
	return nil
}

func (s *OrderService) charge(ctx context.Context, amount int) error {
	_, span := s.tracer.Start(ctx, "charge")
	defer span.End()
	if amount <= 0 {
		return fmt.Errorf("invalid amount %d", amount)
	}
	time.Sleep(2 * time.Millisecond)
	return nil
}
