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

// ImportListHandler shows the list of directories available for import
func (app *App) ImportListHandler(w http.ResponseWriter, r *http.Request) {
	// Check if import is enabled
	if app.importScanner == nil {
		writeError(w, http.StatusServiceUnavailable, "Import functionality is not configured (IMPORT_DIR not set)", nil)
		return
	}

	tmpl := app.getTemplates()

	// Scan import directory
	imports, err := app.importScanner.Scan()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to scan import directory", err)
		return
	}

	data := struct {
		Imports []ImportDirectory
	}{
		Imports: imports,
	}

	err = tmpl.ExecuteTemplate(w, "import_list.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportStartHandler starts an import session for a selected directory
func (app *App) ImportStartHandler(w http.ResponseWriter, r *http.Request) {
	// Check if import is enabled
	if app.importScanner == nil {
		writeError(w, http.StatusServiceUnavailable, "Import functionality is not configured", nil)
		return
	}

	// Get directory name from query parameter
	dirName := r.URL.Query().Get("dir")
	if dirName == "" {
		writeError(w, http.StatusBadRequest, "Directory name is required", nil)
		return
	}

	// Scan to find the directory
	imports, err := app.importScanner.Scan()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to scan import directory", err)
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
	sessionID := app.importSessions.Create(session)

	// Redirect to step 1 (choose media kind)
	http.Redirect(w, r, "/import/step1?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
}

// ImportStep1Handler handles step 1: choose media kind (Film/TV)
func (app *App) ImportStep1Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			writeError(w, http.StatusBadRequest, "Failed to parse form", err)
			return
		}

		mediaKindStr := r.FormValue("media_kind")
		if mediaKindStr == "film" {
			session.MediaKind = Film
		} else if mediaKindStr == "tv" {
			session.MediaKind = TV
		} else if mediaKindStr == "music" {
			session.MediaKind = Music
		} else {
			writeError(w, http.StatusBadRequest, "Invalid media kind", nil)
			return
		}

		// Redirect to step 2 (add to existing or create new)
		http.Redirect(w, r, "/import/step2?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	tmpl := app.getTemplates()

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step1.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportStep2Handler handles step 2: add to existing or create new media
func (app *App) ImportStep2Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			writeError(w, http.StatusBadRequest, "Failed to parse form", err)
			return
		}

		action := r.FormValue("action")
		if action == "new" {
			session.AddToExisting = false
			// Redirect to TMDB search/manual entry (step 3)
			http.Redirect(w, r, "/import/step3?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
			return
		} else if action == "existing" {
			session.AddToExisting = true
			existingSlug := r.FormValue("existing_media")
			if existingSlug == "" {
				writeError(w, http.StatusBadRequest, "Existing media selection is required", nil)
				return
			}

			// Find media by slug
			media := app.findMediaBySlug(existingSlug)
			if media == nil {
				writeError(w, http.StatusNotFound, "Selected media not found", nil)
				return
			}

			session.ExistingMediaPath = media.Path
			// Skip TMDB search and go to disk details (step 4)
			http.Redirect(w, r, "/import/step4?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
			return
		} else {
			writeError(w, http.StatusBadRequest, "Invalid action", nil)
			return
		}
	}

	// Get compatible existing media (same type)
	app.mediaListMutex.RLock()
	var compatibleMedia []Media
	for _, media := range app.mediaList {
		if media.Type == session.MediaKind {
			compatibleMedia = append(compatibleMedia, media)
		}
	}
	app.mediaListMutex.RUnlock()

	tmpl := app.getTemplates()

	data := struct {
		Session         *ImportSession
		SessionID       string
		CompatibleMedia []Media
	}{
		Session:         session,
		SessionID:       sessionID,
		CompatibleMedia: compatibleMedia,
	}

	err := tmpl.ExecuteTemplate(w, "import_step2.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportStep3Handler handles step 3: TMDB search or manual entry (only for new media)
func (app *App) ImportStep3Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Handle "skip TMDB" action
	if r.Method == http.MethodPost && r.FormValue("action") == "skip" {
		// Redirect to manual entry (step 3b)
		http.Redirect(w, r, "/import/step3/manual?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
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
	if query != "" {
		if session.MediaKind == Film && app.tmdbClient != nil {
			movieResults, err := app.tmdbClient.SearchMovies(query, year)
			if err != nil {
				searchErr = err
			} else {
				results = movieResults
			}
		} else if session.MediaKind == TV && app.tmdbClient != nil {
			tvResults, err := app.tmdbClient.SearchTV(query)
			if err != nil {
				searchErr = err
			} else {
				results = tvResults
			}
		} else if session.MediaKind == Music && app.musicBrainzClient != nil {
			musicResults, err := app.musicBrainzClient.SearchReleases(query)
			if err != nil {
				searchErr = err
			} else {
				results = musicResults
			}
		}
	}

	// Prepare error message
	var errorMsg string
	if searchErr != nil {
		log.Printf("TMDB search failed for %q (kind=%s): %v", query, session.MediaKind, searchErr)
		errorMsg = "Search error: unable to reach TMDB"
	}

	tmpl := app.getTemplates()

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

	err := tmpl.ExecuteTemplate(w, "import_step3.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportStep3ConfirmHandler handles TMDB match selection
func (app *App) ImportStep3ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	tmdbID := r.URL.Query().Get("id")

	if sessionID == "" || tmdbID == "" {
		writeError(w, http.StatusBadRequest, "Session ID and TMDB ID are required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Fetch metadata from TMDB or MusicBrainz
	if session.MediaKind == Film {
		if app.tmdbClient == nil {
			writeError(w, http.StatusServiceUnavailable, "TMDB API is not configured", nil)
			return
		}
		movie, err := app.tmdbClient.FetchMovieMetadata(tmdbID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to fetch movie metadata", err)
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
		if app.tmdbClient == nil {
			writeError(w, http.StatusServiceUnavailable, "TMDB API is not configured", nil)
			return
		}
		tv, err := app.tmdbClient.FetchTVMetadata(tmdbID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to fetch TV metadata", err)
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
	} else if session.MediaKind == Music {
		if app.musicBrainzClient == nil {
			writeError(w, http.StatusServiceUnavailable, "MusicBrainz client is not configured", nil)
			return
		}
		// For Music, the ID passed is the MusicBrainz Release ID
		// We fetch the release to get details
		release, err := app.musicBrainzClient.FetchRelease(tmdbID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to fetch release metadata", err)
			return
		}
		session.MusicBrainzID = tmdbID
		session.MusicBrainzTitle = release.Title
		// We don't have artist in the release struct yet, but we can get it from search or add it to release
		// For now, let's just use the title
	}

	// Redirect to step 4 (disk details)
	http.Redirect(w, r, "/import/step4?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
}

// ImportStep3ManualHandler handles step 3b: manual title/year entry (when TMDB is skipped)
func (app *App) ImportStep3ManualHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			writeError(w, http.StatusBadRequest, "Failed to parse form", err)
			return
		}

		title := r.FormValue("title")
		if title == "" {
			writeError(w, http.StatusBadRequest, "Title is required", nil)
			return
		}
		session.Title = title

		// Year is required for films
		if session.MediaKind == Film {
			yearStr := r.FormValue("year")
			year, err := strconv.Atoi(yearStr)
			if err != nil || year <= 0 {
				writeError(w, http.StatusBadRequest, "Valid year is required for films", nil)
				return
			}
			session.Year = year
		}

		// Redirect to step 4 (disk details)
		http.Redirect(w, r, "/import/step4?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	tmpl := app.getTemplates()

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step3_manual.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportStep4Handler handles step 4: disk details (series/disk numbers, disk type)
func (app *App) ImportStep4Handler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			writeError(w, http.StatusBadRequest, "Failed to parse form", err)
			return
		}

		// Parse TV-specific fields
		if session.MediaKind == TV {
			seriesStr := r.FormValue("series_num")
			series, err := strconv.Atoi(seriesStr)
			if err != nil || series <= 0 {
				writeError(w, http.StatusBadRequest, "Valid series number is required for TV shows", nil)
				return
			}
			session.SeriesNum = series

			diskStr := r.FormValue("disk_num")
			disk, err := strconv.Atoi(diskStr)
			if err != nil || disk <= 0 {
				writeError(w, http.StatusBadRequest, "Valid disk number is required for TV shows", nil)
				return
			}
			session.DiskNum = disk
		} else {
			// For films and music, disk number defaults to 1
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
				writeError(w, http.StatusBadRequest, "Custom disk type text is required", nil)
				return
			}
			session.DiskType = DiskTypeCustom
			session.DiskTypeCustom = customType
		default:
			writeError(w, http.StatusBadRequest, "Invalid disk type", nil)
			return
		}

		// Redirect to confirmation page
		http.Redirect(w, r, "/import/confirm?session="+url.QueryEscape(sessionID), http.StatusSeeOther)
		return
	}

	tmpl := app.getTemplates()

	data := struct {
		Session   *ImportSession
		SessionID string
	}{
		Session:   session,
		SessionID: sessionID,
	}

	err := tmpl.ExecuteTemplate(w, "import_step4.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportConfirmHandler shows confirmation before executing import
func (app *App) ImportConfirmHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	tmpl := app.getTemplates()

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
		Session   *ImportSession
		SessionID string
		DestPath  string
	}{
		Session:   session,
		SessionID: sessionID,
		DestPath:  destPath,
	}

	err := tmpl.ExecuteTemplate(w, "import_confirm.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}

// ImportExecuteHandler executes the import
func (app *App) ImportExecuteHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	sessionID := r.FormValue("session")
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, ok := app.importSessions.Get(sessionID)
	if !ok {
		writeError(w, http.StatusNotFound, "Invalid session", nil)
		return
	}

	// Execute the import
	err := ExecuteImport(session, app.mediaDir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Import failed", err)
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

	// Refresh media list to include the newly imported media
	err = app.RefreshMediaList()
	if err != nil {
		log.Printf("Warning: Failed to refresh media list: %v", err)
		// Don't fail the import, just log the warning
	}

	// Store directory name before cleaning up session
	dirName := session.SourceDir.Name

	// Clean up session
	app.importSessions.Delete(sessionID)

	// Redirect to success page with directory name
	http.Redirect(w, r, "/import/success?dir="+url.QueryEscape(dirName), http.StatusSeeOther)
}

// ImportSuccessHandler shows the import success page
func (app *App) ImportSuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Get directory name from query parameter
	dirName := r.URL.Query().Get("dir")

	tmpl := app.getTemplates()

	data := struct {
		DirName string
	}{
		DirName: dirName,
	}

	err := tmpl.ExecuteTemplate(w, "import_success.html", data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error rendering template", err)
		return
	}
}
