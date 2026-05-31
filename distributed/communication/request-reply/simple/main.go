package main

import (
	"fmt"
	"time"
)

func main() {
	reply := make(chan string, 1)

	// Responder: handles the request asynchronously.
	go func(req string, out chan<- string) {
		time.Sleep(10 * time.Millisecond)
		out <- "echo: " + req
	}("ping", reply)

	select {
	case r := <-reply:
		fmt.Println("got reply:", r)
	case <-time.After(time.Second):
		fmt.Println("timed out")
	}
}
