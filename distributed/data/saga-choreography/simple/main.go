package main

import (
	"fmt"
)

type Event struct {
	Type    string
	Payload map[string]any
}

type Bus struct {
	subs map[string][]func(Event)
}

func (b *Bus) Subscribe(topic string, h func(Event)) {
	if b.subs == nil {
		b.subs = map[string][]func(Event){}
	}
	b.subs[topic] = append(b.subs[topic], h)
}

func (b *Bus) Publish(e Event) {
	for _, h := range b.subs[e.Type] {
		h(e)
	}
}

func main() {
	bus := &Bus{}

	// Order service: creates order, emits OrderCreated
	bus.Subscribe("StartOrder", func(e Event) {
		fmt.Println("[order] order created:", e.Payload["id"])
		bus.Publish(Event{Type: "OrderCreated", Payload: e.Payload})
	})

	// Payment service: processes payment on OrderCreated
	bus.Subscribe("OrderCreated", func(e Event) {
		fmt.Println("[payment] payment processed for order:", e.Payload["id"])
		bus.Publish(Event{Type: "PaymentDone", Payload: e.Payload})
	})

	// Inventory service: reserves stock; fails for order 99
	bus.Subscribe("PaymentDone", func(e Event) {
		if e.Payload["id"] == "99" {
			fmt.Println("[inventory] reservation FAILED, emitting compensate")
			bus.Publish(Event{Type: "CompensatePayment", Payload: e.Payload})
			return
		}
		fmt.Println("[inventory] stock reserved for order:", e.Payload["id"])
		bus.Publish(Event{Type: "SagaDone", Payload: e.Payload})
	})

	// Compensation handlers
	bus.Subscribe("CompensatePayment", func(e Event) {
		fmt.Println("[payment] payment REFUNDED for order:", e.Payload["id"])
		bus.Publish(Event{Type: "CompensateOrder", Payload: e.Payload})
	})
	bus.Subscribe("CompensateOrder", func(e Event) {
		fmt.Println("[order] order CANCELLED for:", e.Payload["id"])
	})

	bus.Subscribe("SagaDone", func(e Event) {
		fmt.Println("[saga] completed successfully for order:", e.Payload["id"])
	})

	fmt.Println("--- Happy path ---")
	bus.Publish(Event{Type: "StartOrder", Payload: map[string]any{"id": "42"}})
	fmt.Println("--- Failure path ---")
	bus.Publish(Event{Type: "StartOrder", Payload: map[string]any{"id": "99"}})
}
