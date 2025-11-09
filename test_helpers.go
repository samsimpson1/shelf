package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestData creates a temporary directory with the standard test media structure.
// This ensures tests have a consistent, isolated environment regardless of the state
// of the checked-in testdata directory.
//
// The structure created is:
//   - War of the Worlds (2025) [Film] - 1 disk, TMDB ID 755898
//   - Better Call Saul [TV] - 2 disks (Series 1), TMDB ID 60059
//   - No TMDB (2021) [Film] - 1 disk, no TMDB ID
func setupTestData(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Create War of the Worlds (2025) [Film]
	filmDir := filepath.Join(tmpDir, "War of the Worlds (2025) [Film]")
	if err := os.Mkdir(filmDir, 0755); err != nil {
		t.Fatalf("Failed to create film directory: %v", err)
	}
	filmDiskDir := filepath.Join(filmDir, "Disk [Blu-Ray]")
	if err := os.Mkdir(filmDiskDir, 0755); err != nil {
		t.Fatalf("Failed to create film disk directory: %v", err)
	}
	filmTMDB := filepath.Join(filmDir, "tmdb.txt")
	if err := os.WriteFile(filmTMDB, []byte("755898"), 0644); err != nil {
		t.Fatalf("Failed to create film tmdb.txt: %v", err)
	}

	// Create Better Call Saul [TV]
	tvDir := filepath.Join(tmpDir, "Better Call Saul [TV]")
	if err := os.Mkdir(tvDir, 0755); err != nil {
		t.Fatalf("Failed to create TV directory: %v", err)
	}
	tvDisk1 := filepath.Join(tvDir, "Series 1 Disk 1 [Blu-Ray]")
	if err := os.Mkdir(tvDisk1, 0755); err != nil {
		t.Fatalf("Failed to create TV disk 1 directory: %v", err)
	}
	tvDisk2 := filepath.Join(tvDir, "Series 1 Disk 2 [Blu-Ray]")
	if err := os.Mkdir(tvDisk2, 0755); err != nil {
		t.Fatalf("Failed to create TV disk 2 directory: %v", err)
	}
	tvTMDB := filepath.Join(tvDir, "tmdb.txt")
	if err := os.WriteFile(tvTMDB, []byte("60059"), 0644); err != nil {
		t.Fatalf("Failed to create TV tmdb.txt: %v", err)
	}

	// Create No TMDB (2021) [Film]
	noTMDBDir := filepath.Join(tmpDir, "No TMDB (2021) [Film]")
	if err := os.Mkdir(noTMDBDir, 0755); err != nil {
		t.Fatalf("Failed to create no TMDB directory: %v", err)
	}
	noTMDBDiskDir := filepath.Join(noTMDBDir, "Disk [DVD]")
	if err := os.Mkdir(noTMDBDiskDir, 0755); err != nil {
		t.Fatalf("Failed to create no TMDB disk directory: %v", err)
	}

	return tmpDir
}

// ensureTestDataExists ensures the testdata directory has the expected structure.
// This is useful for tests that need to use the checked-in testdata directory.
// It creates any missing directories but does not remove anything.
func ensureTestDataExists(t *testing.T) {
	t.Helper()

	testdataDir := "testdata"

	// Ensure War of the Worlds disk directory exists
	filmDiskDir := filepath.Join(testdataDir, "War of the Worlds (2025) [Film]", "Disk [Blu-Ray]")
	if err := os.MkdirAll(filmDiskDir, 0755); err != nil {
		t.Fatalf("Failed to create film disk directory: %v", err)
	}

	// Ensure Better Call Saul disk directories exist
	tvDisk1 := filepath.Join(testdataDir, "Better Call Saul [TV]", "Series 1 Disk 1 [Blu-Ray]")
	if err := os.MkdirAll(tvDisk1, 0755); err != nil {
		t.Fatalf("Failed to create TV disk 1 directory: %v", err)
	}
	tvDisk2 := filepath.Join(testdataDir, "Better Call Saul [TV]", "Series 1 Disk 2 [Blu-Ray]")
	if err := os.MkdirAll(tvDisk2, 0755); err != nil {
		t.Fatalf("Failed to create TV disk 2 directory: %v", err)
	}

	// No TMDB directory should already be complete (has .nodelete file)
}
