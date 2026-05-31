package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	s := New()
	s.AddLiveness("self", func(_ context.Context) error { return nil })
	s.AddReadiness("database", func(_ context.Context) error { return nil })

	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	for _, path := range []string{"/healthz/live", "/healthz/ready"} {
		resp, _ := http.Get(srv.URL + path)
		fmt.Printf("%s -> %d\n", path, resp.StatusCode)
		resp.Body.Close()
	}
}
