package main

import (
	"fmt"
	"time"
)

func main() {
	s1 := New(":8080")
	fmt.Printf("s1: addr=%s timeout=%s maxConn=%d\n", s1.addr, s1.timeout, s1.maxConn)

	s2 := New(":9090", WithTimeout(5*time.Second), WithMaxConn(50))
	fmt.Printf("s2: addr=%s timeout=%s maxConn=%d\n", s2.addr, s2.timeout, s2.maxConn)
}
