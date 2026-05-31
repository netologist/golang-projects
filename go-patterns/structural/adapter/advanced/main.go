package main

import (
	"context"
	"fmt"
	"time"
)

type slowReader struct{}

func (s *slowReader) ReadSync(key string) (string, error) {
	time.Sleep(10 * time.Millisecond)
	return "value-of-" + key, nil
}

func main() {
	r := NewAsyncAdapter(&slowReader{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res := <-r.Read(ctx, "config")
	fmt.Printf("value=%s err=%v\n", res.Value, res.Err)
}
