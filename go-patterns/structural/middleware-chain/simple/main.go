package main

import "context"

func main() {
	final := func(_ context.Context, req Request) (Response, error) {
		return Response{Body: "handled: " + req.Body}, nil
	}
	h := Chain(final, LoggingMW, UppercaseMW)
	h(context.Background(), Request{Body: "hello"})
}
