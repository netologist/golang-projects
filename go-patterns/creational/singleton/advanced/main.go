package main

import "fmt"

func main() {
	db1 := GetDB()
	db2 := GetDB()
	fmt.Println("same instance:", db1 == db2)
	fmt.Println("ping:", db1.Ping())
}
