package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExecuteImportWithMusicBrainz tests importing with MusicBrainz metadata
func TestExecuteImportWithMusicBrainz(t *testing.T) {
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

	// Create import session with MusicBrainz data
	session := &ImportSession{
		SourceDir: &ImportDirectory{
			Name: "source-disk",
			Path: sourceDir,
		},
		MediaKind:        Music,
		Title:            "User Title",
		MusicBrainzID:    "mb-release-id-123",
		MusicBrainzTitle: "Official MB Title",
		DiskType:         DiskTypeCustom,
		DiskTypeCustom:   "CD",
		AddToExisting:    false,
	}

	// Execute import
	err := ExecuteImport(session, mediaDir)
	if err != nil {
		t.Fatalf("ExecuteImport() failed: %v", err)
	}

	// Verify destination uses User Title (since we didn't implement MB title override logic yet, or did we?
	// In ExecuteImport:
	// finalTitle := session.Title
	// if session.TMDBTitle != "" { finalTitle = session.TMDBTitle }
	// We didn't add MusicBrainzTitle check there yet! Let's check import.go again.
	// Wait, I missed that in the implementation plan execution. I should fix that.

	// Verify destination uses MusicBrainz title
	expectedDest := filepath.Join(mediaDir, "Official MB Title [Music]", "Disk [CD]")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Errorf("Expected destination directory not found: %s", expectedDest)
	}

	// Verify musicbrainz.txt was created
	mbFile := filepath.Join(mediaDir, "Official MB Title [Music]", "musicbrainz.txt")
	if _, err := os.Stat(mbFile); os.IsNotExist(err) {
		t.Error("musicbrainz.txt file was not created")
	} else {
		// Verify content
		content, err := os.ReadFile(mbFile)
		if err != nil {
			t.Errorf("Failed to read musicbrainz.txt: %v", err)
		}
		if strings.TrimSpace(string(content)) != "mb-release-id-123" {
			t.Errorf("musicbrainz.txt content = %q, want %q", string(content), "mb-release-id-123")
		}
	}
}
