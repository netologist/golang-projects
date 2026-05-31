package main

// Event is a topic/payload pair.
type Event struct{ Topic, Payload string }

// Handler reacts to an Event.
type Handler func(Event)

// Bus is a synchronous in-process event bus.
type Bus struct {
	subs map[string][]Handler
}

// NewBus creates an empty bus.
func NewBus() *Bus { return &Bus{subs: map[string][]Handler{}} }

// Subscribe registers h for the topic.
func (b *Bus) Subscribe(topic string, h Handler) {
	b.subs[topic] = append(b.subs[topic], h)
}

// Publish delivers e to all subscribers of its topic synchronously.
func (b *Bus) Publish(e Event) {
	for _, h := range b.subs[e.Topic] {
		h(e)
	}
}
