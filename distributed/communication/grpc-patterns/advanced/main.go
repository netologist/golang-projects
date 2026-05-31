package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Start an in-process server with the logging interceptor on an ephemeral port.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer(grpc.UnaryInterceptor(UnaryLogging))
	pb.RegisterEchoServer(srv, &echoServer{})
	go func() { _ = srv.Serve(lis) }()
	defer srv.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := pb.NewEchoClient(conn).SayHello(context.Background(), &pb.HelloRequest{Name: "Interceptor"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response:", resp.GetMessage())
}
