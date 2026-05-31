package main

import (
	"context"
	"fmt"
	"strings"
)

// Request is the inbound message.
type Request struct{ Body string }

// Response is the outbound message.
type Response struct{ Body string }

// HandlerFunc processes a Request.
type HandlerFunc func(ctx context.Context, req Request) (Response, error)

// Middleware wraps a HandlerFunc to add behaviour.
type Middleware func(HandlerFunc) HandlerFunc

// Chain applies middleware so mw[0] is the outermost layer.
func Chain(h HandlerFunc, mw ...Middleware) HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

// LoggingMW logs entry and exit.
func LoggingMW(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		fmt.Printf("-> req: %s\n", req.Body)
		resp, err := next(ctx, req)
		fmt.Printf("<- resp: %s err: %v\n", resp.Body, err)
		return resp, err
	}
}

// UppercaseMW upper-cases the request body before passing it on.
func UppercaseMW(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		req.Body = strings.ToUpper(req.Body)
		return next(ctx, req)
	}
}
