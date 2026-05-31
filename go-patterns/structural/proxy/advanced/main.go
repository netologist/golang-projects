package main

import (
	"context"
	"fmt"
	"time"
)

type realLoader struct{}

func (r *realLoader) Load(_ context.Context, url string) ([]byte, error) {
	return []byte("payload-for-" + url), nil
}

func main() {
	// Compose: Auth( Caching( real ) )
	loader := NewAuthProxy(NewCachingProxy(&realLoader{}, time.Minute), "secret")

	ctx := context.WithValue(context.Background(), TokenKey, "secret")
	data, err := loader.Load(ctx, "/resource")
	fmt.Printf("authorized: data=%s err=%v\n", data, err)

	_, err = loader.Load(context.Background(), "/resource")
	fmt.Printf("unauthorized: err=%v\n", err)
}
