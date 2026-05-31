package main

import (
	"errors"
	"fmt"
)

func main() {
	err := loadProfile(7)
	fmt.Println("error chain:", err)
	fmt.Println("is db error:", errors.Is(err, errDB))
}
