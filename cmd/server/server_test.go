package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestHealthz(t *testing.T) {
	log, _ := zap.NewDevelopment()
	ts := httptest.NewServer(routes(log))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status: %d", resp.StatusCode)
	}
}
