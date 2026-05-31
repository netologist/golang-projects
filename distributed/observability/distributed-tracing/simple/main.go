package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	defer func() { _ = tp.Shutdown(context.Background()) }()

	tracer := otel.Tracer("demo")
	ctx, parent := tracer.Start(context.Background(), "handle-request")
	defer parent.End()

	_, child := tracer.Start(ctx, "db-query")
	time.Sleep(5 * time.Millisecond)
	child.End()

	fmt.Println("created parent span 'handle-request' with child span 'db-query'")
}
