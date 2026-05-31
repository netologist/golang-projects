package main

import "fmt"

type event struct {
	id        int
	body      string
	published bool
}

type store struct {
	records []string
	outbox  []*event
	nextID  int
}

// save writes a record and its event atomically (same slice append here).
func (s *store) save(record, eventBody string) {
	s.records = append(s.records, record)
	s.nextID++
	s.outbox = append(s.outbox, &event{id: s.nextID, body: eventBody})
}

// relay publishes all pending events.
func (s *store) relay(publish func(string)) {
	for _, e := range s.outbox {
		if !e.published {
			publish(e.body)
			e.published = true
		}
	}
}

func main() {
	s := &store{}
	s.save("order:1", "OrderCreated{1}")
	s.save("order:2", "OrderCreated{2}")

	s.relay(func(body string) { fmt.Println("published:", body) })
	s.relay(func(body string) { fmt.Println("published again?:", body) }) // nothing new
}
