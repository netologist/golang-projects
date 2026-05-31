package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
)

// echoServer implements the generated Echo service.
type echoServer struct {
	pb.UnimplementedEchoServer
}

func (s *echoServer) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message:   fmt.Sprintf("Hello, %s!", req.GetName()),
		Timestamp: time.Now().UnixMilli(),
	}, nil
}
