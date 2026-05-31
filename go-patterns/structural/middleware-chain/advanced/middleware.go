package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type ctxKey string

// RequestIDKey carries the per-request ID in the context.
const RequestIDKey ctxKey = "request-id"

// Request is the inbound message.
type Request struct{ Body string }

// Response is the outbound message.
type Response struct {
	Body       string
	StatusCode int
}

// HandlerFunc processes a Request.
type HandlerFunc func(ctx context.Context, req Request) (Response, error)

// Middleware wraps a HandlerFunc.
type Middleware func(HandlerFunc) HandlerFunc

// Chain applies middleware so mw[0] is the outermost layer.
func Chain(h HandlerFunc, mw ...Middleware) HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

// RequestID injects a unique request ID into the context.
func RequestID(next HandlerFunc) HandlerFunc {
	var counter uint64
	return func(ctx context.Context, req Request) (Response, error) {
		id := fmt.Sprintf("req-%d", atomic.AddUint64(&counter, 1))
		ctx = context.WithValue(ctx, RequestIDKey, id)
		return next(ctx, req)
	}
}

// Logging logs entry, exit, duration, and error.
func Logging(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		id, _ := ctx.Value(RequestIDKey).(string)
		start := time.Now()
		resp, err := next(ctx, req)
		fmt.Printf("id=%s body=%q status=%d dur=%s err=%v\n",
			id, req.Body, resp.StatusCode, time.Since(start), err)
		return resp, err
	}
}

// Recovery catches panics and converts them to a 500 response.
func Recovery(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (resp Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered: %v\n%s", r, debug.Stack())
				resp = Response{StatusCode: 500}
				err = fmt.Errorf("internal error")
			}
		}()
		return next(ctx, req)
	}
}
