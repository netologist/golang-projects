package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := NewBus()
	done := make(chan struct{}, 2)
	bus.Subscribe(ctx, "order.placed", func(e Event) {
		fmt.Printf("fulfilment: ship order %s\n", e.Payload)
		done <- struct{}{}
	}, 4)
	bus.Subscribe(ctx, "order.placed", func(e Event) {
		fmt.Printf("billing: charge order %s\n", e.Payload)
		done <- struct{}{}
	}, 4)

	bus.Publish(ctx, Event{Topic: "order.placed", Payload: "A-100"})
	<-done
	<-done
	time.Sleep(5 * time.Millisecond)
}
