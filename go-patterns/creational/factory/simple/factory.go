package main

import (
	"encoding/json"
	"fmt"
)

// Codec serialises and deserialises values.
type Codec interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
	Name() string
}

var registry = map[string]func() Codec{}

// Register makes a codec available by name.
func Register(name string, fn func() Codec) { registry[name] = fn }

// New returns a new instance of the named codec.
func New(name string) (Codec, error) {
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("codec %q not registered", name)
	}
	return fn(), nil
}

type jsonCodec struct{}

func (c *jsonCodec) Name() string                    { return "json" }
func (c *jsonCodec) Encode(v any) ([]byte, error)    { return json.Marshal(v) }
func (c *jsonCodec) Decode(data []byte, v any) error { return json.Unmarshal(data, v) }

func init() {
	Register("json", func() Codec { return &jsonCodec{} })
}
