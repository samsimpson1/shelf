package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestImportScanner tests the ImportScanner functionality
func TestImportScanner(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create some test directories
	dir1 := filepath.Join(tmpDir, "test-import-1")
	dir2 := filepath.Join(tmpDir, "test-import-2")
	file1 := filepath.Join(tmpDir, "not-a-dir.txt")

	if err := os.Mkdir(dir1, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.Mkdir(dir2, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.WriteFile(file1, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create scanner
	scanner := NewImportScanner(tmpDir)

	// Scan directory
	imports, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	// Should find 2 directories (not the file)
	if len(imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(imports))
	}

	// Check that both directories are present
	foundDir1 := false
	foundDir2 := false
	for _, imp := range imports {
		if imp.Name == "test-import-1" {
			foundDir1 = true
		}
		if imp.Name == "test-import-2" {
			foundDir2 = true
		}
	}

	if !foundDir1 {
		t.Error("Expected to find test-import-1")
	}
	if !foundDir2 {
		t.Error("Expected to find test-import-2")
	}
}

// TestImportScannerNonexistentDir tests scanning a non-existent directory
func TestImportScannerNonexistentDir(t *testing.T) {
	scanner := NewImportScanner("/nonexistent/directory")
	_, err := scanner.Scan()
	if err == nil {
		t.Error("Expected error when scanning non-existent directory")
	}
}

// TestDetectDiskType tests disk type detection
func TestDetectDiskType(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		expectedType DiskType
		expectedConf bool
	}{
		{
			name: "Blu-ray BDMV",
			setupFunc: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, "BDMV"), 0755)
			},
			expectedType: DiskTypeBluRay,
			expectedConf: true,
		},
		{
			name: "DVD VIDEO_TS",
			setupFunc: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, "VIDEO_TS"), 0755)
			},
			expectedType: DiskTypeDVD,
			expectedConf: true,
		},
		{
			name: "Unknown structure",
			setupFunc: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, "other"), 0755)
			},
			expectedType: "",
			expectedConf: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := tt.setupFunc(tmpDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			diskType, confident := DetectDiskType(tmpDir)
			if diskType != tt.expectedType {
				t.Errorf("Expected type %q, got %q", tt.expectedType, diskType)
			}
			if confident != tt.expectedConf {
				t.Errorf("Expected confidence %v, got %v", tt.expectedConf, confident)
			}
		})
	}
}

// TestSanitizeName tests the name sanitization function
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Normal Title", "Normal Title"},
		{"Title: Subtitle", "Title_ Subtitle"},
		{"Title/Subtitle", "Title_Subtitle"},
		{"Title\\Subtitle", "Title_Subtitle"},
		{"Title<>Subtitle", "Title_Subtitle"},
		{"Title\"Subtitle", "Title'Subtitle"},
		{"Title|Subtitle", "Title_Subtitle"},
		{"Title?Subtitle", "Title_Subtitle"},
		{"Title*Subtitle", "Title_Subtitle"},
		{"Title   With   Spaces", "Title   With   Spaces"},
		{"Title___Multiple___Underscores", "Title_Multiple_Underscores"},
		{"  Leading and Trailing  ", "Leading and Trailing"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGenerateMediaDirName tests media directory name generation
func TestGenerateMediaDirName(t *testing.T) {
	tests := []struct {
		title    string
		year     int
		mediaType MediaType
		expected string
	}{
		{"The Matrix", 1999, Film, "The Matrix (1999) [Film]"},
		{"Breaking Bad", 0, TV, "Breaking Bad [TV]"},
		{"Title: Subtitle", 2020, Film, "Title_ Subtitle (2020) [Film]"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GenerateMediaDirName(tt.title, tt.year, tt.mediaType)
			if result != tt.expected {
				t.Errorf("GenerateMediaDirName(%q, %d, %v) = %q, want %q",
					tt.title, tt.year, tt.mediaType, result, tt.expected)
			}
		})
	}
}

// TestGenerateDiskDirName tests disk directory name generation
func TestGenerateDiskDirName(t *testing.T) {
	tests := []struct {
		diskType  string
		seriesNum int
		diskNum   int
		mediaType MediaType
		expected  string
	}{
		{"Blu-Ray", 0, 1, Film, "Disk [Blu-Ray]"},
		{"DVD", 1, 2, TV, "Series 1 Disk 2 [DVD]"},
		{"Blu-Ray UHD", 3, 1, TV, "Series 3 Disk 1 [Blu-Ray UHD]"},
		{"Custom: Type", 0, 1, Film, "Disk [Custom_ Type]"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GenerateDiskDirName(tt.diskType, tt.seriesNum, tt.diskNum, tt.mediaType)
			if result != tt.expected {
				t.Errorf("GenerateDiskDirName(%q, %d, %d, %v) = %q, want %q",
					tt.diskType, tt.seriesNum, tt.diskNum, tt.mediaType, result, tt.expected)
			}
		})
	}
}

// TestExecuteImportNewFilm tests importing a new film
func TestExecuteImportNewFilm(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("Failed to create media directory: %v", err)
	}

	// Create a test file in source
	testFile := filepath.Join(sourceDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create import session
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind: Film,
		Title:     "Test Film",
		Year:      2020,
		DiskType:  DiskTypeBluRay,
		AddToExisting: false,
	}

	// Execute import
	err := ExecuteImport(session, mediaDir)
	if err != nil {
		t.Fatalf("ExecuteImport() failed: %v", err)
	}

	// Verify destination exists
	expectedDest := filepath.Join(mediaDir, "Test Film (2020) [Film]", "Disk [Blu-Ray]")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Errorf("Expected destination directory not found: %s", expectedDest)
	}

	// Verify test file was moved
	movedFile := filepath.Join(expectedDest, "test.txt")
	if _, err := os.Stat(movedFile); os.IsNotExist(err) {
		t.Error("Test file was not moved to destination")
	}

	// Verify source no longer exists
	if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
		t.Error("Source directory still exists after import")
	}
}

// TestExecuteImportNewTV tests importing a new TV show
func TestExecuteImportNewTV(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("Failed to create media directory: %v", err)
	}

	// Create import session
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind: TV,
		Title:     "Test Show",
		SeriesNum: 1,
		DiskNum:   2,
		DiskType:  DiskTypeDVD,
		AddToExisting: false,
	}

	// Execute import
	err := ExecuteImport(session, mediaDir)
	if err != nil {
		t.Fatalf("ExecuteImport() failed: %v", err)
	}

	// Verify destination exists
	expectedDest := filepath.Join(mediaDir, "Test Show [TV]", "Series 1 Disk 2 [DVD]")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Errorf("Expected destination directory not found: %s", expectedDest)
	}
}

// TestExecuteImportAddToExisting tests adding a disk to existing media
func TestExecuteImportAddToExisting(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")
	existingMedia := filepath.Join(mediaDir, "Existing Film (2020) [Film]")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(existingMedia, 0755); err != nil {
		t.Fatalf("Failed to create existing media directory: %v", err)
	}

	// Create import session
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind: Film,
		Title:     "Existing Film",
		Year:      2020,
		DiskType:  DiskTypeBluRay,
		AddToExisting: true,
		ExistingMediaPath: existingMedia,
	}

	// Execute import
	err := ExecuteImport(session, mediaDir)
	if err != nil {
		t.Fatalf("ExecuteImport() failed: %v", err)
	}

	// Verify destination exists within existing media
	expectedDest := filepath.Join(existingMedia, "Disk [Blu-Ray]")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Errorf("Expected destination directory not found: %s", expectedDest)
	}
}

// TestExecuteImportWithTMDB tests importing with TMDB metadata
func TestExecuteImportWithTMDB(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("Failed to create media directory: %v", err)
	}

	// Create import session with TMDB data
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind: Film,
		Title:     "User Title",
		Year:      2020,
		TMDBID:    "12345",
		TMDBTitle: "Official TMDB Title",
		TMDBYear:  2021,
		DiskType:  DiskTypeBluRay,
		AddToExisting: false,
	}

	// Execute import
	err := ExecuteImport(session, mediaDir)
	if err != nil {
		t.Fatalf("ExecuteImport() failed: %v", err)
	}

	// Verify destination uses TMDB title and year
	expectedDest := filepath.Join(mediaDir, "Official TMDB Title (2021) [Film]", "Disk [Blu-Ray]")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Errorf("Expected destination directory not found: %s", expectedDest)
	}

	// Verify tmdb.txt was created
	tmdbFile := filepath.Join(mediaDir, "Official TMDB Title (2021) [Film]", "tmdb.txt")
	if _, err := os.Stat(tmdbFile); os.IsNotExist(err) {
		t.Error("tmdb.txt file was not created")
	} else {
		// Verify content
		content, err := os.ReadFile(tmdbFile)
		if err != nil {
			t.Errorf("Failed to read tmdb.txt: %v", err)
		}
		if strings.TrimSpace(string(content)) != "12345" {
			t.Errorf("tmdb.txt content = %q, want %q", string(content), "12345")
		}
	}
}

// TestExecuteImportDuplicateDest tests error when destination already exists
func TestExecuteImportDuplicateDest(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	importDir := filepath.Join(tmpDir, "import")
	mediaDir := filepath.Join(tmpDir, "media")
	sourceDir := filepath.Join(importDir, "source-disk")
	existingDest := filepath.Join(mediaDir, "Test Film (2020) [Film]", "Disk [Blu-Ray]")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(existingDest, 0755); err != nil {
		t.Fatalf("Failed to create existing destination: %v", err)
	}

	// Create import session
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind: Film,
		Title:     "Test Film",
		Year:      2020,
		DiskType:  DiskTypeBluRay,
		AddToExisting: false,
	}

	// Execute import - should fail
	err := ExecuteImport(session, mediaDir)
	if err == nil {
		t.Error("Expected error when destination already exists, got nil")
	}
}
