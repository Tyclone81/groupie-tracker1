package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BaseAPIURL is the primary entry point for fetching remote groupie trackers data.
const BaseAPIURL = "https://groupietrackers.herokuapp.com/api"

// Create a reusable HTTP client with a 10-second timeout to prevent the application
// from hanging indefinitely if the external API is slow or unresponsive.
var client = &http.Client{Timeout: 10 * time.Second}

// FetchArtists makes an HTTP GET request to the external API to retrieve all band data.
func FetchArtists() ([]Artist, error) {
	// Execute the HTTP GET request to /api/artists
	resp, err := client.Get(BaseAPIURL + "/artists")
	if err != nil {
		return nil, fmt.Errorf("failed to reach external API: %w", err)
	}
	// Always close the response body when finished to prevent network resource memory leaks
	defer resp.Body.Close()

	// Verify that the external API returned an HTTP 200 OK status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API status: %d", resp.StatusCode)
	}

	// Decode the raw JSON response body directly into a slice of Artist structs
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return artists, nil
}

// FetchRelation fetches the concert dates/locations mapping for a single artist by ID.
func FetchRelation(id int) (Relation, error) {
	// Construct the dynamic URL for the relation endpoint (e.g., /api/relation/1)
	url := fmt.Sprintf("%s/relation/%d", BaseAPIURL, id)
	
	resp, err := client.Get(url)
	if err != nil {
		return Relation{}, fmt.Errorf("failed to reach relation API: %w", err)
	}
	defer resp.Body.Close()

	// Ensure the response was successful
	if resp.StatusCode != http.StatusOK {
		return Relation{}, fmt.Errorf("relation API status: %d", resp.StatusCode)
	}

	// Decode the relation JSON data into our Relation struct
	var relation Relation
	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return Relation{}, fmt.Errorf("failed to decode relation JSON: %w", err)
	}

	return relation, nil
}