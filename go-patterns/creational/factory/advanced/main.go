package main

import "fmt"

func main() {
	c, _ := New("json")
	data, _ := c.Encode(map[string]any{"pattern": "factory", "version": 2})
	fmt.Printf("[%s] %s\n", c.Name(), data)
}
