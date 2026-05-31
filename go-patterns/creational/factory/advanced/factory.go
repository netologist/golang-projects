package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// ErrUnknownCodec is returned when the requested codec is not registered.
var ErrUnknownCodec = errors.New("unknown codec")

// Codec serialises and deserialises values.
type Codec interface {
	Name() string
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type registry struct {
	mu     sync.RWMutex
	codecs map[string]func() Codec
}

var global = &registry{codecs: map[string]func() Codec{}}

// Register makes a codec available by name.
func Register(name string, fn func() Codec) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.codecs[name] = fn
}

// New returns a new instance of the named codec.
func New(name string) (Codec, error) {
	global.mu.RLock()
	fn, ok := global.codecs[name]
	global.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("factory: %w: %q", ErrUnknownCodec, name)
	}
	return fn(), nil
}

type jsonCodec struct{}

func (c *jsonCodec) Name() string                 { return "json" }
func (c *jsonCodec) Encode(v any) ([]byte, error) { return json.Marshal(v) }
func (c *jsonCodec) Decode(d []byte, v any) error { return json.Unmarshal(d, v) }

func init() { Register("json", func() Codec { return &jsonCodec{} }) }
