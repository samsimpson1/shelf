package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationScanAndServe(t *testing.T) {
	// Use programmatically created test directory for integration testing
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Failed to scan testdata: %v", err)
	}

	// Load templates
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Create app
	app := NewApp(mediaList, tmpl, testDir)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.IndexHandler(w, req)

	// Check response
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}

	body := w.Body.String()

	// Verify expected content is present
	expectedContent := []string{
		"Shelf",
		"War of the Worlds",
		"Better Call Saul",
		"No TMDB",
		"grid", // Check for grid layout
	}

	for _, expected := range expectedContent {
		if !strings.Contains(body, expected) {
			t.Errorf("Response body missing expected content: %q", expected)
		}
	}
}

func TestIntegrationWithEmptyDirectory(t *testing.T) {
	// Create temporary empty directory
	tmpDir, err := os.MkdirTemp("", "empty-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Scan empty directory
	scanner := NewScanner(tmpDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Failed to scan empty directory: %v", err)
	}

	if len(mediaList) != 0 {
		t.Errorf("Expected empty media list, got %d items", len(mediaList))
	}

	// Load templates
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Create app and serve
	app := NewApp(mediaList, tmpl, tmpDir)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.IndexHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "No Media Found") {
		t.Error("Expected empty state message in response")
	}
}

func TestIntegrationWithCustomDirectory(t *testing.T) {
	// Create temporary directory with test media
	tmpDir, err := os.MkdirTemp("", "custom-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a film directory
	filmDir := filepath.Join(tmpDir, "Test Film (2020) [Film]")
	err = os.Mkdir(filmDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	diskDir := filepath.Join(filmDir, "Disk [Blu-Ray]")
	err = os.Mkdir(diskDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	tmdbFile := filepath.Join(filmDir, "tmdb.txt")
	err = os.WriteFile(tmdbFile, []byte("123456"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Scan custom directory
	scanner := NewScanner(tmpDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Failed to scan custom directory: %v", err)
	}

	if len(mediaList) != 1 {
		t.Errorf("Expected 1 media item, got %d", len(mediaList))
	}

	if len(mediaList) > 0 {
		media := mediaList[0]
		if media.Title != "Test Film" {
			t.Errorf("Expected title 'Test Film', got %q", media.Title)
		}
		if media.Year != 2020 {
			t.Errorf("Expected year 2020, got %d", media.Year)
		}
		if media.DiskCount != 1 {
			t.Errorf("Expected 1 disk, got %d", media.DiskCount)
		}
		if media.TMDBID != "123456" {
			t.Errorf("Expected TMDB ID '123456', got %q", media.TMDBID)
		}
	}
}

func TestIntegrationInvalidDirectory(t *testing.T) {
	scanner := NewScanner("/nonexistent/directory/path")
	_, err := scanner.Scan()
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}

func TestIntegrationMixedMedia(t *testing.T) {
	// Create temporary directory with mixed media types
	tmpDir, err := os.MkdirTemp("", "mixed-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create 2 films
	for i := 1; i <= 2; i++ {
		filmDir := filepath.Join(tmpDir, "Film "+string(rune('A'+i-1))+" (2020) [Film]")
		os.Mkdir(filmDir, 0755)
		os.Mkdir(filepath.Join(filmDir, "Disk [DVD]"), 0755)
	}

	// Create 2 TV shows
	for i := 1; i <= 2; i++ {
		tvDir := filepath.Join(tmpDir, "TV Show "+string(rune('A'+i-1))+" [TV]")
		os.Mkdir(tvDir, 0755)
		os.Mkdir(filepath.Join(tvDir, "Series 1 Disk 1 [DVD]"), 0755)
	}

	// Scan
	scanner := NewScanner(tmpDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Failed to scan mixed directory: %v", err)
	}

	if len(mediaList) != 4 {
		t.Errorf("Expected 4 media items, got %d", len(mediaList))
	}

	// Count by type
	filmCount := 0
	tvCount := 0
	for _, media := range mediaList {
		if media.Type == Film {
			filmCount++
		} else if media.Type == TV {
			tvCount++
		}
	}

	if filmCount != 2 {
		t.Errorf("Expected 2 films, got %d", filmCount)
	}
	if tvCount != 2 {
		t.Errorf("Expected 2 TV shows, got %d", tvCount)
	}
}

func TestIntegrationEndToEnd(t *testing.T) {
	// This test simulates the full application flow:
	// 1. Scan directory
	// 2. Load templates
	// 3. Create app
	// 4. Serve HTTP request
	// 5. Verify response

	// Scan testdata
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Verify scan results
	if len(mediaList) != 3 {
		t.Fatalf("Expected 3 media items, got %d", len(mediaList))
	}

	// Load templates
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		t.Fatalf("Template parsing failed: %v", err)
	}

	// Create app
	app := NewApp(mediaList, tmpl, testDir)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.IndexHandler)

	// Test server
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected HTML content type, got %q", contentType)
	}
}
