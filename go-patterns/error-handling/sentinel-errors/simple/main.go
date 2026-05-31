package main

import (
	"errors"
	"fmt"
)

func main() {
	s := &store{data: map[string]string{"a": "1"}}

	_, err := s.get("missing")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("handled not-found:", err)
	}
}
