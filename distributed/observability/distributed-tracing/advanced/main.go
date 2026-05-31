package main

import (
	"context"
	"fmt"
)

func main() {
	tracer, tp, shutdown := InitTracer("order-service")
	defer func() { _ = shutdown(context.Background()) }()
	_ = tp

	svc := NewOrderService(tracer)
	if err := svc.PlaceOrder(context.Background(), "order-100", 250); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("order placed; spans emitted (PlaceOrder -> reserveStock, charge)")
}
