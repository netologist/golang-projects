package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("unary:", resp.GetMessage())

	stream, err := c.SayManyHellos(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println("stream:", msg.GetMessage())
	}
}
