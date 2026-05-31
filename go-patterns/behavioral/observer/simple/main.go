package main

import "fmt"

func main() {
	bus := NewBus()
	bus.Subscribe("user.created", func(e Event) {
		fmt.Printf("email service: welcome %s\n", e.Payload)
	})
	bus.Subscribe("user.created", func(e Event) {
		fmt.Printf("audit log: created %s\n", e.Payload)
	})
	bus.Publish(Event{Topic: "user.created", Payload: "alice"})
}
