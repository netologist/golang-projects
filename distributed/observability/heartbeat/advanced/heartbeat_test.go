package main

import (
	"testing"
	"time"
)

func TestMonitor_NodeAlive(t *testing.T) {
	events := make([]string, 0)
	m := NewMonitor(200*time.Millisecond, 2, func(id string, _ NodeStatus, new NodeStatus) {
		events = append(events, id+":"+new.String())
	})
	m.Register("n1")
	m.Start(50 * time.Millisecond)
	defer m.Stop()

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		m.Beat("n1")
	}
	time.Sleep(100 * time.Millisecond)
	if !m.IsAlive("n1") {
		t.Fatal("node should be alive")
	}
}

func TestMonitor_NodeTimeout(t *testing.T) {
	var lastStatus NodeStatus
	m := NewMonitor(100*time.Millisecond, 2, func(id string, _ NodeStatus, new NodeStatus) {
		lastStatus = new
	})
	m.Register("n1")
	m.Start(50 * time.Millisecond)
	defer m.Stop()

	// no heartbeats → should go DOWN
	time.Sleep(400 * time.Millisecond)
	if m.IsAlive("n1") {
		t.Fatal("node should be down")
	}
	if lastStatus != StatusDown {
		t.Fatalf("want DOWN, got %s", lastStatus)
	}
}

func TestMonitor_Recovery(t *testing.T) {
	m := NewMonitor(100*time.Millisecond, 2, nil)
	m.Register("n1")
	m.Start(50 * time.Millisecond)
	defer m.Stop()

	time.Sleep(400 * time.Millisecond) // go down
	if m.IsAlive("n1") {
		t.Fatal("should be down before recovery")
	}
	m.Beat("n1") // recover
	if !m.IsAlive("n1") {
		t.Fatal("should be alive after beat")
	}
}
