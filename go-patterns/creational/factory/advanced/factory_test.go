package main

import (
	"errors"
	"testing"
)

func TestNew_knownCodec(t *testing.T) {
	c, err := New("json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name() != "json" {
		t.Errorf("got %s, want json", c.Name())
	}
}

func TestNew_roundtrip(t *testing.T) {
	c, _ := New("json")
	type payload struct{ X int }
	data, err := c.Encode(payload{X: 42})
	if err != nil {
		t.Fatal(err)
	}
	var out payload
	if err := c.Decode(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.X != 42 {
		t.Errorf("got %d, want 42", out.X)
	}
}

func TestNew_unknownCodec(t *testing.T) {
	_, err := New("nope")
	if !errors.Is(err, ErrUnknownCodec) {
		t.Errorf("expected ErrUnknownCodec, got %v", err)
	}
}
