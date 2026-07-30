package models

import (
	"testing"
)

// TestFetchArtists verifies that our API client can successfully fetch 
// and parse the list of bands from the remote /artists endpoint.
func TestFetchArtists(t *testing.T) {
	// Execute the function that makes the live HTTP GET call to /api/artists
	artists, err := FetchArtists()
	
	// Fail immediately if network or JSON unmarshaling fails
	if err != nil {
		t.Fatalf("Expected no error when fetching artists, got: %v", err)
	}

	// Ensure the returned slice actually contains parsed records
	if len(artists) == 0 {
		t.Errorf("Expected non-empty list of artists, got 0 items")
	}

	// Verify that structural fields (like Name) were properly mapped from JSON
	if artists[0].Name == "" {
		t.Errorf("Expected first artist to have a valid Name, got empty string")
	}
}

// TestFetchRelation checks that concert locations and dates for a specific 
// artist ID are correctly retrieved and parsed into a Relation struct.
func TestFetchRelation(t *testing.T) {
	// Query relation data for the first band (ID = 1)
	relation, err := FetchRelation(1)
	if err != nil {
		t.Fatalf("Expected valid relation data for ID 1, got error: %v", err)
	}

	// Verify that the returned relation ID matches what was requested
	if relation.ID != 1 {
		t.Errorf("Expected relation ID to be 1, got %d", relation.ID)
	}
}