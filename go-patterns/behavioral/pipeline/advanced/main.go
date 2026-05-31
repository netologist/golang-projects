package main

import (
	"context"
	"fmt"
	"strings"
)

func main() {
	in := make(chan string, 4)
	for _, w := range []string{"go", "rust", "zig", "ocaml"} {
		in <- w
	}
	close(in)

	out, errs := RunStage(context.Background(), in, 3, func(_ context.Context, s string) (string, error) {
		return strings.ToUpper(s), nil
	})

	for v := range out {
		fmt.Println(v)
	}
	for e := range errs {
		fmt.Println("error:", e.Err)
	}
}
