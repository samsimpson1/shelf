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
	"sync"
)

// App holds the application state
type App struct {
	mediaList         []Media
	mediaListMutex    sync.RWMutex // Protects concurrent access to mediaList
	scanner           *Scanner     // Scanner for refreshing media list
	templates         *template.Template
	mediaDir          string
	importDir         string // Path to import directory
	importScanner     *ImportScanner
	devMode           bool // Enable template hot-reloading in development
	tmdbClient        *TMDBClient
	musicBrainzClient *MusicBrainzClient
	playURLPrefix     string // URL prefix for play commands
	importSessions    *ImportSessionStore
}

// NewApp creates a new App instance
func NewApp(mediaList []Media, scanner *Scanner, templates *template.Template, mediaDir, importDir string) *App {
	var importScanner *ImportScanner
	if importDir != "" {
		importScanner = NewImportScanner(importDir)
	}

	return &App{
		mediaList:         mediaList,
		scanner:           scanner,
		templates:         templates,
		mediaDir:          mediaDir,
		importDir:         importDir,
		importScanner:     importScanner,
		devMode:           false,
		tmdbClient:        nil,
		musicBrainzClient: nil,
		playURLPrefix:     "",
		importSessions:    NewImportSessionStore(),
	}
}

// SetTMDBClient sets the TMDB client for the app
func (app *App) SetTMDBClient(client *TMDBClient) {
	app.tmdbClient = client
}

// SetMusicBrainzClient sets the MusicBrainz client for the app
func (app *App) SetMusicBrainzClient(client *MusicBrainzClient) {
	app.musicBrainzClient = client
}

// SetDevMode enables or disables development mode (template hot-reloading)
func (app *App) SetDevMode(enabled bool) {
	app.devMode = enabled
}

// SetPlayURLPrefix sets the URL prefix for play commands
func (app *App) SetPlayURLPrefix(prefix string) {
	app.playURLPrefix = prefix
}

// RefreshMediaList re-scans the media directory and updates the media list
func (app *App) RefreshMediaList() error {
	// Re-scan the media directory
	newMediaList, err := app.scanner.Scan()
	if err != nil {
		return fmt.Errorf("failed to refresh media list: %w", err)
	}

	// Update the media list with write lock
	app.mediaListMutex.Lock()
	app.mediaList = newMediaList
	app.mediaListMutex.Unlock()

	log.Printf("Media list refreshed: %d items", len(newMediaList))
	return nil
}

// templatesGlob is the single source of truth for which template files are loaded.
const templatesGlob = "templates/*.html"

// parseTemplates loads every template under templatesGlob.
func parseTemplates() (*template.Template, error) {
	return template.ParseGlob(templatesGlob)
}

// getTemplates returns the template set to use for rendering. In dev mode it
// re-parses from disk on every call so edits are picked up without a restart;
// otherwise it returns the cached set loaded at startup.
func (app *App) getTemplates() *template.Template {
	if !app.devMode {
		return app.templates
	}
	tmpl, err := parseTemplates()
	if err != nil {
		log.Printf("Error reloading templates: %v", err)
		return app.templates
	}
	return tmpl
}

// IndexHandler handles the main page request
func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := app.getTemplates()

	// Get a copy of the media list with read lock
	app.mediaListMutex.RLock()
	sorted := make([]Media, len(app.mediaList))
	copy(sorted, app.mediaList)
	app.mediaListMutex.RUnlock()

	// Sort media list: Films first, then TV shows, then Music, alphabetically within each type
	sort.Slice(sorted, func(i, j int) bool {
		// Sort by type first: Film < TV < Music
		if sorted[i].Type != sorted[j].Type {
			return sorted[i].Type < sorted[j].Type
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
	tmpl := app.getTemplates()

	slug := r.PathValue("slug")
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
	trackList := media.LoadTrackList()

	data := struct {
		Media         *Media
		Description   string
		Genres        []string
		HasPoster     bool
		TrackList     *TrackList
		PlayURLPrefix string
	}{
		Media:         media,
		Description:   description,
		Genres:        genres,
		HasPoster:     hasPoster,
		TrackList:     trackList,
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
	app.mediaListMutex.RLock()
	defer app.mediaListMutex.RUnlock()

	for i := range app.mediaList {
		if app.mediaList[i].Slug() == slug {
			// Return a copy to avoid concurrent access issues
			mediaCopy := app.mediaList[i]
			return &mediaCopy
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

	tmpl := app.getTemplates()

	slug := r.PathValue("slug")
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

	// Music does not support TMDB integration
	if media.Type == Music {
		http.Error(w, "TMDB search is not supported for Music media type", http.StatusBadRequest)
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

	tmpl := app.getTemplates()

	slug := r.PathValue("slug")
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

	// Music does not support TMDB integration
	if media.Type == Music {
		http.Error(w, "TMDB confirmation is not supported for Music media type", http.StatusBadRequest)
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
	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	slug := r.PathValue("slug")
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

	// Music does not support TMDB integration
	if media.Type == Music {
		http.Error(w, "Setting TMDB ID is not supported for Music media type", http.StatusBadRequest)
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

// SelectPosterHandler handles the poster selection page
func (app *App) SelectPosterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	tmpl := app.getTemplates()

	slug := r.PathValue("slug")
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

	// Check if media has TMDB ID
	if media.TMDBID == "" {
		http.Error(w, "Media does not have a TMDB ID", http.StatusBadRequest)
		return
	}

	// Fetch poster images from TMDB
	posters, err := app.tmdbClient.FetchPosterImages(media.TMDBID, media.Type)
	if err != nil {
		log.Printf("Failed to fetch poster images for %s: %v", media.Title, err)
		http.Error(w, "Failed to fetch poster images from TMDB", http.StatusInternalServerError)
		return
	}

	// Find current poster
	currentPosterPath, hasPoster := media.FindPosterFile()

	data := struct {
		Media             *Media
		Posters           []PosterImage
		HasCurrentPoster  bool
		CurrentPosterPath string
	}{
		Media:             media,
		Posters:           posters,
		HasCurrentPoster:  hasPoster,
		CurrentPosterPath: currentPosterPath,
	}

	err = tmpl.ExecuteTemplate(w, "select_poster.html", data)
	if err != nil {
		log.Printf("Error rendering select_poster template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// SavePosterHandler handles saving a selected poster
func (app *App) SavePosterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if TMDB client is available
	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	slug := r.PathValue("slug")
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

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get poster path from form
	posterPath := r.FormValue("poster_path")
	if posterPath == "" {
		http.Error(w, "Poster path is required", http.StatusBadRequest)
		return
	}

	// Delete existing poster (handles all extensions)
	err = DeleteExistingPoster(media.Path)
	if err != nil {
		log.Printf("Warning: Failed to delete existing poster for %s: %v", media.Title, err)
		// Continue anyway - we'll try to download the new poster
	}

	// Download the new poster
	err = app.tmdbClient.DownloadPoster(posterPath, media.Path)
	if err != nil {
		log.Printf("Failed to download new poster for %s: %v", media.Title, err)
		http.Error(w, "Failed to download poster", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully saved new poster for %s", media.Title)

	// Redirect back to detail page
	http.Redirect(w, r, "/media/"+url.PathEscape(slug), http.StatusSeeOther)
}
