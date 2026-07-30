package models

// Artist represents the primary JSON structure returned by the /api/artists endpoint.
// The `json:"..."` tags tell Go's JSON parser how to map JSON key names to Go struct fields.
type Artist struct {
	ID           int      `json:"id"`           // Unique identifier for the band/artist
	Image        string   `json:"image"`        // URL to the band's image
	Name         string   `json:"name"`         // Name of the band or artist
	Members      []string `json:"members"`      // List of individual band member names
	CreationDate int      `json:"creationDate"` // Year the band was formed
	FirstAlbum   string   `json:"firstAlbum"`   // Release date of their first album
	Locations    string   `json:"locations"`    // API sub-endpoint URL for concert locations
	ConcertDates string   `json:"concertDates"` // API sub-endpoint URL for concert dates
	Relations    string   `json:"relations"`    // API sub-endpoint URL for location/date relations
}

// Relation maps the JSON response from /api/relation/{id}.
// It links concert dates directly to specific geographic locations.
type Relation struct {
	ID             int                 `json:"id"`             // Unique identifier matching the Artist ID
	DatesLocations map[string][]string `json:"datesLocations"` // Map where Key = Location name, Value = Slice of date strings
}

// ArtistDetail is a composite data structure used specifically for template rendering.
// It bundles both the core Artist info and their specific Relation data so it can be passed to artist.html at once.
type ArtistDetail struct {
	Artist   Artist   // The main band information
	Relation Relation // The combined concert dates and locations for this band
}