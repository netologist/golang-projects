package main

import (
	"fmt"
	"log"
)

func main() {
	q, err := NewQueryBuilder("users").
		Select("id", "name", "email").
		Where("active = true").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
}
