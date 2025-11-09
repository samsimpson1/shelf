package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// App holds the application state
type App struct {
	mediaList      []Media
	templates      *template.Template
	mediaDir       string
	importDir      string // Path to import directory
	importScanner  *ImportScanner
	devMode        bool // Enable template hot-reloading in development
	tmdbClient     *TMDBClient
	playURLPrefix  string // URL prefix for play commands
}

// NewApp creates a new App instance
func NewApp(mediaList []Media, templates *template.Template, mediaDir, importDir string) *App {
	var importScanner *ImportScanner
	if importDir != "" {
		importScanner = NewImportScanner(importDir)
	}

	return &App{
		mediaList:     mediaList,
		templates:     templates,
		mediaDir:      mediaDir,
		importDir:     importDir,
		importScanner: importScanner,
		devMode:       false,
		tmdbClient:    nil,
		playURLPrefix: "",
	}
}

// SetTMDBClient sets the TMDB client for the app
func (app *App) SetTMDBClient(client *TMDBClient) {
	app.tmdbClient = client
}

// SetDevMode enables or disables development mode (template hot-reloading)
func (app *App) SetDevMode(enabled bool) {
	app.devMode = enabled
}

// SetPlayURLPrefix sets the URL prefix for play commands
func (app *App) SetPlayURLPrefix(prefix string) {
	app.playURLPrefix = prefix
}

// loadTemplates reloads templates from disk (used in dev mode)
func (app *App) loadTemplates() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/index.html",
		"templates/detail.html",
		"templates/search.html",
		"templates/confirm.html",
		"templates/import_list.html",
		"templates/import_step1.html",
		"templates/import_step2.html",
		"templates/import_step3.html",
		"templates/import_step4.html",
		"templates/import_step5.html",
		"templates/import_confirm.html",
		"templates/import_success.html",
	)
	if err != nil {
		log.Printf("Error reloading templates: %v", err)
		return app.templates // Fall back to cached templates
	}
	return tmpl
}

// IndexHandler handles the main page request
func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Sort media list: Films first, then TV shows, alphabetically within each type
	sorted := make([]Media, len(app.mediaList))
	copy(sorted, app.mediaList)

	sort.Slice(sorted, func(i, j int) bool {
		// Films come before TV shows
		if sorted[i].Type != sorted[j].Type {
			return sorted[i].Type == Film
		}
		// Within same type, sort alphabetically by title
		return sorted[i].Title < sorted[j].Title
	})

	data := struct {
		MediaList     []Media
		ImportEnabled bool
	}{
		MediaList:     sorted,
		ImportEnabled: app.importScanner != nil,
	}

	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// PosterHandler serves poster images for media items
func (app *App) PosterHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL: /posters/{slug}
	slug := strings.TrimPrefix(r.URL.Path, "/posters/")
	slug = strings.TrimSuffix(slug, "/")

	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// Find media by slug
	media := app.findMediaBySlug(slug)
	if media == nil {
		http.NotFound(w, r)
		return
	}

	// Find poster file
	posterPath, found := media.FindPosterFile()
	if !found {
		http.NotFound(w, r)
		return
	}

	// Validate path is within media directory (security check)
	cleanPath := filepath.Clean(posterPath)
	if !strings.HasPrefix(cleanPath, filepath.Clean(app.mediaDir)) {
		log.Printf("Security warning: attempted access to path outside media dir: %s", cleanPath)
		http.NotFound(w, r)
		return
	}

	// Serve the file with appropriate content type
	http.ServeFile(w, r, posterPath)
}

// DetailHandler handles individual media detail pages
func (app *App) DetailHandler(w http.ResponseWriter, r *http.Request) {
	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Extract slug from URL: /media/{slug}
	slug := strings.TrimPrefix(r.URL.Path, "/media/")
	slug = strings.TrimSuffix(slug, "/")

	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// Find media by slug
	media := app.findMediaBySlug(slug)
	if media == nil {
		http.NotFound(w, r)
		return
	}

	// Load additional metadata
	description := media.LoadDescription()
	genres := media.LoadGenres()
	_, hasPoster := media.FindPosterFile()

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		PlayURLPrefix string
	}{
		Media:         media,
		Description:   description,
		Genres:        genres,
		HasPoster:     hasPoster,
		PlayURLPrefix: app.playURLPrefix,
	}

	err := tmpl.ExecuteTemplate(w, "detail.html", data)
	if err != nil {
		log.Printf("Error rendering detail template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// findMediaBySlug finds a media item by its slug
func (app *App) findMediaBySlug(slug string) *Media {
	for i := range app.mediaList {
		if app.mediaList[i].Slug() == slug {
			return &app.mediaList[i]
		}
	}
	return nil
}

// SearchTMDBHandler handles the TMDB search page
func (app *App) SearchTMDBHandler(w http.ResponseWriter, r *http.Request) {
	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Extract slug from URL: /media/{slug}/search-tmdb
	path := strings.TrimPrefix(r.URL.Path, "/media/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	slug := parts[0]

	// Find media by slug
	media := app.findMediaBySlug(slug)
	if media == nil {
		http.NotFound(w, r)
		return
	}

	// Get query parameters
	query := r.URL.Query().Get("query")
	yearStr := r.URL.Query().Get("year")
	year := 0

	// Parse year if provided
	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil && parsedYear > 0 {
			year = parsedYear
		}
	}

	// Pre-fill query with media title if not provided
	if query == "" && r.Method == "GET" {
		query = media.Title
	}

	var results interface{}
	var searchErr error

	// Perform search if query is provided
	if query != "" && r.Method == "GET" {
		if media.Type == Film {
			// Search for movies
			if year == 0 && media.Year > 0 {
				year = media.Year
			}
			movieResults, err := app.tmdbClient.SearchMovies(query, year)
			if err != nil {
				searchErr = err
			} else {
				results = movieResults
			}
		} else if media.Type == TV {
			// Search for TV shows
			tvResults, err := app.tmdbClient.SearchTV(query)
			if err != nil {
				searchErr = err
			} else {
				results = tvResults
			}
		}
	}

	// Prepare error message
	var errorMsg string
	if searchErr != nil {
		errorMsg = fmt.Sprintf("Search error: %v", searchErr)
	}

	data := struct {
		Media   *Media
		Query   string
		Year    int
		Results interface{}
		Error   string
	}{
		Media:   media,
		Query:   query,
		Year:    year,
		Results: results,
		Error:   errorMsg,
	}

	err := tmpl.ExecuteTemplate(w, "search.html", data)
	if err != nil {
		log.Printf("Error rendering search template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ConfirmTMDBHandler handles the TMDB confirmation page
func (app *App) ConfirmTMDBHandler(w http.ResponseWriter, r *http.Request) {
	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Extract slug from URL: /media/{slug}/confirm-tmdb
	path := strings.TrimPrefix(r.URL.Path, "/media/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	slug := parts[0]

	// Find media by slug
	media := app.findMediaBySlug(slug)
	if media == nil {
		http.NotFound(w, r)
		return
	}

	// Get TMDB ID from query parameter
	tmdbID := r.URL.Query().Get("id")
	if tmdbID == "" {
		http.Error(w, "TMDB ID is required", http.StatusBadRequest)
		return
	}

	// Get optional query parameter for back navigation
	query := r.URL.Query().Get("query")

	// Fetch TMDB match details
	var tmdbMatch interface{}
	var fetchErr error

	if media.Type == Film {
		movieData, err := app.tmdbClient.FetchMovieMetadata(tmdbID)
		if err != nil {
			fetchErr = err
		} else {
			// Convert to search result format for template
			tmdbMatch = MovieSearchResult{
				ID:          movieData.ID,
				Title:       movieData.Title,
				ReleaseDate: movieData.ReleaseDate,
				Overview:    movieData.Overview,
				PosterPath:  movieData.PosterPath,
			}
		}
	} else if media.Type == TV {
		tvData, err := app.tmdbClient.FetchTVMetadata(tmdbID)
		if err != nil {
			fetchErr = err
		} else {
			// Convert to search result format for template
			tmdbMatch = TVSearchResult{
				ID:           tvData.ID,
				Name:         tvData.Name,
				FirstAirDate: tvData.FirstAirDate,
				Overview:     tvData.Overview,
				PosterPath:   tvData.PosterPath,
			}
		}
	}

	// Load current media metadata
	description := media.LoadDescription()
	_, hasPoster := media.FindPosterFile()

	// Prepare error message
	var errorMsg string
	if fetchErr != nil {
		errorMsg = fmt.Sprintf("Failed to fetch TMDB details: %v", fetchErr)
	}

	data := struct {
		Media       *Media
		TMDBID      string
		TMDBMatch   interface{}
		Query       string
		Description string
		HasPoster   bool
		Error       string
	}{
		Media:       media,
		TMDBID:      tmdbID,
		TMDBMatch:   tmdbMatch,
		Query:       query,
		Description: description,
		HasPoster:   hasPoster,
		Error:       errorMsg,
	}

	err := tmpl.ExecuteTemplate(w, "confirm.html", data)
	if err != nil {
		log.Printf("Error rendering confirm template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// SaveTMDBHandler handles saving the TMDB ID
func (app *App) SaveTMDBHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	// Extract slug from URL: /media/{slug}/set-tmdb
	path := strings.TrimPrefix(r.URL.Path, "/media/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	slug := parts[0]

	// Find media by slug
	media := app.findMediaBySlug(slug)
	if media == nil {
		http.NotFound(w, r)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get TMDB ID from form
	tmdbID := r.FormValue("tmdb_id")
	if tmdbID == "" {
		http.Error(w, "TMDB ID is required", http.StatusBadRequest)
		return
	}

	// Validate TMDB ID
	err = app.tmdbClient.ValidateTMDBID(tmdbID, media.Type)
	if err != nil {
		log.Printf("Invalid TMDB ID %s for %s: %v", tmdbID, media.Title, err)
		http.Error(w, fmt.Sprintf("Invalid TMDB ID: %v", err), http.StatusBadRequest)
		return
	}

	// Write TMDB ID to file
	err = WriteTMDBID(tmdbID, media.Path)
	if err != nil {
		log.Printf("Failed to write TMDB ID for %s: %v", media.Title, err)
		http.Error(w, "Failed to save TMDB ID", http.StatusInternalServerError)
		return
	}

	// Update media object
	media.TMDBID = tmdbID

	// Check if metadata should be downloaded now
	downloadMetadata := r.FormValue("download_metadata") == "true"
	if downloadMetadata {
		err = app.tmdbClient.FetchAndSaveMetadata(media)
		if err != nil {
			log.Printf("Warning: Failed to fetch metadata for %s: %v", media.Title, err)
			// Don't fail the request, just log the warning
		} else {
			log.Printf("Successfully fetched metadata for %s", media.Title)
		}
	}

	// Redirect back to detail page
	http.Redirect(w, r, "/media/"+url.PathEscape(slug), http.StatusSeeOther)
}
