package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	upstream := NewLoggingUpstream(1) // fail once, then succeed
	sc := New(upstream, 3)

	srv := httptest.NewServer(sc)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/")
	resp.Body.Close()
	fmt.Printf("final status=%d requests=%d retries=%d\n", resp.StatusCode, sc.Requests(), sc.Retries())
}
