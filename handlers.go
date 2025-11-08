package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

// App holds the application state
type App struct {
	mediaList []Media
	templates *template.Template
	mediaDir  string
}

// NewApp creates a new App instance
func NewApp(mediaList []Media, templates *template.Template, mediaDir string) *App {
	return &App{
		mediaList: mediaList,
		templates: templates,
		mediaDir:  mediaDir,
	}
}

// IndexHandler handles the main page request
func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
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
		MediaList []Media
	}{
		MediaList: sorted,
	}

	err := app.templates.ExecuteTemplate(w, "index.html", data)
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
		Media       *Media
		Description string
		Genres      []string
		HasPoster   bool
	}{
		Media:       media,
		Description: description,
		Genres:      genres,
		HasPoster:   hasPoster,
	}

	err := app.templates.ExecuteTemplate(w, "detail.html", data)
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
