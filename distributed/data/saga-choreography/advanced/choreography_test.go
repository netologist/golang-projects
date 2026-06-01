package main

import (
	"testing"
)

func buildSagaHappyPath(bus *EventBus, saga *ChoreographySaga) {
	bus.Subscribe("OrderCreated", func(e Event) []Event {
		saga.Record("order", true, "")
		return []Event{{Topic: "PaymentRequested", Payload: e.Payload}}
	})
	bus.Subscribe("PaymentRequested", func(e Event) []Event {
		saga.Record("payment", true, "")
		return []Event{{Topic: "StockReserved", Payload: e.Payload}}
	})
	bus.Subscribe("StockReserved", func(e Event) []Event {
		saga.Record("inventory", true, "")
		return nil
	})
}

func TestChoreography_HappyPath(t *testing.T) {
	bus := NewBus()
	saga := NewChoreographySaga(bus)
	buildSagaHappyPath(bus, saga)

	bus.Publish(Event{Topic: "OrderCreated", Payload: map[string]any{"id": "1"}})

	log := saga.Log()
	if len(log) != 3 {
		t.Fatalf("want 3 steps, got %d", len(log))
	}
	for _, r := range log {
		if !r.Success {
			t.Fatalf("step %s failed", r.Step)
		}
	}
}

func TestChoreography_CompensationOnFailure(t *testing.T) {
	bus := NewBus()
	saga := NewChoreographySaga(bus)

	bus.Subscribe("OrderCreated", func(e Event) []Event {
		saga.Record("order", true, "")
		return []Event{{Topic: "PaymentRequested", Payload: e.Payload}}
	})
	bus.Subscribe("PaymentRequested", func(e Event) []Event {
		saga.Record("payment", true, "")
		// fail on inventory
		return []Event{{Topic: "InventoryFailed", Payload: e.Payload}}
	})
	bus.Subscribe("InventoryFailed", func(e Event) []Event {
		saga.Record("inventory", false, "out of stock")
		return []Event{{Topic: "PaymentRefunded", Payload: e.Payload}}
	})
	bus.Subscribe("PaymentRefunded", func(e Event) []Event {
		saga.Record("compensate-payment", true, "")
		return nil
	})

	bus.Publish(Event{Topic: "OrderCreated", Payload: map[string]any{"id": "99"}})

	log := saga.Log()
	failFound := false
	compFound := false
	for _, r := range log {
		if r.Step == "inventory" && !r.Success {
			failFound = true
		}
		if r.Step == "compensate-payment" {
			compFound = true
		}
	}
	if !failFound {
		t.Fatal("expected inventory failure")
	}
	if !compFound {
		t.Fatal("expected compensation step")
	}
}
