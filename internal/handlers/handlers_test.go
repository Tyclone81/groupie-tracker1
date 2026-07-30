package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomeHandler tests the routing and HTTP status code behaviors for the home handler.
func TestHomeHandler(t *testing.T) {
	// Initialize our Handler with the relative path to HTML template files
	h := NewHandler("../../templates/html")

	// Sub-test 1: Verify valid GET request to "/" returns HTTP 200 OK
	t.Run("Valid Home Route Status 200", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder() // Record the response headers and body

		h.HomeHandler(rr, req)
		
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("HomeHandler returned wrong status: got %v want %v", status, http.StatusOK)
		}
	})

	// Sub-test 2: Verify non-existent sub-paths (e.g. "/unknown") trigger HTTP 404 Not Found
	t.Run("Unknown Route Status 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/invalid-route", nil)
		rr := httptest.NewRecorder()

		h.HomeHandler(rr, req)
		
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("HomeHandler returned wrong status for 404: got %v want %v", status, http.StatusNotFound)
		}
	})
}

// TestArtistHandlerBadRequest tests edge-case URL parameters on the /artist handler.
func TestArtistHandlerBadRequest(t *testing.T) {
	h := NewHandler("../../templates/html")

	// Pass a non-numeric string as the artist ID (?id=invalid)
	req, _ := http.NewRequest("GET", "/artist?id=invalid", nil)
	rr := httptest.NewRecorder()

	h.ArtistHandler(rr, req)

	// Ensure handler safely responds with HTTP 400 Bad Request instead of crashing
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("ArtistHandler returned wrong status for bad ID: got %v want %v", status, http.StatusBadRequest)
	}
}