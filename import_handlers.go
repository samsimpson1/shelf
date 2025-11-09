package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
)

// ImportSessionStore stores active import sessions
type ImportSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*ImportSession
	counter  uint64
}

// NewImportSessionStore creates a new import session store
func NewImportSessionStore() *ImportSessionStore {
	return &ImportSessionStore{
		sessions: make(map[string]*ImportSession),
		counter:  0,
	}
}

// Create creates a new import session and returns its ID
func (s *ImportSessionStore) Create(session *ImportSession) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate simple sequential ID
	id := fmt.Sprintf("import-%d", atomic.AddUint64(&s.counter, 1))
	s.sessions[id] = session
	return id
}

// Get retrieves an import session by ID
func (s *ImportSessionStore) Get(id string) (*ImportSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	return session, ok
}

// Delete removes an import session
func (s *ImportSessionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
}

// Global session store
var importSessionStore = NewImportSessionStore()

// ImportListHandler shows the list of directories available for import
func (app *App) ImportListHandler(w http.ResponseWriter, r *http.Request) {
	// Check if import is enabled
	if app.importScanner == nil {
		http.Error(w, "Import functionality is not configured (IMPORT_DIR not set)", http.StatusServiceUnavailable)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Scan import directory
	imports, err := app.importScanner.Scan()
	if err != nil {
		log.Printf("Error scanning import directory: %v", err)
		http.Error(w, "Failed to scan import directory", http.StatusInternalServerError)
		return
	}

	data := struct {
		Imports []ImportDirectory
	}{
		Imports: imports,
	}

	err = tmpl.ExecuteTemplate(w, "import_list.html", data)
	if err != nil {
		log.Printf("Error rendering import_list template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportStartHandler starts an import session for a selected directory
func (app *App) ImportStartHandler(w http.ResponseWriter, r *http.Request) {
	// Check if import is enabled
	if app.importScanner == nil {
		http.Error(w, "Import functionality is not configured", http.StatusServiceUnavailable)
		return
	}

	// Get directory name from query parameter
	dirName := r.URL.Query().Get("dir")
	if dirName == "" {
		http.Error(w, "Directory name is required", http.StatusBadRequest)
		return
	}

	// Scan to find the directory
	imports, err := app.importScanner.Scan()
	if err != nil {
		http.Error(w, "Failed to scan import directory", http.StatusInternalServerError)
		return
	}

	var selectedDir *ImportDirectory
	for i := range imports {
		if imports[i].Name == dirName {
			selectedDir = &imports[i]
			break
		}
	}

	if selectedDir == nil {
		http.NotFound(w, r)
		return
	}

	// Detect disk type
	detectedType, _ := DetectDiskType(selectedDir.Path)

	// Create new import session
	session := &ImportSession{
		SourceDir:    selectedDir,
		DetectedType: detectedType,
	}

	// Store session and get ID
	sessionID := importSessionStore.Create(session)

	// Redirect to step 1 (choose media kind)
	http.Redirect(w, r, "/import/step1?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
}

// ImportStep1Handler handles step 1: choose media kind (Film/TV)
func (app *App) ImportStep1Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		mediaKindStr := r.FormValue("media_kind")
		if mediaKindStr == "film" {
			session.MediaKind = Film
		} else if mediaKindStr == "tv" {
			session.MediaKind = TV
		} else {
			http.Error(w, "Invalid media kind", http.StatusBadRequest)
			return
		}

		// Redirect to step 2 (TMDB search or manual entry)
		http.Redirect(w, r, "/import/step2?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step1.html", data)
	if err != nil {
		log.Printf("Error rendering import_step1 template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportStep2Handler handles step 2: TMDB search or manual entry
func (app *App) ImportStep2Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Handle "skip TMDB" action
	if r.Method == http.MethodPost && r.FormValue("action") == "skip" {
		// Redirect to manual entry (step 3)
		http.Redirect(w, r, "/import/step3?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Handle TMDB search
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

	var results interface{}
	var searchErr error

	// Perform search if query is provided
	if query != "" && app.tmdbClient != nil {
		if session.MediaKind == Film {
			movieResults, err := app.tmdbClient.SearchMovies(query, year)
			if err != nil {
				searchErr = err
			} else {
				results = movieResults
			}
		} else if session.MediaKind == TV {
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

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	data := struct {
		Session       *ImportSession
		SessionID     string
		Query         string
		Year          int
		Results       interface{}
		Error         string
		TMDBAvailable bool
	}{
		Session:       session,
		SessionID:     sessionID,
		Query:         query,
		Year:          year,
		Results:       results,
		Error:         errorMsg,
		TMDBAvailable: app.tmdbClient != nil,
	}

	err := tmpl.ExecuteTemplate(w, "import_step2.html", data)
	if err != nil {
		log.Printf("Error rendering import_step2 template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportStep2ConfirmHandler handles TMDB match selection
func (app *App) ImportStep2ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	tmdbID := r.URL.Query().Get("id")

	if sessionID == "" || tmdbID == "" {
		http.Error(w, "Session ID and TMDB ID are required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	if app.tmdbClient == nil {
		http.Error(w, "TMDB API is not configured", http.StatusServiceUnavailable)
		return
	}

	// Fetch metadata from TMDB
	if session.MediaKind == Film {
		movie, err := app.tmdbClient.FetchMovieMetadata(tmdbID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch movie metadata: %v", err), http.StatusInternalServerError)
			return
		}
		session.TMDBID = tmdbID
		session.TMDBTitle = movie.Title
		session.TMDBOverview = movie.Overview
		// Extract year from release date
		if len(movie.ReleaseDate) >= 4 {
			if year, err := strconv.Atoi(movie.ReleaseDate[:4]); err == nil {
				session.TMDBYear = year
			}
		}
		// Extract genre names
		session.TMDBGenres = make([]string, len(movie.Genres))
		for i, genre := range movie.Genres {
			session.TMDBGenres[i] = genre.Name
		}
	} else if session.MediaKind == TV {
		tv, err := app.tmdbClient.FetchTVMetadata(tmdbID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch TV metadata: %v", err), http.StatusInternalServerError)
			return
		}
		session.TMDBID = tmdbID
		session.TMDBTitle = tv.Name
		session.TMDBOverview = tv.Overview
		// Extract genre names
		session.TMDBGenres = make([]string, len(tv.Genres))
		for i, genre := range tv.Genres {
			session.TMDBGenres[i] = genre.Name
		}
	}

	// Redirect to step 4 (disk details)
	http.Redirect(w, r, "/import/step4?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
}

// ImportStep3Handler handles step 3: manual title/year entry
func (app *App) ImportStep3Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		if title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		session.Title = title

		// Year is required for films
		if session.MediaKind == Film {
			yearStr := r.FormValue("year")
			year, err := strconv.Atoi(yearStr)
			if err != nil || year <= 0 {
				http.Error(w, "Valid year is required for films", http.StatusBadRequest)
				return
			}
			session.Year = year
		}

		// Redirect to step 4 (disk details)
		http.Redirect(w, r, "/import/step4?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step3.html", data)
	if err != nil {
		log.Printf("Error rendering import_step3 template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportStep4Handler handles step 4: disk details (series/disk numbers, disk type)
func (app *App) ImportStep4Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Parse TV-specific fields
		if session.MediaKind == TV {
			seriesStr := r.FormValue("series_num")
			series, err := strconv.Atoi(seriesStr)
			if err != nil || series <= 0 {
				http.Error(w, "Valid series number is required for TV shows", http.StatusBadRequest)
				return
			}
			session.SeriesNum = series

			diskStr := r.FormValue("disk_num")
			disk, err := strconv.Atoi(diskStr)
			if err != nil || disk <= 0 {
				http.Error(w, "Valid disk number is required for TV shows", http.StatusBadRequest)
				return
			}
			session.DiskNum = disk
		} else {
			// For films, disk number defaults to 1
			session.DiskNum = 1
		}

		// Parse disk type
		diskTypeStr := r.FormValue("disk_type")
		switch diskTypeStr {
		case "bluray":
			session.DiskType = DiskTypeBluRay
		case "bluray_uhd":
			session.DiskType = DiskTypeBluRayUHD
		case "dvd":
			session.DiskType = DiskTypeDVD
		case "custom":
			customType := r.FormValue("disk_type_custom")
			if customType == "" {
				http.Error(w, "Custom disk type text is required", http.StatusBadRequest)
				return
			}
			session.DiskType = DiskTypeCustom
			session.DiskTypeCustom = customType
		default:
			http.Error(w, "Invalid disk type", http.StatusBadRequest)
			return
		}

		// Redirect to step 5 (add to existing or create new)
		http.Redirect(w, r, "/import/step5?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step4.html", data)
	if err != nil {
		log.Printf("Error rendering import_step4 template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportStep5Handler handles step 5: add to existing or create new media
func (app *App) ImportStep5Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		action := r.FormValue("action")
		if action == "new" {
			session.AddToExisting = false
		} else if action == "existing" {
			session.AddToExisting = true
			existingSlug := r.FormValue("existing_media")
			if existingSlug == "" {
				http.Error(w, "Existing media selection is required", http.StatusBadRequest)
				return
			}

			// Find media by slug
			media := app.findMediaBySlug(existingSlug)
			if media == nil {
				http.Error(w, "Selected media not found", http.StatusNotFound)
				return
			}

			session.ExistingMediaPath = media.Path
		} else {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		// Redirect to confirmation page
		http.Redirect(w, r, "/import/confirm?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	// Get compatible existing media (same type)
	var compatibleMedia []Media
	for _, media := range app.mediaList {
		if media.Type == session.MediaKind {
			compatibleMedia = append(compatibleMedia, media)
		}
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	data := struct {
		Session         *ImportSession
		SessionID       string
		CompatibleMedia []Media
	}{
		Session:         session,
		SessionID:       sessionID,
		CompatibleMedia: compatibleMedia,
	}

	err := tmpl.ExecuteTemplate(w, "import_step5.html", data)
	if err != nil {
		log.Printf("Error rendering import_step5 template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportConfirmHandler shows confirmation before executing import
func (app *App) ImportConfirmHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	// Generate preview of destination path
	finalTitle := session.Title
	if session.TMDBTitle != "" {
		finalTitle = session.TMDBTitle
	}
	finalYear := session.Year
	if session.TMDBYear > 0 {
		finalYear = session.TMDBYear
	}

	var destPath string
	if session.AddToExisting {
		diskTypeText := session.DiskType.String()
		if session.DiskType == DiskTypeCustom {
			diskTypeText = session.DiskTypeCustom
		}
		diskDir := GenerateDiskDirName(diskTypeText, session.SeriesNum, session.DiskNum, session.MediaKind)
		destPath = session.ExistingMediaPath + "/" + diskDir
	} else {
		mediaDir := GenerateMediaDirName(finalTitle, finalYear, session.MediaKind)
		diskTypeText := session.DiskType.String()
		if session.DiskType == DiskTypeCustom {
			diskTypeText = session.DiskTypeCustom
		}
		diskDir := GenerateDiskDirName(diskTypeText, session.SeriesNum, session.DiskNum, session.MediaKind)
		destPath = app.mediaDir + "/" + mediaDir + "/" + diskDir
	}

	data := struct {
		Session  *ImportSession
		SessionID string
		DestPath string
	}{
		Session:  session,
		SessionID: sessionID,
		DestPath: destPath,
	}

	err := tmpl.ExecuteTemplate(w, "import_confirm.html", data)
	if err != nil {
		log.Printf("Error rendering import_confirm template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// ImportExecuteHandler executes the import
func (app *App) ImportExecuteHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.FormValue("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	session, ok := importSessionStore.Get(sessionID)
	if !ok {
		http.Error(w, "Invalid session", http.StatusNotFound)
		return
	}

	// Execute the import
	err := ExecuteImport(session, app.mediaDir)
	if err != nil {
		log.Printf("Import failed: %v", err)
		http.Error(w, fmt.Sprintf("Import failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch and save TMDB metadata if available
	if session.TMDBID != "" && app.tmdbClient != nil {
		// Determine the media path
		var mediaPath string
		if session.AddToExisting {
			mediaPath = session.ExistingMediaPath
		} else {
			finalTitle := session.Title
			if session.TMDBTitle != "" {
				finalTitle = session.TMDBTitle
			}
			finalYear := session.Year
			if session.TMDBYear > 0 {
				finalYear = session.TMDBYear
			}
			mediaDir := GenerateMediaDirName(finalTitle, finalYear, session.MediaKind)
			mediaPath = app.mediaDir + "/" + mediaDir
		}

		// Create a temporary Media object for metadata fetching
		media := &Media{
			Type:   session.MediaKind,
			TMDBID: session.TMDBID,
			Path:   mediaPath,
		}

		err = app.tmdbClient.FetchAndSaveMetadata(media)
		if err != nil {
			log.Printf("Warning: Failed to fetch metadata: %v", err)
			// Don't fail the import, just log the warning
		}
	}

	// Clean up session
	importSessionStore.Delete(sessionID)

	// Redirect to success page
	http.Redirect(w, r, "/import/success", http.StatusSeeOther)
}

// ImportSuccessHandler shows the import success page
func (app *App) ImportSuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Reload templates in dev mode
	tmpl := app.templates
	if app.devMode {
		tmpl = app.loadTemplates()
	}

	err := tmpl.ExecuteTemplate(w, "import_success.html", nil)
	if err != nil {
		log.Printf("Error rendering import_success template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
