package main

import (
	"context"
	"fmt"
)

func handler(ctx context.Context) {
	service(ctx)
}

func service(ctx context.Context) {
	if id, ok := RequestID(ctx); ok {
		fmt.Println("service sees request id:", id)
	}
}

func main() {
	ctx := WithRequestID(context.Background(), "abc-123")
	handler(ctx)
}
