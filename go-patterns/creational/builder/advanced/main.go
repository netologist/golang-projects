package main

import (
	"fmt"
	"log"
)

func main() {
	q, err := NewQueryBuilder("orders").
		Select("id", "total", "status").
		Where("status = $1", "pending").
		Where("total > $2", 100).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SQL:", q.SQL())
	fmt.Println("Args:", q.Args())
}
