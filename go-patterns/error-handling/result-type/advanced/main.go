package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parse(s string) Result[int] {
	n, err := strconv.Atoi(s)
	if err != nil {
		return Err[int](fmt.Errorf("parse %q: %w", s, err))
	}
	return Ok(n)
}

func validate(n int) Result[int] {
	if n < 0 {
		return Err[int](errors.New("must be non-negative"))
	}
	return Ok(n)
}

func main() {
	good := parse("42").FlatMap(validate).Map(func(n int) int { return n * 2 })
	v, err := good.Unwrap()
	fmt.Printf("good: val=%d err=%v\n", v, err)

	bad := parse("oops").FlatMap(validate)
	_, err = bad.Unwrap()
	fmt.Printf("bad: err=%v\n", err)
}
