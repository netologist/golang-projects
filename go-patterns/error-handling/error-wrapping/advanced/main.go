package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	cause := errors.New("dial tcp: connection refused")
	err := New(http.StatusBadGateway, "UPSTREAM", "payment service unreachable", cause)

	fmt.Println("error:", err)
	fmt.Println("stack:", err.Stack())

	data, _ := err.MarshalJSON()
	fmt.Println("json:", string(data))
}
