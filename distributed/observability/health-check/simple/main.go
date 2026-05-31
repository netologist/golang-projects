package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	// Demonstrate via httptest so the demo is self-contained.
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	healthHandler(rr, req)
	fmt.Printf("status=%d body=%s", rr.Code, rr.Body.String())
}
