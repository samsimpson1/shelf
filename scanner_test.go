package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner("/test/path")
	if scanner.mediaDir != "/test/path" {
		t.Errorf("NewScanner() mediaDir = %v, want /test/path", scanner.mediaDir)
	}
	if scanner.tmdbClient != nil {
		t.Error("NewScanner() should not have TMDB client")
	}
}

func TestNewScannerWithTMDB(t *testing.T) {
	tmdbClient := NewTMDBClient("test-key")
	scanner := NewScannerWithTMDB("/test/path", tmdbClient)

	if scanner.mediaDir != "/test/path" {
		t.Errorf("NewScannerWithTMDB() mediaDir = %v, want /test/path", scanner.mediaDir)
	}
	if scanner.tmdbClient == nil {
		t.Error("NewScannerWithTMDB() should have TMDB client")
	}
}

func TestScanNonexistentDirectory(t *testing.T) {
	scanner := NewScanner("/nonexistent/directory")
	_, err := scanner.Scan()
	if err == nil {
		t.Error("Scan() expected error for nonexistent directory, got nil")
	}
}

func TestScanFileNotDirectory(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	scanner := NewScanner(tmpFile.Name())
	_, err = scanner.Scan()
	if err == nil {
		t.Error("Scan() expected error for file path, got nil")
	}
}

func TestScanTestdata(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// We expect 3 media items from our testdata
	expectedCount := 3
	if len(mediaList) != expectedCount {
		t.Errorf("Scan() returned %d items, want %d", len(mediaList), expectedCount)
	}

	// Create a map for easy lookup
	mediaMap := make(map[string]Media)
	for _, media := range mediaList {
		mediaMap[media.Title] = media
	}

	// Test War of the Worlds film
	if media, ok := mediaMap["War of the Worlds"]; ok {
		if media.Type != Film {
			t.Errorf("War of the Worlds: Type = %v, want Film", media.Type)
		}
		if media.Year != 2025 {
			t.Errorf("War of the Worlds: Year = %v, want 2025", media.Year)
		}
		if media.DiskCount != 1 {
			t.Errorf("War of the Worlds: DiskCount = %v, want 1", media.DiskCount)
		}
		if len(media.Disks) != 1 {
			t.Errorf("War of the Worlds: len(Disks) = %v, want 1", len(media.Disks))
		}
		if media.TMDBID != "755898" {
			t.Errorf("War of the Worlds: TMDBID = %v, want 755898", media.TMDBID)
		}
	} else {
		t.Error("War of the Worlds not found in scan results")
	}

	// Test Better Call Saul TV show
	if media, ok := mediaMap["Better Call Saul"]; ok {
		if media.Type != TV {
			t.Errorf("Better Call Saul: Type = %v, want TV", media.Type)
		}
		if media.Year != 0 {
			t.Errorf("Better Call Saul: Year = %v, want 0", media.Year)
		}
		if media.DiskCount != 2 {
			t.Errorf("Better Call Saul: DiskCount = %v, want 2", media.DiskCount)
		}
		if len(media.Disks) != 2 {
			t.Errorf("Better Call Saul: len(Disks) = %v, want 2", len(media.Disks))
		}
		if media.TMDBID != "60059" {
			t.Errorf("Better Call Saul: TMDBID = %v, want 60059", media.TMDBID)
		}
	} else {
		t.Error("Better Call Saul not found in scan results")
	}

	// Test No TMDB film (no TMDB ID)
	if media, ok := mediaMap["No TMDB"]; ok {
		if media.Type != Film {
			t.Errorf("No TMDB: Type = %v, want Film", media.Type)
		}
		if media.Year != 2021 {
			t.Errorf("No TMDB: Year = %v, want 2021", media.Year)
		}
		if media.DiskCount != 1 {
			t.Errorf("No TMDB: DiskCount = %v, want 1", media.DiskCount)
		}
		if media.TMDBID != "" {
			t.Errorf("No TMDB: TMDBID = %v, want empty string", media.TMDBID)
		}
	} else {
		t.Error("No TMDB not found in scan results")
	}
}

func TestParseFilmName(t *testing.T) {
	tests := []struct {
		name      string
		dirName   string
		shouldMatch bool
		wantTitle string
		wantYear  int
	}{
		{
			name:      "Valid film",
			dirName:   "War of the Worlds (2025) [Film]",
			shouldMatch: true,
			wantTitle: "War of the Worlds",
			wantYear:  2025,
		},
		{
			name:      "Film with special chars",
			dirName:   "The Lord of the Rings: The Fellowship of the Ring (2001) [Film]",
			shouldMatch: true,
			wantTitle: "The Lord of the Rings: The Fellowship of the Ring",
			wantYear:  2001,
		},
		{
			name:      "TV show (should not match)",
			dirName:   "Breaking Bad [TV]",
			shouldMatch: false,
		},
		{
			name:      "Invalid format - no year",
			dirName:   "Some Movie [Film]",
			shouldMatch: false,
		},
		{
			name:      "Invalid format - no brackets",
			dirName:   "Some Movie (2020)",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("/tmp")
			media, ok := scanner.parseFilm(tt.dirName, "/fake/path")

			if ok != tt.shouldMatch {
				t.Errorf("parseFilm() matched = %v, want %v", ok, tt.shouldMatch)
			}

			if tt.shouldMatch {
				if media.Title != tt.wantTitle {
					t.Errorf("parseFilm() Title = %v, want %v", media.Title, tt.wantTitle)
				}
				if media.Year != tt.wantYear {
					t.Errorf("parseFilm() Year = %v, want %v", media.Year, tt.wantYear)
				}
				if media.Type != Film {
					t.Errorf("parseFilm() Type = %v, want Film", media.Type)
				}
			}
		})
	}
}

func TestParseTVName(t *testing.T) {
	tests := []struct {
		name      string
		dirName   string
		shouldMatch bool
		wantTitle string
	}{
		{
			name:      "Valid TV show",
			dirName:   "Better Call Saul [TV]",
			shouldMatch: true,
			wantTitle: "Better Call Saul",
		},
		{
			name:      "TV show with special chars",
			dirName:   "Game of Thrones: Season One [TV]",
			shouldMatch: true,
			wantTitle: "Game of Thrones: Season One",
		},
		{
			name:      "Film (should not match)",
			dirName:   "Some Movie (2020) [Film]",
			shouldMatch: false,
		},
		{
			name:      "Invalid format - no brackets",
			dirName:   "Some Show TV",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner("/tmp")
			media, ok := scanner.parseTV(tt.dirName, "/fake/path")

			if ok != tt.shouldMatch {
				t.Errorf("parseTV() matched = %v, want %v", ok, tt.shouldMatch)
			}

			if tt.shouldMatch {
				if media.Title != tt.wantTitle {
					t.Errorf("parseTV() Title = %v, want %v", media.Title, tt.wantTitle)
				}
				if media.Type != TV {
					t.Errorf("parseTV() Type = %v, want TV", media.Type)
				}
				if media.Year != 0 {
					t.Errorf("parseTV() Year = %v, want 0", media.Year)
				}
			}
		})
	}
}

func TestCountFilmDisks(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// War of the Worlds should have 1 disk
	path := filepath.Join(testDir, "War of the Worlds (2025) [Film]")
	count := scanner.countFilmDisks(path)
	if count != 1 {
		t.Errorf("countFilmDisks() = %v, want 1", count)
	}

	// Nonexistent path should return 0
	count = scanner.countFilmDisks("/nonexistent/path")
	if count != 0 {
		t.Errorf("countFilmDisks() for nonexistent path = %v, want 0", count)
	}
}

func TestCountTVDisks(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// Better Call Saul should have 2 disks
	path := filepath.Join(testDir, "Better Call Saul [TV]")
	count := scanner.countTVDisks(path)
	if count != 2 {
		t.Errorf("countTVDisks() = %v, want 2", count)
	}

	// Nonexistent path should return 0
	count = scanner.countTVDisks("/nonexistent/path")
	if count != 0 {
		t.Errorf("countTVDisks() for nonexistent path = %v, want 0", count)
	}
}

func TestReadTMDBID(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// War of the Worlds should have TMDB ID
	path := filepath.Join(testDir, "War of the Worlds (2025) [Film]")
	id := scanner.readTMDBID(path)
	if id != "755898" {
		t.Errorf("readTMDBID() = %v, want 755898", id)
	}

	// No TMDB film should not have TMDB ID
	path = filepath.Join(testDir, "No TMDB (2021) [Film]")
	id = scanner.readTMDBID(path)
	if id != "" {
		t.Errorf("readTMDBID() = %v, want empty string", id)
	}

	// Nonexistent path should return empty string
	id = scanner.readTMDBID("/nonexistent/path")
	if id != "" {
		t.Errorf("readTMDBID() for nonexistent path = %v, want empty string", id)
	}
}

func TestReadTMDBIDWithWhitespace(t *testing.T) {
	// Create a temporary directory with a tmdb.txt file containing whitespace
	tmpDir, err := os.MkdirTemp("", "tmdb-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmdbFile := filepath.Join(tmpDir, "tmdb.txt")
	err = os.WriteFile(tmdbFile, []byte("  12345  \n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner("/tmp")
	id := scanner.readTMDBID(tmpDir)
	if id != "12345" {
		t.Errorf("readTMDBID() = %v, want 12345 (whitespace should be trimmed)", id)
	}
}

func TestExtractFormat(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected string
	}{
		{
			name:     "Blu-Ray format",
			dirName:  "Disk [Blu-Ray]",
			expected: "Blu-Ray",
		},
		{
			name:     "Blu-Ray UHD format",
			dirName:  "Series 1 Disk 1 [Blu-Ray UHD]",
			expected: "Blu-Ray UHD",
		},
		{
			name:     "DVD format",
			dirName:  "Disk [DVD]",
			expected: "DVD",
		},
		{
			name:     "No brackets",
			dirName:  "Disk",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFormat(tt.dirName)
			if result != tt.expected {
				t.Errorf("extractFormat(%q) = %q, want %q", tt.dirName, result, tt.expected)
			}
		})
	}
}

func TestCollectFilmDisks(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// War of the Worlds should have 1 disk with details
	path := filepath.Join(testDir, "War of the Worlds (2025) [Film]")
	disks := scanner.collectFilmDisks(path)

	if len(disks) != 1 {
		t.Fatalf("collectFilmDisks() returned %d disks, want 1", len(disks))
	}

	disk := disks[0]
	if disk.Name != "Disk 1" {
		t.Errorf("Disk Name = %q, want %q", disk.Name, "Disk 1")
	}
	if disk.Format != "Blu-Ray" {
		t.Errorf("Disk Format = %q, want %q", disk.Format, "Blu-Ray")
	}
	if disk.SizeGB < 0 {
		t.Errorf("Disk SizeGB = %v, should be >= 0", disk.SizeGB)
	}

	// Nonexistent path should return empty slice
	disks = scanner.collectFilmDisks("/nonexistent/path")
	if len(disks) != 0 {
		t.Errorf("collectFilmDisks() for nonexistent path returned %d disks, want 0", len(disks))
	}
}

func TestCollectTVDisks(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// Better Call Saul should have 2 disks with details
	path := filepath.Join(testDir, "Better Call Saul [TV]")
	disks := scanner.collectTVDisks(path)

	if len(disks) != 2 {
		t.Fatalf("collectTVDisks() returned %d disks, want 2", len(disks))
	}

	// First disk
	if disks[0].Name != "Series 1 Disk 1" {
		t.Errorf("Disk[0] Name = %q, want %q", disks[0].Name, "Series 1 Disk 1")
	}
	if disks[0].Format != "Blu-Ray" {
		t.Errorf("Disk[0] Format = %q, want %q", disks[0].Format, "Blu-Ray")
	}
	if disks[0].SizeGB < 0 {
		t.Errorf("Disk[0] SizeGB = %v, should be >= 0", disks[0].SizeGB)
	}

	// Second disk
	if disks[1].Name != "Series 1 Disk 2" {
		t.Errorf("Disk[1] Name = %q, want %q", disks[1].Name, "Series 1 Disk 2")
	}
	if disks[1].Format != "Blu-Ray UHD" {
		t.Errorf("Disk[1] Format = %q, want %q", disks[1].Format, "Blu-Ray UHD")
	}
	if disks[1].SizeGB < 0 {
		t.Errorf("Disk[1] SizeGB = %v, should be >= 0", disks[1].SizeGB)
	}

	// Nonexistent path should return empty slice
	disks = scanner.collectTVDisks("/nonexistent/path")
	if len(disks) != 0 {
		t.Errorf("collectTVDisks() for nonexistent path returned %d disks, want 0", len(disks))
	}
}

func TestScanIgnoresNonMediaDirectories(t *testing.T) {
	// Create a temporary directory with some non-media directories
	tmpDir, err := os.MkdirTemp("", "scan-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some directories that don't match patterns
	os.Mkdir(filepath.Join(tmpDir, "Random Folder"), 0755)
	os.Mkdir(filepath.Join(tmpDir, "Another Directory"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644)

	// Create one valid film
	filmDir := filepath.Join(tmpDir, "Test Film (2020) [Film]")
	os.Mkdir(filmDir, 0755)
	os.Mkdir(filepath.Join(filmDir, "Disk [DVD]"), 0755)

	scanner := NewScanner(tmpDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Should only find 1 media item
	if len(mediaList) != 1 {
		t.Errorf("Scan() returned %d items, want 1", len(mediaList))
	}

	if len(mediaList) > 0 && mediaList[0].Title != "Test Film" {
		t.Errorf("Scan() found wrong media: %v", mediaList[0].Title)
	}
}

// Test cache miss scenario - sizes calculated when cache missing
func TestSizeCacheMiss(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// Scan for War of the Worlds (first scan, no cache)
	filmPath := filepath.Join(testDir, "War of the Worlds (2025) [Film]")
	disks := scanner.collectFilmDisks(filmPath)

	if len(disks) != 1 {
		t.Fatalf("collectFilmDisks() returned %d disks, want 1", len(disks))
	}

	// Check that cache file was created
	cachePath := filepath.Join(filmPath, "sizes.json")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Cache file was not created after first scan")
	}

	// Verify cache content
	cache := loadSizeCache(filmPath)
	if len(cache) != 1 {
		t.Errorf("Cache contains %d entries, want 1", len(cache))
	}

	if _, exists := cache["Disk [Blu-Ray]"]; !exists {
		t.Error("Cache does not contain expected disk entry")
	}
}

// Test cache hit scenario - sizes loaded from existing valid cache
func TestSizeCacheHit(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// First scan to create cache
	filmPath := filepath.Join(testDir, "War of the Worlds (2025) [Film]")
	disks1 := scanner.collectFilmDisks(filmPath)

	if len(disks1) != 1 {
		t.Fatalf("First scan returned %d disks, want 1", len(disks1))
	}

	originalSize := disks1[0].SizeGB

	// Second scan should use cache
	disks2 := scanner.collectFilmDisks(filmPath)

	if len(disks2) != 1 {
		t.Fatalf("Second scan returned %d disks, want 1", len(disks2))
	}

	// Size should be the same (loaded from cache)
	if disks2[0].SizeGB != originalSize {
		t.Errorf("Cached size %v differs from original %v", disks2[0].SizeGB, originalSize)
	}
}

// Test cache update scenario - cache file created/updated after calculation
func TestSizeCacheUpdate(t *testing.T) {
	testDir := setupTestData(t)

	// Create a film directory with multiple disks
	filmPath := filepath.Join(testDir, "Multi Disk Film (2020) [Film]")
	os.Mkdir(filmPath, 0755)
	disk1Path := filepath.Join(filmPath, "Disk [Blu-Ray]")
	disk2Path := filepath.Join(filmPath, "Disk [DVD]")
	os.Mkdir(disk1Path, 0755)
	os.Mkdir(disk2Path, 0755)

	// Add some test files to each disk
	os.WriteFile(filepath.Join(disk1Path, "file1.txt"), []byte("test data 1"), 0644)
	os.WriteFile(filepath.Join(disk2Path, "file2.txt"), []byte("test data 2"), 0644)

	scanner := NewScanner(testDir)

	// First scan - should create cache with both disks
	disks := scanner.collectFilmDisks(filmPath)
	if len(disks) != 2 {
		t.Fatalf("collectFilmDisks() returned %d disks, want 2", len(disks))
	}

	// Verify cache was created and contains both disks
	cache := loadSizeCache(filmPath)
	if len(cache) != 2 {
		t.Errorf("Cache contains %d entries, want 2", len(cache))
	}

	if _, exists := cache["Disk [Blu-Ray]"]; !exists {
		t.Error("Cache missing Disk [Blu-Ray]")
	}
	if _, exists := cache["Disk [DVD]"]; !exists {
		t.Error("Cache missing Disk [DVD]")
	}
}

// Test invalid cache handling - graceful fallback to calculation
func TestSizeCacheInvalid(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	filmPath := filepath.Join(testDir, "War of the Worlds (2025) [Film]")

	// Create an invalid cache file
	cachePath := filepath.Join(filmPath, "sizes.json")
	os.WriteFile(cachePath, []byte("invalid json{{{"), 0644)

	// Scan should handle invalid cache gracefully
	disks := scanner.collectFilmDisks(filmPath)

	if len(disks) != 1 {
		t.Fatalf("collectFilmDisks() with invalid cache returned %d disks, want 1", len(disks))
	}

	// Cache should be regenerated with valid data
	cache := loadSizeCache(filmPath)
	if len(cache) == 0 {
		t.Error("Cache was not regenerated after encountering invalid cache")
	}
}

// Test TV disk caching
func TestTVDiskCaching(t *testing.T) {
	testDir := setupTestData(t)
	scanner := NewScanner(testDir)

	// Scan Better Call Saul TV show
	tvPath := filepath.Join(testDir, "Better Call Saul [TV]")
	disks := scanner.collectTVDisks(tvPath)

	if len(disks) != 2 {
		t.Fatalf("collectTVDisks() returned %d disks, want 2", len(disks))
	}

	// Check that cache file was created
	cachePath := filepath.Join(tvPath, "sizes.json")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Cache file was not created for TV show")
	}

	// Verify cache content
	cache := loadSizeCache(tvPath)
	if len(cache) != 2 {
		t.Errorf("TV cache contains %d entries, want 2", len(cache))
	}

	if _, exists := cache["Series 1 Disk 1 [Blu-Ray]"]; !exists {
		t.Error("Cache missing Series 1 Disk 1")
	}
	if _, exists := cache["Series 1 Disk 2 [Blu-Ray UHD]"]; !exists {
		t.Error("Cache missing Series 1 Disk 2")
	}

	// Second scan should use cache
	disks2 := scanner.collectTVDisks(tvPath)
	if len(disks2) != 2 {
		t.Fatalf("Second TV scan returned %d disks, want 2", len(disks2))
	}

	// Sizes should match (loaded from cache)
	for i := range disks {
		if disks2[i].SizeGB != disks[i].SizeGB {
			t.Errorf("Disk %d cached size %v differs from original %v", i, disks2[i].SizeGB, disks[i].SizeGB)
		}
	}
}

// Test loadSizeCache with missing file
func TestLoadSizeCacheMissingFile(t *testing.T) {
	tmpDir := t.TempDir()

	cache := loadSizeCache(tmpDir)

	if cache == nil {
		t.Error("loadSizeCache() returned nil, want empty map")
	}

	if len(cache) != 0 {
		t.Errorf("loadSizeCache() returned %d entries, want 0", len(cache))
	}
}

// Test saveSizeCache
func TestSaveSizeCache(t *testing.T) {
	tmpDir := t.TempDir()

	testCache := map[string]int64{
		"Disk [Blu-Ray]": 23456789012,
		"Disk [DVD]":     4567890123,
	}

	err := saveSizeCache(tmpDir, testCache)
	if err != nil {
		t.Fatalf("saveSizeCache() error = %v", err)
	}

	// Verify file was created
	cachePath := filepath.Join(tmpDir, "sizes.json")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created")
	}

	// Load and verify content
	loadedCache := loadSizeCache(tmpDir)

	if len(loadedCache) != 2 {
		t.Errorf("Loaded cache contains %d entries, want 2", len(loadedCache))
	}

	if loadedCache["Disk [Blu-Ray]"] != 23456789012 {
		t.Errorf("Loaded cache has wrong size for Blu-Ray disk")
	}

	if loadedCache["Disk [DVD]"] != 4567890123 {
		t.Errorf("Loaded cache has wrong size for DVD disk")
	}
}
