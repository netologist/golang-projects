package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// sidecar wraps a handler to inject a header and log the request.
func sidecar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Sidecar", "v1")
		fmt.Printf("sidecar: %s %s (header injected)\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello; sidecar header=%s", r.Header.Get("X-Sidecar"))
	})

	srv := httptest.NewServer(sidecar(app))
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/")
	buf := make([]byte, 128)
	n, _ := resp.Body.Read(buf)
	resp.Body.Close()
	fmt.Println("response:", string(buf[:n]))
}
