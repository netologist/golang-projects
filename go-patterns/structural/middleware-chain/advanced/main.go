package main

import "context"

func main() {
	final := func(_ context.Context, req Request) (Response, error) {
		return Response{Body: "ok: " + req.Body, StatusCode: 200}, nil
	}
	// Recovery outermost, then RequestID, then Logging.
	h := Chain(final, Recovery, RequestID, Logging)
	h(context.Background(), Request{Body: "ping"})
}
