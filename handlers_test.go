package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create a simple template for testing
	tmpl := template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html>
<body>
<h1>Media Backup Manager</h1>
<table>
{{range .MediaList}}
<tr>
<td>{{.DisplayTitle}}</td>
<td>{{.Type}}</td>
<td>{{.DiskCount}}</td>
<td>{{.TMDBID}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`))

	tests := []struct {
		name           string
		mediaList      []Media
		expectedStatus int
		expectedInBody []string
	}{
		{
			name: "Multiple media items",
			mediaList: []Media{
				{
					Title:     "War of the Worlds",
					Type:      Film,
					Year:      2025,
					DiskCount: 1,
					TMDBID:    "755898",
					Path:      "/test/path1",
				},
				{
					Title:     "Better Call Saul",
					Type:      TV,
					Year:      0,
					DiskCount: 5,
					TMDBID:    "60059",
					Path:      "/test/path2",
				},
			},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Media Backup Manager",
				"War of the Worlds (2025)",
				"Better Call Saul",
				"Film",
				"TV",
				"755898",
				"60059",
			},
		},
		{
			name:           "Empty media list",
			mediaList:      []Media{},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"Media Backup Manager"},
		},
		{
			name: "Film without TMDB ID",
			mediaList: []Media{
				{
					Title:     "Unknown Film",
					Type:      Film,
					Year:      2020,
					DiskCount: 1,
					TMDBID:    "",
					Path:      "/test/path",
				},
			},
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Unknown Film (2020)",
				"Film",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("/test/media")
			app := NewApp(tt.mediaList, scanner, tmpl, "/test/media", "")

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			app.IndexHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("IndexHandler() status = %v, want %v", res.StatusCode, tt.expectedStatus)
			}

			body := w.Body.String()
			for _, expected := range tt.expectedInBody {
				if !strings.Contains(body, expected) {
					t.Errorf("IndexHandler() body does not contain %q", expected)
				}
			}
		})
	}
}

func TestIndexHandlerSorting(t *testing.T) {
	// Create a simple template for testing
	tmpl := template.Must(template.New("index.html").Parse(`{{range .MediaList}}{{.Title}},{{end}}`))

	mediaList := []Media{
		{Title: "Zebra Show", Type: TV, DiskCount: 1, Path: "/test/z"},
		{Title: "Alpha Film", Type: Film, Year: 2020, DiskCount: 1, Path: "/test/a"},
		{Title: "Beta Show", Type: TV, DiskCount: 1, Path: "/test/b"},
		{Title: "Gamma Film", Type: Film, Year: 2021, DiskCount: 1, Path: "/test/g"},
	}

	scanner := NewScanner("/test/media")
	app := NewApp(mediaList, scanner, tmpl, "/test/media", "")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	app.IndexHandler(w, req)

	body := w.Body.String()

	// Expected order: Films first (alphabetically), then TV shows (alphabetically)
	// Alpha Film, Gamma Film, Beta Show, Zebra Show
	expected := "Alpha Film,Gamma Film,Beta Show,Zebra Show,"
	if body != expected {
		t.Errorf("IndexHandler() sorting = %q, want %q", body, expected)
	}
}

func TestNewApp(t *testing.T) {
	mediaList := []Media{
		{Title: "Test", Type: Film, Year: 2020, DiskCount: 1, Path: "/test"},
	}
	tmpl := template.Must(template.New("test").Parse("test"))
	scanner := NewScanner("/test/media")

	app := NewApp(mediaList, scanner, tmpl, "/test/media", "")

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	if len(app.mediaList) != 1 {
		t.Errorf("NewApp() mediaList length = %v, want 1", len(app.mediaList))
	}
	if app.templates == nil {
		t.Error("NewApp() templates is nil")
	}
	if app.scanner == nil {
		t.Error("NewApp() scanner is nil")
	}
}

func TestDetailHandler(t *testing.T) {
	// Create a realistic detail template
	tmpl := template.Must(template.New("detail.html").Parse(`
<!DOCTYPE html>
<html>
<body>
<h1>{{.Media.DisplayTitle}}</h1>
<p>Type: {{.Media.Type}}</p>
<p>Disks: {{.Media.DiskCount}}</p>
{{if .Media.Disks}}
<table>
<tr><th>Name</th><th>Format</th><th>Size</th></tr>
{{range .Media.Disks}}
<tr><td>{{.Name}}</td><td>{{.Format}}</td><td>{{printf "%.1f GB" .SizeGB}}</td></tr>
{{end}}
</table>
{{end}}
<p>Description: {{.Description}}</p>
{{range .Genres}}<span>{{.}}</span>{{end}}
</body>
</html>
`))

	tests := []struct {
		name           string
		mediaList      []Media
		requestPath    string
		expectedStatus int
		expectedInBody []string
		notInBody      []string
	}{
		{
			name: "Film with disks",
			mediaList: []Media{
				{
					Title:     "The Thing",
					Type:      Film,
					Year:      1982,
					DiskCount: 2,
					Disks: []Disk{
						{Name: "Disk 1", Format: "Blu-Ray", SizeGB: 45.2},
						{Name: "Disk 2", Format: "DVD", SizeGB: 4.7},
					},
					TMDBID: "1091",
					Path:   "/test/the-thing",
				},
			},
			requestPath:    "/media/the-thing-1982",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"The Thing (1982)",
				"Film",
				"Disk 1",
				"Disk 2",
				"Blu-Ray",
				"DVD",
				"45.2 GB",
				"4.7 GB",
			},
		},
		{
			name: "TV show with disks",
			mediaList: []Media{
				{
					Title:     "Better Call Saul",
					Type:      TV,
					Year:      0,
					DiskCount: 2,
					Disks: []Disk{
						{Name: "Series 1 Disk 1", Format: "Blu-Ray", SizeGB: 23.5},
						{Name: "Series 1 Disk 2", Format: "Blu-Ray UHD", SizeGB: 66.8},
					},
					TMDBID: "60059",
					Path:   "/test/better-call-saul",
				},
			},
			requestPath:    "/media/better-call-saul",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Better Call Saul",
				"TV",
				"Series 1 Disk 1",
				"Series 1 Disk 2",
				"Blu-Ray UHD",
				"23.5 GB",
				"66.8 GB",
			},
		},
		{
			name: "Film with no disks",
			mediaList: []Media{
				{
					Title:     "Empty Film",
					Type:      Film,
					Year:      2020,
					DiskCount: 0,
					Disks:     []Disk{},
					Path:      "/test/empty-film",
				},
			},
			requestPath:    "/media/empty-film-2020",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				"Empty Film (2020)",
				"Film",
			},
			notInBody: []string{
				"<table>",
			},
		},
		{
			name:           "Invalid slug - not found",
			mediaList:      []Media{},
			requestPath:    "/media/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty slug",
			mediaList:      []Media{},
			requestPath:    "/media/",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("/test/media")
			app := NewApp(tt.mediaList, scanner, tmpl, "/test/media", "")

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			w := httptest.NewRecorder()

			app.DetailHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("DetailHandler() status = %v, want %v", res.StatusCode, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				body := w.Body.String()
				for _, expected := range tt.expectedInBody {
					if !strings.Contains(body, expected) {
						t.Errorf("DetailHandler() body does not contain %q", expected)
					}
				}
				for _, notExpected := range tt.notInBody {
					if strings.Contains(body, notExpected) {
						t.Errorf("DetailHandler() body should not contain %q", notExpected)
					}
				}
			}
		})
	}
}

func TestFindMediaBySlug(t *testing.T) {
	mediaList := []Media{
		{Title: "The Thing", Type: Film, Year: 1982, Path: "/test/thing"},
		{Title: "Better Call Saul", Type: TV, Year: 0, Path: "/test/bcs"},
	}

	tmpl := template.Must(template.New("test").Parse("test"))
	scanner := NewScanner("/test/media")
	app := NewApp(mediaList, scanner, tmpl, "/test/media", "")

	tests := []struct {
		name      string
		slug      string
		wantTitle string
		wantNil   bool
	}{
		{
			name:      "Find film by slug",
			slug:      "the-thing-1982",
			wantTitle: "The Thing",
			wantNil:   false,
		},
		{
			name:      "Find TV show by slug",
			slug:      "better-call-saul",
			wantTitle: "Better Call Saul",
			wantNil:   false,
		},
		{
			name:    "Nonexistent slug",
			slug:    "nonexistent",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := app.findMediaBySlug(tt.slug)

			if tt.wantNil {
				if media != nil {
					t.Errorf("findMediaBySlug(%q) = %v, want nil", tt.slug, media)
				}
			} else {
				if media == nil {
					t.Errorf("findMediaBySlug(%q) = nil, want media", tt.slug)
				} else if media.Title != tt.wantTitle {
					t.Errorf("findMediaBySlug(%q).Title = %q, want %q", tt.slug, media.Title, tt.wantTitle)
				}
			}
		})
	}
}

<<<<<<< HEAD
func TestRefreshMediaList(t *testing.T) {
	// Create a temporary test directory
	testDir := t.TempDir()

	// Set up test data with one film
	setupTestMediaDir(t, testDir, []string{
		"The Matrix (1999) [Film]/Disk [Blu-Ray]",
	})

	// Create scanner and initial app
	scanner := NewScanner(testDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Initial scan failed: %v", err)
	}

	tmpl := template.Must(template.New("test").Parse("test"))
	app := NewApp(mediaList, scanner, tmpl, testDir, "")

	// Verify initial state
	if len(app.mediaList) != 1 {
		t.Fatalf("Initial media list length = %d, want 1", len(app.mediaList))
	}
	if app.mediaList[0].Title != "The Matrix" {
		t.Errorf("Initial media title = %q, want %q", app.mediaList[0].Title, "The Matrix")
	}

	// Add a new media item to the directory
	setupTestMediaDir(t, testDir, []string{
		"Fight Club (1999) [Film]/Disk [DVD]",
	})

	// Refresh the media list
	err = app.RefreshMediaList()
	if err != nil {
		t.Fatalf("RefreshMediaList() failed: %v", err)
	}

	// Verify the list was refreshed
	app.mediaListMutex.RLock()
	mediaCount := len(app.mediaList)
	app.mediaListMutex.RUnlock()

	if mediaCount != 2 {
		t.Errorf("After refresh, media list length = %d, want 2", mediaCount)
	}

	// Verify both media items are present
	titles := make(map[string]bool)
	app.mediaListMutex.RLock()
	for _, media := range app.mediaList {
		titles[media.Title] = true
	}
	app.mediaListMutex.RUnlock()

	if !titles["The Matrix"] {
		t.Error("The Matrix not found after refresh")
	}
	if !titles["Fight Club"] {
		t.Error("Fight Club not found after refresh")
	}
}

func TestRefreshMediaListThreadSafety(t *testing.T) {
	// Create a temporary test directory
	testDir := t.TempDir()

	// Set up test data
	setupTestMediaDir(t, testDir, []string{
		"The Matrix (1999) [Film]/Disk [Blu-Ray]",
		"Fight Club (1999) [Film]/Disk [DVD]",
	})

	// Create scanner and app
	scanner := NewScanner(testDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Initial scan failed: %v", err)
	}

	tmpl := template.Must(template.New("test").Parse("test"))
	app := NewApp(mediaList, scanner, tmpl, testDir, "")

	// Run concurrent operations to test thread safety
	done := make(chan bool)
	errors := make(chan error, 3)

	// Goroutine 1: Repeatedly refresh media list
	go func() {
		for i := 0; i < 10; i++ {
			if err := app.RefreshMediaList(); err != nil {
				errors <- err
				return
			}
		}
		done <- true
	}()

	// Goroutine 2: Repeatedly read media list via findMediaBySlug
	go func() {
		for i := 0; i < 100; i++ {
			_ = app.findMediaBySlug("the-matrix-1999")
		}
		done <- true
	}()

	// Goroutine 3: Repeatedly iterate over media list
	go func() {
		for i := 0; i < 100; i++ {
			app.mediaListMutex.RLock()
			count := len(app.mediaList)
			app.mediaListMutex.RUnlock()
			if count < 0 {
				errors <- fmt.Errorf("invalid media count: %d", count)
				return
			}
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		select {
		case <-done:
			// Success
		case err := <-errors:
			t.Fatalf("Concurrent operation failed: %v", err)
		}
	}
}

// Helper function to set up test media directories
func setupTestMediaDir(t *testing.T, baseDir string, paths []string) {
	t.Helper()
	for _, path := range paths {
		fullPath := filepath.Join(baseDir, path)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", fullPath, err)
		}
=======
func TestSelectPosterHandler(t *testing.T) {
	// Create a simple template
	tmpl := template.Must(template.New("select_poster.html").Parse(`Poster selection page`))

	// Test without TMDB client
	mediaList := []Media{
		{Title: "Test Film", Type: Film, Year: 2020, TMDBID: "123", Path: "/test/path", DiskCount: 1},
	}

	app := NewApp(mediaList, tmpl, "/test/media", "")

	req := httptest.NewRequest(http.MethodGet, "/media/test-film-2020/select-poster", nil)
	w := httptest.NewRecorder()

	app.SelectPosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Should return 503 when TMDB client is not available
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("SelectPosterHandler() without TMDB client status = %v, want %v", res.StatusCode, http.StatusServiceUnavailable)
	}
}

func TestSelectPosterHandlerNoTMDBID(t *testing.T) {
	// Create a simple template
	tmpl := template.Must(template.New("select_poster.html").Parse(`Poster selection page`))

	// Media without TMDB ID
	mediaList := []Media{
		{Title: "Test Film", Type: Film, Year: 2020, TMDBID: "", Path: "/test/path", DiskCount: 1},
	}

	app := NewApp(mediaList, tmpl, "/test/media", "")
	client := NewTMDBClient("test-key")
	app.SetTMDBClient(client)

	req := httptest.NewRequest(http.MethodGet, "/media/test-film-2020/select-poster", nil)
	w := httptest.NewRecorder()

	app.SelectPosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Should return 400 when media has no TMDB ID
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("SelectPosterHandler() without TMDB ID status = %v, want %v", res.StatusCode, http.StatusBadRequest)
	}
}

func TestSelectPosterHandlerNotFound(t *testing.T) {
	tmpl := template.Must(template.New("select_poster.html").Parse(`Poster selection page`))

	mediaList := []Media{}
	app := NewApp(mediaList, tmpl, "/test/media", "")
	client := NewTMDBClient("test-key")
	app.SetTMDBClient(client)

	req := httptest.NewRequest(http.MethodGet, "/media/nonexistent/select-poster", nil)
	w := httptest.NewRecorder()

	app.SelectPosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("SelectPosterHandler() for nonexistent media status = %v, want %v", res.StatusCode, http.StatusNotFound)
	}
}

func TestSavePosterHandlerMethodNotAllowed(t *testing.T) {
	mediaList := []Media{
		{Title: "Test Film", Type: Film, Year: 2020, TMDBID: "123", Path: "/test/path", DiskCount: 1},
	}

	app := NewApp(mediaList, nil, "/test/media", "")
	client := NewTMDBClient("test-key")
	app.SetTMDBClient(client)

	// Test with GET (should only accept POST)
	req := httptest.NewRequest(http.MethodGet, "/media/test-film-2020/save-poster", nil)
	w := httptest.NewRecorder()

	app.SavePosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("SavePosterHandler() with GET status = %v, want %v", res.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestSavePosterHandlerNoTMDBClient(t *testing.T) {
	mediaList := []Media{
		{Title: "Test Film", Type: Film, Year: 2020, TMDBID: "123", Path: "/test/path", DiskCount: 1},
	}

	app := NewApp(mediaList, nil, "/test/media", "")

	req := httptest.NewRequest(http.MethodPost, "/media/test-film-2020/save-poster", nil)
	w := httptest.NewRecorder()

	app.SavePosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("SavePosterHandler() without TMDB client status = %v, want %v", res.StatusCode, http.StatusServiceUnavailable)
	}
}

func TestSavePosterHandlerMissingPosterPath(t *testing.T) {
	mediaList := []Media{
		{Title: "Test Film", Type: Film, Year: 2020, TMDBID: "123", Path: "/test/path", DiskCount: 1},
	}

	app := NewApp(mediaList, nil, "/test/media", "")
	client := NewTMDBClient("test-key")
	app.SetTMDBClient(client)

	// POST without poster_path form value
	req := httptest.NewRequest(http.MethodPost, "/media/test-film-2020/save-poster", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	app.SavePosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("SavePosterHandler() without poster_path status = %v, want %v", res.StatusCode, http.StatusBadRequest)
	}
}

func TestSavePosterHandlerNotFound(t *testing.T) {
	mediaList := []Media{}
	app := NewApp(mediaList, nil, "/test/media", "")
	client := NewTMDBClient("test-key")
	app.SetTMDBClient(client)

	req := httptest.NewRequest(http.MethodPost, "/media/nonexistent/save-poster",
		strings.NewReader("poster_path=/test.jpg"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	app.SavePosterHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("SavePosterHandler() for nonexistent media status = %v, want %v", res.StatusCode, http.StatusNotFound)
>>>>>>> b3b3297 (add poster selector functionality)
	}
}
