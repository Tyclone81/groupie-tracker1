package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"groupie-trackers/internal/handlers"
)

func TestRouting(t *testing.T) {
	h := handlers.NewHandler("../../templates/html")
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomeHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}