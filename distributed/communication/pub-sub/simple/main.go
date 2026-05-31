package main

import (
	"fmt"
	"sync"
)

type bus struct {
	mu   sync.RWMutex
	subs map[string][]chan string
}

func newBus() *bus { return &bus{subs: map[string][]chan string{}} }

func (b *bus) subscribe(topic string) <-chan string {
	ch := make(chan string, 4)
	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], ch)
	b.mu.Unlock()
	return ch
}

func (b *bus) publish(topic, msg string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[topic] {
		ch <- msg
	}
}

func main() {
	b := newBus()
	s1 := b.subscribe("news")
	s2 := b.subscribe("news")

	b.publish("news", "hello")
	fmt.Println("s1:", <-s1)
	fmt.Println("s2:", <-s2)
}
