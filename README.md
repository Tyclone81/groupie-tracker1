# groupie-tracker1
It is a Go web application that fetches data from an external RESTful API and visualizes band profiles, creation dates, first album releases, band members, and concert relations (locations & dates) in a responsive, user-friendly interface.
Built using strictly Go standard library packages (no external dependencies) following server-side rendering (SSR) best practices.

## Overview & Key Features

### Data Integration
Consumes and maps data from four RESTful API endpoints (artists, locations, dates, and relations).
 
### Client-Server Events
Implements dynamic search/filter handling via native HTML GET requests processed on the Go backend.
 
### Server Resilience
Comprehensive error handling for 400 Bad Request, 404 Not Found, and 500 Internal Server Error scenarios to ensure the server never panics or crashes.
### Zero External Dependencies
Built using standard Go packages (net/http, html/template, encoding/json, strconv, etc.).

### Unit Test Coverage
Automated package unit tests for API integration, data parsing, and HTTP status codes.

## Project Architecture
The project follows standard Go workspace conventions (cmd/, internal/) to enforce strict modularity and separation of concerns.

    groupie-tracker1/
    ├── cmd/
    │   └── web/
    │       ├── main.go              # Application entry point & route definitions
    │       └── main_test.go         # End-to-end integration and server routing tests
    ├── internal/
    │   ├── handlers/
    │   │   ├── handlers.go          # HTTP request handlers & search filter logic
    │   │   └── handlers_test.go     # HTTP status code and response unit tests
    │   └── models/
    │       ├── fetcher.go           # External API HTTP client logic
    │       ├── fetcher_test.go      # Data parsing and JSON unmarshaling unit tests
    │       └── structs.go           # Data blueprints mapping API JSON structures
    ├── static/
    │   └── css/
    │       └── styles.css           # Responsive grid layout and UI styling
    ├── templates/
    │   ├── artist.html              # Individual band details view
    │   ├── error.html               # Custom error page view (400, 404, 500)
    │   └── index.html               # Main page layout & search interface
    ├── go.mod                       # Go module definition file
    └── README.md                    # Project documentation


## How Files Connect (Data Flow)
When a client sends an HTTP request to the application, data flows through the file structure in a controlled loop as follows:

[ Browser ]
    │
    ▼
1. cmd/web/main.go  ────────────────► Serves static files (/static/) & routes URLs
    │
    ▼
2. internal/handlers/handlers.go ───► Receives HTTP request, checks paths & search params
    │
    ▼
3. internal/models/fetcher.go ──────► Makes external GET requests to the REST API
    │
    ▼
4. internal/models/structs.go ──────► Decodes raw JSON into typed Go data structs
    │
    ▼
5. internal/handlers/handlers.go ───► Injects Go structs into HTML templates
    │
    ▼
6. templates/*.html & CSS ──────────► Compiles HTML + applies styles.css and returns page


## Detailed Component Roles
cmd/web/main.go: Registers the static file server (./static), binds routes (/ and /artist), and starts http.ListenAndServe(":8080").

internal/models/structs.go: Defines the Go structs (Artist, Relation, ArtistDetail) matching incoming API JSON payloads.

internal/models/fetcher.go: Contains FetchArtists() and FetchRelation(id). Handles outgoing client calls with safety timeouts.

internal/handlers/handlers.go: Orchestrates incoming web requests. Filters artist data upon search query submission and safely calls RenderError if invalid routes/parameters are supplied.

templates/ & static/: Houses raw template blueprints (index.html, artist.html, error.html) and styling assets (styles.css).

## How to Run the Program
### Prerequisites
Go installed on your system.

### Execution Steps
Clone or navigate to the project root directory:
```Bash

cd groupie-tracker1
```
Start the web server:
```Bash

go run cmd/web/main.go
```
Open your browser and navigate to: http://localhost:8080

## How to Run Tests
The test suite validates backend API fetching, data integrity, error response handling, and route setup.
To execute all unit tests across all project packages, run:

```Bash

    go test -v ./...
```

## Manual Edge-Case Testing Guide
You can test server resilience and status code handling manually in your browser or terminal as follows:

Scenario                Test Request/Action                         Expected Result
Home Overview           GET http://localhost:8080/                  Status 200 OK. Displays band cards in grid view.
Search Action           Search Queen in home search bar             Status 200 OK. Filters cards matching term "Queen"
Artist Detail Page      GET http://localhost:8080/artist?id=1       Status 200 OK. Renders band details & concert relations.
404 Not Found           GET http://localhost:8080/non-existent      Status 404 Not Found. Renders custom error.html.
400 Bad Request         GET http://localhost:8080/artist?id=invalid Status 400 Bad Request. Renders custom error.html.
404 Out of Bounds       GET http://localhost:8080/artist?id=99999   Status 404 Not Found. Renders custom error.html.