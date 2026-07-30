package main

import (
	"log"
	"net/http"

	"groupie-tracker1/internal/handlers"
)

func main() {
	// Initialize handler helper with relative path to HTML templates
	h := handlers.NewHandler("templates/html")
	
	// Create a new request multiplexer (HTTP router)
	mux := http.NewServeMux()

	// Correct File Server Setup:
	// Serve files from the "./static" directory.
	// StripPrefix removes "/static" from requested URLs like "/static/css/styles.css",
	// allowing http.FileServer to look up "css/styles.css" inside "./static".
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static", http.StripPrefix("/static", fileServer))

	// Register page handlers
	mux.HandleFunc("/", h.HomeHandler)        // Home page & Search
	mux.HandleFunc("/artist", h.ArtistHandler) // Individual Artist Detail page

	// Start listening for HTTP connections
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server stopped unexpectedly: %v", err)
	}
}