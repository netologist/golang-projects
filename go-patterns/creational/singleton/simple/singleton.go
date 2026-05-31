package main

import (
	"fmt"
	"sync"
)

type logger struct{ prefix string }

func (l *logger) Log(msg string) { fmt.Printf("[%s] %s\n", l.prefix, msg) }

var (
	instance *logger
	once     sync.Once
)

// GetLogger returns the singleton logger, initialised on first call.
func GetLogger() *logger {
	once.Do(func() {
		instance = &logger{prefix: "APP"}
	})
	return instance
}
