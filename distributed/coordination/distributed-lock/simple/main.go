package main

import (
	"fmt"
	"sync"
)

type lock struct {
	mu    sync.Mutex
	held  bool
	owner string
}

func (l *lock) acquire(owner string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.held {
		return false
	}
	l.held = true
	l.owner = owner
	return true
}

func (l *lock) release(owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.owner == owner {
		l.held = false
		l.owner = ""
	}
}

func main() {
	l := &lock{}
	fmt.Println("A acquires:", l.acquire("A"))
	fmt.Println("B acquires:", l.acquire("B"))
	l.release("A")
	fmt.Println("B acquires after A release:", l.acquire("B"))
}
