package main

import "fmt"

func main() {
	s := NewStore()
	_ = s.Create("user:1", "alice")

	_, err := s.Get("user:2")
	fmt.Printf("err=%v status=%d\n", err, HTTPStatus(err))

	err = s.Create("user:1", "bob")
	fmt.Printf("err=%v status=%d\n", err, HTTPStatus(err))
}
