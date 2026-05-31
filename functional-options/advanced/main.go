package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	s, err := New(":8080", WithTimeout(10*time.Second), WithMaxConn(200))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server ready: addr=%s timeout=%s maxConn=%d\n",
		s.Addr(), s.Timeout(), s.MaxConn())

	_, err = New("", WithTimeout(-1*time.Second))
	fmt.Printf("Expected error: %v\n", err)
}
