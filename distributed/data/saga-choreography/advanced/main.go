package main

import "fmt"

func main() {
	runSaga := func(orderID string, failInventory bool) {
		bus := NewBus()
		saga := NewChoreographySaga(bus)

		// Order service
		bus.Subscribe("OrderPlaced", func(e Event) []Event {
			fmt.Printf("[order] created %s\n", e.Payload["id"])
			saga.Record("order-created", true, "")
			return []Event{{Topic: "PaymentRequested", Payload: e.Payload}}
		})

		// Payment service
		bus.Subscribe("PaymentRequested", func(e Event) []Event {
			fmt.Printf("[payment] charged for %s\n", e.Payload["id"])
			saga.Record("payment-done", true, "")
			return []Event{{Topic: "InventoryReserve", Payload: e.Payload}}
		})

		// Inventory service
		bus.Subscribe("InventoryReserve", func(e Event) []Event {
			if failInventory {
				fmt.Printf("[inventory] FAILED for %s\n", e.Payload["id"])
				saga.Record("inventory-reserve", false, "out of stock")
				return []Event{{Topic: "PaymentRefund", Payload: e.Payload}}
			}
			fmt.Printf("[inventory] reserved for %s\n", e.Payload["id"])
			saga.Record("inventory-reserve", true, "")
			return nil
		})

		// Compensation
		bus.Subscribe("PaymentRefund", func(e Event) []Event {
			fmt.Printf("[payment] refunded for %s\n", e.Payload["id"])
			saga.Record("payment-refund", true, "")
			return []Event{{Topic: "OrderCancelled", Payload: e.Payload}}
		})
		bus.Subscribe("OrderCancelled", func(e Event) []Event {
			fmt.Printf("[order] cancelled %s\n", e.Payload["id"])
			saga.Record("order-cancelled", true, "")
			return nil
		})

		bus.Publish(Event{Topic: "OrderPlaced", Payload: map[string]any{"id": orderID}})
		fmt.Println("saga log:")
		saga.PrintLog()
		fmt.Println()
	}

	fmt.Println("=== Happy path ===")
	runSaga("order-42", false)

	fmt.Println("=== Failure path ===")
	runSaga("order-99", true)
}
