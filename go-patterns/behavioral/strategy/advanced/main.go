package main

import (
	"fmt"
	"strings"
)

func main() {
	payloads := map[string][]byte{
		"tiny":  []byte("ok"),
		"large": []byte(strings.Repeat("go ", 5000)),
	}
	for name, p := range payloads {
		c := SelectStrategy(len(p))
		out, err := c(p)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("%s: %d bytes -> %d bytes\n", name, len(p), len(out))
	}
}
