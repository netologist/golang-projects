package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message:   fmt.Sprintf("Hello, %s!", req.GetName()),
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

func (s *server) SayManyHellos(req *pb.HelloRequest, stream pb.Echo_SayManyHellosServer) error {
	for i := 1; i <= 3; i++ {
		if err := stream.Send(&pb.HelloResponse{
			Message:   fmt.Sprintf("Hello #%d, %s!", i, req.GetName()),
			Timestamp: time.Now().UnixMilli(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	pb.RegisterEchoServer(srv, &server{})
	log.Println("gRPC server listening on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
