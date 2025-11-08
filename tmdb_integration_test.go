package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These tests make real API calls to TMDB and require a valid API key.
// Set TMDB_API_KEY environment variable to run these tests.
// Skip with: go test -short

// Helper function to find a poster file in a directory
func findPosterFileInDir(dir string) string {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		posterPath := filepath.Join(dir, "poster"+ext)
		if _, err := os.Stat(posterPath); err == nil {
			return posterPath
		}
	}
	return ""
}

func skipIfNoAPIKey(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: TMDB_API_KEY not set. Set it with: export TMDB_API_KEY=your_api_key_here")
	}
}

func TestIntegrationFetchMovieMetadata(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Use Fight Club (1999) - ID: 550
	metadata, err := client.FetchMovieMetadata("550")
	if err != nil {
		t.Fatalf("FetchMovieMetadata(550) failed: %v", err)
	}

	// Verify metadata fields
	if metadata.ID != 550 {
		t.Errorf("ID = %d, want 550", metadata.ID)
	}

	if metadata.Title != "Fight Club" {
		t.Errorf("Title = %q, want %q", metadata.Title, "Fight Club")
	}

	if metadata.PosterPath == "" {
		t.Error("PosterPath is empty")
	}

	if metadata.Overview == "" {
		t.Error("Overview is empty")
	}

	if len(metadata.Genres) == 0 {
		t.Error("Genres is empty")
	}

	// Verify release date
	if metadata.ReleaseDate == "" {
		t.Error("ReleaseDate is empty")
	}
}

func TestIntegrationFetchTVMetadata(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Use Better Call Saul - ID: 60059
	metadata, err := client.FetchTVMetadata("60059")
	if err != nil {
		t.Fatalf("FetchTVMetadata(60059) failed: %v", err)
	}

	// Verify metadata fields
	if metadata.ID != 60059 {
		t.Errorf("ID = %d, want 60059", metadata.ID)
	}

	if metadata.Name != "Better Call Saul" {
		t.Errorf("Name = %q, want %q", metadata.Name, "Better Call Saul")
	}

	if metadata.PosterPath == "" {
		t.Error("PosterPath is empty")
	}

	if metadata.Overview == "" {
		t.Error("Overview is empty")
	}

	if len(metadata.Genres) == 0 {
		t.Error("Genres is empty")
	}

	// Verify first air date
	if metadata.FirstAirDate == "" {
		t.Error("FirstAirDate is empty")
	}
}

func TestIntegrationSearchMovies(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Search for "The Matrix" (1999)
	results, err := client.SearchMovies("The Matrix", 1999)
	if err != nil {
		t.Fatalf("SearchMovies() failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("SearchMovies returned no results")
	}

	// Verify first result
	found := false
	for _, result := range results {
		if result.Title == "The Matrix" && strings.HasPrefix(result.ReleaseDate, "1999") {
			found = true
			if result.ID == 0 {
				t.Error("Result ID is 0")
			}
			if result.Overview == "" {
				t.Error("Result Overview is empty")
			}
			break
		}
	}

	if !found {
		t.Error("The Matrix (1999) not found in search results")
	}
}

func TestIntegrationSearchMoviesNoYear(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Search without year
	results, err := client.SearchMovies("Inception", 0)
	if err != nil {
		t.Fatalf("SearchMovies() without year failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("SearchMovies returned no results")
	}

	// Should find at least one result
	found := false
	for _, result := range results {
		if strings.Contains(result.Title, "Inception") {
			found = true
			break
		}
	}

	if !found {
		t.Error("No results containing 'Inception' found")
	}
}

func TestIntegrationSearchTV(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Search for "Breaking Bad"
	results, err := client.SearchTV("Breaking Bad")
	if err != nil {
		t.Fatalf("SearchTV() failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("SearchTV returned no results")
	}

	// Verify first result
	found := false
	for _, result := range results {
		if result.Name == "Breaking Bad" {
			found = true
			if result.ID == 0 {
				t.Error("Result ID is 0")
			}
			if result.Overview == "" {
				t.Error("Result Overview is empty")
			}
			break
		}
	}

	if !found {
		t.Error("Breaking Bad not found in search results")
	}
}

func TestIntegrationDownloadPoster(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "poster-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// First fetch metadata to get a poster path
	metadata, err := client.FetchMovieMetadata("550")
	if err != nil {
		t.Fatalf("FetchMovieMetadata failed: %v", err)
	}

	if metadata.PosterPath == "" {
		t.Fatal("No poster path in metadata")
	}

	// Download the poster
	err = client.DownloadPoster(metadata.PosterPath, tmpDir)
	if err != nil {
		t.Fatalf("DownloadPoster() failed: %v", err)
	}

	// Verify poster file was created
	posterFile := findPosterFileInDir(tmpDir)
	if posterFile == "" {
		t.Fatal("Poster file not found after download")
	}

	// Verify file exists and has content
	info, err := os.Stat(posterFile)
	if err != nil {
		t.Fatalf("Poster file stat failed: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Poster file is empty")
	}

	// Verify it's an image file (check extension)
	ext := strings.ToLower(filepath.Ext(posterFile))
	validExts := []string{".jpg", ".jpeg", ".png", ".webp"}
	valid := false
	for _, validExt := range validExts {
		if ext == validExt {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Poster file has invalid extension: %s", ext)
	}
}

func TestIntegrationSaveDescription(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "description-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Fetch metadata
	metadata, err := client.FetchMovieMetadata("550")
	if err != nil {
		t.Fatalf("FetchMovieMetadata failed: %v", err)
	}

	// Save description
	err = client.saveDescription(metadata.Overview, tmpDir)
	if err != nil {
		t.Fatalf("saveDescription() failed: %v", err)
	}

	// Verify file was created
	descPath := filepath.Join(tmpDir, "description.txt")
	data, err := os.ReadFile(descPath)
	if err != nil {
		t.Fatalf("Failed to read description file: %v", err)
	}

	if string(data) != metadata.Overview {
		t.Error("Description file content doesn't match metadata overview")
	}
}

func TestIntegrationSaveGenres(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "genres-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Fetch metadata
	metadata, err := client.FetchMovieMetadata("550")
	if err != nil {
		t.Fatalf("FetchMovieMetadata failed: %v", err)
	}

	// Save genres
	err = client.saveGenres(metadata.Genres, tmpDir)
	if err != nil {
		t.Fatalf("saveGenres() failed: %v", err)
	}

	// Verify file was created
	genrePath := filepath.Join(tmpDir, "genre.txt")
	data, err := os.ReadFile(genrePath)
	if err != nil {
		t.Fatalf("Failed to read genre file: %v", err)
	}

	// Verify format (comma-separated genre names)
	content := string(data)
	if content == "" {
		t.Error("Genre file is empty")
	}

	// Should contain at least one genre
	genres := strings.Split(content, ", ")
	if len(genres) == 0 {
		t.Error("No genres found in file")
	}

	// Verify genres match metadata
	for _, genre := range metadata.Genres {
		if !strings.Contains(content, genre.Name) {
			t.Errorf("Genre %q not found in saved file", genre.Name)
		}
	}
}

func TestIntegrationFetchAndSaveMetadata(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "metadata-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a media item
	media := &Media{
		Title:     "Fight Club",
		Type:      Film,
		Year:      1999,
		TMDBID:    "550",
		Path:      tmpDir,
		DiskCount: 1,
	}

	// Fetch and save all metadata
	err = client.FetchAndSaveMetadata(media)
	if err != nil {
		t.Fatalf("FetchAndSaveMetadata() failed: %v", err)
	}

	// Verify poster was saved
	posterFile, found := media.FindPosterFile()
	if !found {
		t.Error("Poster file not created")
	} else {
		info, err := os.Stat(posterFile)
		if err != nil {
			t.Errorf("Poster file error: %v", err)
		} else if info.Size() == 0 {
			t.Error("Poster file is empty")
		}
	}

	// Verify description was saved
	descPath := filepath.Join(tmpDir, "description.txt")
	descData, err := os.ReadFile(descPath)
	if err != nil {
		t.Errorf("Description file not created: %v", err)
	} else if len(descData) == 0 {
		t.Error("Description file is empty")
	}

	// Verify genres were saved
	genrePath := filepath.Join(tmpDir, "genre.txt")
	genreData, err := os.ReadFile(genrePath)
	if err != nil {
		t.Errorf("Genre file not created: %v", err)
	} else if len(genreData) == 0 {
		t.Error("Genre file is empty")
	}
}

func TestIntegrationFetchAndSaveMetadataTV(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "metadata-tv-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a TV media item
	media := &Media{
		Title:     "Better Call Saul",
		Type:      TV,
		TMDBID:    "60059",
		Path:      tmpDir,
		DiskCount: 5,
	}

	// Fetch and save all metadata
	err = client.FetchAndSaveMetadata(media)
	if err != nil {
		t.Fatalf("FetchAndSaveMetadata() for TV failed: %v", err)
	}

	// Verify all metadata files were created
	posterFile, found := media.FindPosterFile()
	if !found {
		t.Error("Poster file not created for TV show")
	} else if posterFile == "" {
		t.Error("Poster file path is empty")
	}

	descPath := filepath.Join(tmpDir, "description.txt")
	if _, err := os.Stat(descPath); err != nil {
		t.Error("Description file not created for TV show")
	}

	genrePath := filepath.Join(tmpDir, "genre.txt")
	if _, err := os.Stat(genrePath); err != nil {
		t.Error("Genre file not created for TV show")
	}
}

func TestIntegrationValidateTMDBIDMovie(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Valid movie ID
	err := client.ValidateTMDBID("550", Film)
	if err != nil {
		t.Errorf("ValidateTMDBID(550, Film) failed: %v", err)
	}

	// Invalid movie ID
	err = client.ValidateTMDBID("999999999", Film)
	if err == nil {
		t.Error("Expected error for invalid movie ID")
	}
}

func TestIntegrationValidateTMDBIDTV(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Valid TV ID
	err := client.ValidateTMDBID("60059", TV)
	if err != nil {
		t.Errorf("ValidateTMDBID(60059, TV) failed: %v", err)
	}

	// Invalid TV ID
	err = client.ValidateTMDBID("999999999", TV)
	if err == nil {
		t.Error("Expected error for invalid TV ID")
	}
}

func TestIntegrationErrorHandlingInvalidMovieID(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Test with completely invalid ID
	_, err := client.FetchMovieMetadata("not-a-number")
	if err == nil {
		t.Error("Expected error for non-numeric movie ID")
	}

	// Test with very high ID that doesn't exist
	_, err = client.FetchMovieMetadata("999999999")
	if err == nil {
		t.Error("Expected error for non-existent movie ID")
	}
}

func TestIntegrationErrorHandlingInvalidTVID(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Test with completely invalid ID
	_, err := client.FetchTVMetadata("not-a-number")
	if err == nil {
		t.Error("Expected error for non-numeric TV ID")
	}

	// Test with very high ID that doesn't exist
	_, err = client.FetchTVMetadata("999999999")
	if err == nil {
		t.Error("Expected error for non-existent TV ID")
	}
}

func TestIntegrationFetchAndSaveMetadataSkipsExisting(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory with existing metadata
	tmpDir, err := os.MkdirTemp("", "skip-existing-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing files
	posterPath := filepath.Join(tmpDir, "poster.jpg")
	err = os.WriteFile(posterPath, []byte("existing-poster-data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	descPath := filepath.Join(tmpDir, "description.txt")
	err = os.WriteFile(descPath, []byte("existing-description"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	genrePath := filepath.Join(tmpDir, "genre.txt")
	err = os.WriteFile(genrePath, []byte("existing-genres"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		TMDBID:    "550",
		Path:      tmpDir,
		DiskCount: 1,
	}

	// This should not fetch anything since all files exist
	err = client.FetchAndSaveMetadata(media)
	if err != nil {
		t.Errorf("FetchAndSaveMetadata() failed with existing files: %v", err)
	}

	// Verify files weren't modified
	posterData, _ := os.ReadFile(posterPath)
	if string(posterData) != "existing-poster-data" {
		t.Error("Existing poster was modified")
	}

	descData, _ := os.ReadFile(descPath)
	if string(descData) != "existing-description" {
		t.Error("Existing description was modified")
	}

	genreData, _ := os.ReadFile(genrePath)
	if string(genreData) != "existing-genres" {
		t.Error("Existing genres were modified")
	}
}

func TestIntegrationSearchMoviesLimit(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Search for a common term that will have many results
	results, err := client.SearchMovies("love", 0)
	if err != nil {
		t.Fatalf("SearchMovies() failed: %v", err)
	}

	// Should be limited to 20 results
	if len(results) > 20 {
		t.Errorf("SearchMovies returned %d results, expected max 20", len(results))
	}
}

func TestIntegrationSearchTVLimit(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Search for a common term that will have many results
	results, err := client.SearchTV("love")
	if err != nil {
		t.Fatalf("SearchTV() failed: %v", err)
	}

	// Should be limited to 20 results
	if len(results) > 20 {
		t.Errorf("SearchTV returned %d results, expected max 20", len(results))
	}
}

func TestIntegrationSaveTitleMovie(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "title-movie-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Fetch movie metadata to get the title
	metadata, err := client.FetchMovieMetadata("550")
	if err != nil {
		t.Fatalf("FetchMovieMetadata failed: %v", err)
	}

	// Save title
	err = client.saveTitle(metadata.Title, tmpDir)
	if err != nil {
		t.Fatalf("saveTitle() failed: %v", err)
	}

	// Verify file was created
	titlePath := filepath.Join(tmpDir, "title.txt")
	data, err := os.ReadFile(titlePath)
	if err != nil {
		t.Fatalf("Failed to read title file: %v", err)
	}

	if string(data) != metadata.Title {
		t.Errorf("Title file content = %q, want %q", string(data), metadata.Title)
	}

	if string(data) != "Fight Club" {
		t.Errorf("Expected title 'Fight Club', got %q", string(data))
	}
}

func TestIntegrationSaveTitleTV(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "title-tv-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Fetch TV metadata to get the name
	metadata, err := client.FetchTVMetadata("60059")
	if err != nil {
		t.Fatalf("FetchTVMetadata failed: %v", err)
	}

	// Save title
	err = client.saveTitle(metadata.Name, tmpDir)
	if err != nil {
		t.Fatalf("saveTitle() failed: %v", err)
	}

	// Verify file was created
	titlePath := filepath.Join(tmpDir, "title.txt")
	data, err := os.ReadFile(titlePath)
	if err != nil {
		t.Fatalf("Failed to read title file: %v", err)
	}

	if string(data) != metadata.Name {
		t.Errorf("Title file content = %q, want %q", string(data), metadata.Name)
	}

	if string(data) != "Better Call Saul" {
		t.Errorf("Expected title 'Better Call Saul', got %q", string(data))
	}
}

func TestIntegrationFetchAndSaveMetadataWithTitle(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "metadata-with-title-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a media item
	media := &Media{
		Title:     "Fight Club",
		Type:      Film,
		Year:      1999,
		TMDBID:    "550",
		Path:      tmpDir,
		DiskCount: 1,
	}

	// Fetch and save all metadata including title
	err = client.FetchAndSaveMetadata(media)
	if err != nil {
		t.Fatalf("FetchAndSaveMetadata() failed: %v", err)
	}

	// Verify poster was saved
	posterFile, found := media.FindPosterFile()
	if !found {
		t.Error("Poster file not created")
	} else {
		info, err := os.Stat(posterFile)
		if err != nil {
			t.Errorf("Poster file error: %v", err)
		} else if info.Size() == 0 {
			t.Error("Poster file is empty")
		}
	}

	// Verify description was saved
	descPath := filepath.Join(tmpDir, "description.txt")
	descData, err := os.ReadFile(descPath)
	if err != nil {
		t.Errorf("Description file not created: %v", err)
	} else if len(descData) == 0 {
		t.Error("Description file is empty")
	}

	// Verify genres were saved
	genrePath := filepath.Join(tmpDir, "genre.txt")
	genreData, err := os.ReadFile(genrePath)
	if err != nil {
		t.Errorf("Genre file not created: %v", err)
	} else if len(genreData) == 0 {
		t.Error("Genre file is empty")
	}

	// Verify title was saved
	titlePath := filepath.Join(tmpDir, "title.txt")
	titleData, err := os.ReadFile(titlePath)
	if err != nil {
		t.Errorf("Title file not created: %v", err)
	} else if len(titleData) == 0 {
		t.Error("Title file is empty")
	} else if string(titleData) != "Fight Club" {
		t.Errorf("Title content = %q, want 'Fight Club'", string(titleData))
	}
}

func TestIntegrationFetchAndSaveMetadataSkipsExistingTitle(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewTMDBClient(os.Getenv("TMDB_API_KEY"))

	// Create a temporary directory with existing metadata including title
	tmpDir, err := os.MkdirTemp("", "skip-existing-title-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create existing files
	posterPath := filepath.Join(tmpDir, "poster.jpg")
	err = os.WriteFile(posterPath, []byte("existing-poster-data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	descPath := filepath.Join(tmpDir, "description.txt")
	err = os.WriteFile(descPath, []byte("existing-description"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	genrePath := filepath.Join(tmpDir, "genre.txt")
	err = os.WriteFile(genrePath, []byte("existing-genres"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	titlePath := filepath.Join(tmpDir, "title.txt")
	err = os.WriteFile(titlePath, []byte("existing-title"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	media := &Media{
		Title:     "Test Film",
		Type:      Film,
		TMDBID:    "550",
		Path:      tmpDir,
		DiskCount: 1,
	}

	// This should not fetch anything since all files exist
	err = client.FetchAndSaveMetadata(media)
	if err != nil {
		t.Errorf("FetchAndSaveMetadata() failed with existing files: %v", err)
	}

	// Verify files weren't modified
	posterData, _ := os.ReadFile(posterPath)
	if string(posterData) != "existing-poster-data" {
		t.Error("Existing poster was modified")
	}

	descData, _ := os.ReadFile(descPath)
	if string(descData) != "existing-description" {
		t.Error("Existing description was modified")
	}

	genreData, _ := os.ReadFile(genrePath)
	if string(genreData) != "existing-genres" {
		t.Error("Existing genres were modified")
	}

	titleData, _ := os.ReadFile(titlePath)
	if string(titleData) != "existing-title" {
		t.Error("Existing title was modified")
	}
}
