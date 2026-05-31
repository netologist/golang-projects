package main

import "fmt"

func main() {
	c, err := New("json")
	if err != nil {
		panic(err)
	}
	data, _ := c.Encode(map[string]int{"a": 1})
	fmt.Printf("Encoded: %s\n", data)

	_, err = New("msgpack")
	fmt.Printf("Unknown codec: %v\n", err)
}
