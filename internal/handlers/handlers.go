package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"groupie-tracker1/internal/models"
)

// Handler holds dependencies needed across all HTTP request handlers,
// such as the directory path where HTML template files are located.
type Handler struct {
	TemplateDir string
}

// NewHandler initializes a new Handler instance with the specified template directory.
func NewHandler(templateDir string) *Handler {
	return &Handler{TemplateDir: templateDir}
}

// RenderError renders custom error pages (e.g., 400, 404, 500) without crashing the server.
func (h *Handler) RenderError(w http.ResponseWriter, statusCode int, message string) {
	// 1. Try to load custom error page
	tmpl, err := template.ParseFiles(filepath.Join(h.TemplateDir, "error.html"))
	if err != nil {
		// Fallback to standard HTTP error if template parsing fails (avoids double writeHeader)
		http.Error(w, message, statusCode)
		return
	}

	// 2. Set status header BEFORE writing the body
	w.WriteHeader(statusCode)
	
	// 3. Render error template
	tmpl.Execute(w, map[string]interface{}{
		"StatusCode": statusCode,
		"Message":    message,
	})
}

// HomeHandler serves the main landing page ("/") and processes search bar queries.
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Enforce exact route matching so "/unknown-path" doesn't trigger the home page
	if r.URL.Path != "/" {
		h.RenderError(w, http.StatusNotFound, "Page Not Found (404)")
		return
	}

	// Enforce HTTP GET method
	if r.Method != http.MethodGet {
		h.RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed (405)")
		return
	}

	// Retrieve artist data from external API via our models package
	artists, err := models.FetchArtists()
	if err != nil {
		h.RenderError(w, http.StatusInternalServerError, "Failed to load band data from external API")
		return
	}

	// Extract search query parameter if user submitted the search form (?search=...)
	query := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("search")))
	if query != "" {
		var filtered []models.Artist
		for _, a := range artists {
			// Check if search query matches Band Name, Creation Date, or First Album
			match := strings.Contains(strings.ToLower(a.Name), query) ||
				strings.Contains(strconv.Itoa(a.CreationDate), query) ||
				strings.Contains(strings.ToLower(a.FirstAlbum), query)

			// Check if search query matches any Member Name
			if !match {
				for _, m := range a.Members {
					if strings.Contains(strings.ToLower(m), query) {
						match = true
						break
					}
				}
			}

			if match {
				filtered = append(filtered, a)
			}
		}
		artists = filtered // Replace main list with filtered results
	}

	// Parse and execute index.html, injecting the filtered/unfiltered artists list
	tmpl, err := template.ParseFiles(filepath.Join(h.TemplateDir, "index.html"))
	if err != nil {
		h.RenderError(w, http.StatusInternalServerError, "Failed to render home template")
		return
	}

	tmpl.Execute(w, artists)
}

// ArtistHandler handles single band detail requests ("/artist?id=X").
func (h *Handler) ArtistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed (405)")
		return
	}

	// Extract and validate the 'id' parameter from the request URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.RenderError(w, http.StatusBadRequest, "Invalid Artist ID (400)")
		return
	}

	// Fetch all artists to verify the requested ID exists
	artists, err := models.FetchArtists()
	if err != nil {
		h.RenderError(w, http.StatusInternalServerError, "Failed to connect to backend service")
		return
	}

	var selectedArtist *models.Artist
	for _, a := range artists {
		if a.ID == id {
			selectedArtist = &a
			break
		}
	}

	// If ID is numeric but doesn't correspond to any artist (e.g., id=99999)
	if selectedArtist == nil {
		h.RenderError(w, http.StatusNotFound, "Artist Not Found (404)")
		return
	}

	// Fetch concert relations for this specific artist ID
	relation, err := models.FetchRelation(id)
	if err != nil {
		h.RenderError(w, http.StatusInternalServerError, "Failed to load concert relation data")
		return
	}

	// Combine artist information and concert relations into one struct
	data := models.ArtistDetail{
		Artist:   *selectedArtist,
		Relation: relation,
	}

	// Execute artist.html with the bundled data
	tmpl, err := template.ParseFiles(filepath.Join(h.TemplateDir, "artist.html"))
	if err != nil {
		h.RenderError(w, http.StatusInternalServerError, "Failed to render artist page")
		return
	}

	tmpl.Execute(w, data)
}