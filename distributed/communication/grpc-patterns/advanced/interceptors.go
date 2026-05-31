package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// UnaryLogging logs the method, duration, and error of every unary call.
func UnaryLogging(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	fmt.Printf("grpc method=%s dur=%s err=%v\n", info.FullMethod, time.Since(start), err)
	return resp, err
}
