package main

import (
	"context"
	"net"
	"testing"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() (pb.EchoClient, func()) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer(grpc.UnaryInterceptor(UnaryLogging))
	pb.RegisterEchoServer(srv, &echoServer{})
	go func() { _ = srv.Serve(lis) }()

	conn, _ := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	cleanup := func() {
		_ = conn.Close()
		srv.Stop()
	}
	return pb.NewEchoClient(conn), cleanup
}

func TestEcho_SayHello(t *testing.T) {
	client, cleanup := dialer()
	defer cleanup()

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Test"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetMessage() != "Hello, Test!" {
		t.Errorf("got %q, want %q", resp.GetMessage(), "Hello, Test!")
	}
	if resp.GetTimestamp() == 0 {
		t.Error("expected non-zero timestamp")
	}
}
