package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestConvertToTrackList tests the conversion from MusicBrainz release to TrackList
func TestConvertToTrackList(t *testing.T) {
	mbRelease := &MBRelease{
		ID:    "test-release-id",
		Title: "Test Album",
		Media: []MBMedium{
			{
				Position: 1,
				Format:   "CD",
				Tracks: []MBTrack{
					{
						ID:       "track1",
						Position: 1,
						Title:    "Track One",
						Length:   180000, // 3 minutes
						Recording: MBRecording{
							ID:     "rec1",
							Title:  "Track One Recording",
							Length: 180000,
						},
					},
					{
						ID:       "track2",
						Position: 2,
						Title:    "Track Two",
						Length:   240000, // 4 minutes
						Recording: MBRecording{
							ID:     "rec2",
							Title:  "Track Two Recording",
							Length: 240000,
						},
					},
				},
			},
		},
	}

	trackList := ConvertToTrackList(mbRelease)

	if trackList == nil {
		t.Fatal("Expected trackList, got nil")
	}

	if len(trackList.Media) != 1 {
		t.Fatalf("Expected 1 medium, got %d", len(trackList.Media))
	}

	medium := trackList.Media[0]
	if medium.Position != 1 {
		t.Errorf("Expected medium position 1, got %d", medium.Position)
	}
	if medium.Format != "CD" {
		t.Errorf("Expected format 'CD', got '%s'", medium.Format)
	}
	if len(medium.Tracks) != 2 {
		t.Fatalf("Expected 2 tracks, got %d", len(medium.Tracks))
	}

	// Check first track
	track1 := medium.Tracks[0]
	if track1.Position != "1" {
		t.Errorf("Expected track position '1', got '%s'", track1.Position)
	}
	if track1.Title != "Track One" {
		t.Errorf("Expected track title 'Track One', got '%s'", track1.Title)
	}
	if track1.Duration != 180000 {
		t.Errorf("Expected track duration 180000, got %d", track1.Duration)
	}

	// Check second track
	track2 := medium.Tracks[1]
	if track2.Position != "2" {
		t.Errorf("Expected track position '2', got '%s'", track2.Position)
	}
	if track2.Title != "Track Two" {
		t.Errorf("Expected track title 'Track Two', got '%s'", track2.Title)
	}
	if track2.Duration != 240000 {
		t.Errorf("Expected track duration 240000, got %d", track2.Duration)
	}
}

// TestConvertToTrackListFallbackToRecording tests fallback to recording data
func TestConvertToTrackListFallbackToRecording(t *testing.T) {
	mbRelease := &MBRelease{
		ID:    "test-release-id",
		Title: "Test Album",
		Media: []MBMedium{
			{
				Position: 1,
				Format:   "CD",
				Tracks: []MBTrack{
					{
						ID:       "track1",
						Position: 1,
						Title:    "", // Empty title, should fall back to recording title
						Length:   0,  // Zero length, should fall back to recording length
						Recording: MBRecording{
							ID:     "rec1",
							Title:  "Recording Title",
							Length: 200000,
						},
					},
				},
			},
		},
	}

	trackList := ConvertToTrackList(mbRelease)

	if trackList == nil {
		t.Fatal("Expected trackList, got nil")
	}

	track := trackList.Media[0].Tracks[0]
	if track.Title != "Recording Title" {
		t.Errorf("Expected track title 'Recording Title', got '%s'", track.Title)
	}
	if track.Duration != 200000 {
		t.Errorf("Expected track duration 200000, got %d", track.Duration)
	}
}

// TestConvertToTrackListNil tests conversion with nil input
func TestConvertToTrackListNil(t *testing.T) {
	trackList := ConvertToTrackList(nil)
	if trackList != nil {
		t.Errorf("Expected nil trackList, got %v", trackList)
	}
}

// TestSaveTrackList tests saving track list to JSON file
func TestSaveTrackList(t *testing.T) {
	tmpDir := t.TempDir()

	trackList := &TrackList{
		Media: []Medium{
			{
				Position: 1,
				Format:   "CD",
				Tracks: []Track{
					{Position: "1", Title: "Test Track", Duration: 180000},
				},
			},
		},
	}

	err := SaveTrackList(trackList, tmpDir)
	if err != nil {
		t.Fatalf("Failed to save track list: %v", err)
	}

	// Verify file was created
	tracksPath := filepath.Join(tmpDir, "tracks.json")
	data, err := os.ReadFile(tracksPath)
	if err != nil {
		t.Fatalf("Failed to read tracks.json: %v", err)
	}

	// Verify JSON can be parsed back
	var loadedTrackList TrackList
	err = json.Unmarshal(data, &loadedTrackList)
	if err != nil {
		t.Fatalf("Failed to parse tracks.json: %v", err)
	}

	if len(loadedTrackList.Media) != 1 {
		t.Errorf("Expected 1 medium, got %d", len(loadedTrackList.Media))
	}
	if len(loadedTrackList.Media[0].Tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(loadedTrackList.Media[0].Tracks))
	}
	if loadedTrackList.Media[0].Tracks[0].Title != "Test Track" {
		t.Errorf("Expected track title 'Test Track', got '%s'", loadedTrackList.Media[0].Tracks[0].Title)
	}
}

// TestSaveTrackListNil tests saving nil track list
func TestSaveTrackListNil(t *testing.T) {
	tmpDir := t.TempDir()

	err := SaveTrackList(nil, tmpDir)
	if err == nil {
		t.Error("Expected error when saving nil track list, got nil")
	}
}

// TestFetchRelease tests fetching release from MusicBrainz API
func TestFetchRelease(t *testing.T) {
	// Create a test server that returns mock MusicBrainz data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify User-Agent header is set
		if userAgent := r.Header.Get("User-Agent"); userAgent == "" {
			t.Error("Expected User-Agent header to be set")
		}

		// Return mock release data
		release := MBRelease{
			ID:    "test-release-id",
			Title: "Test Album",
			Media: []MBMedium{
				{
					Position: 1,
					Format:   "CD",
					Tracks: []MBTrack{
						{
							ID:       "track1",
							Position: 1,
							Title:    "Test Track",
							Length:   180000,
							Recording: MBRecording{
								ID:     "rec1",
								Title:  "Test Track",
								Length: 180000,
							},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// For now, skip this test in favor of integration tests
	// The test structure is valid, but we'd need to refactor the client
	// to accept a custom base URL for unit testing
	t.Skip("Skipping unit test - use integration test instead")
	_ = server // Use server variable
}

// TestFetchAndSaveTrackListCaching tests that existing tracks.json files are not overwritten
func TestFetchAndSaveTrackListCaching(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test media with MusicBrainz ID
	media := &Media{
		Title:         "Test Album",
		Type:          Music,
		MusicBrainzID: "test-release-id",
		Path:          tmpDir,
	}

	// Create an existing tracks.json file
	existingTrackList := &TrackList{
		Media: []Medium{
			{
				Position: 1,
				Format:   "Existing Format",
				Tracks: []Track{
					{Position: "1", Title: "Existing Track", Duration: 999999},
				},
			},
		},
	}
	err := SaveTrackList(existingTrackList, tmpDir)
	if err != nil {
		t.Fatalf("Failed to create existing tracks.json: %v", err)
	}

	// Create client and try to fetch (should skip because file exists)
	client := NewMusicBrainzClient()
	err = client.FetchAndSaveTrackList(media)
	if err != nil {
		t.Fatalf("FetchAndSaveTrackList failed: %v", err)
	}

	// Load the tracks.json and verify it wasn't overwritten
	loadedTrackList := media.LoadTrackList()
	if loadedTrackList == nil {
		t.Fatal("Expected track list to be loaded")
	}

	if len(loadedTrackList.Media) != 1 {
		t.Errorf("Expected 1 medium, got %d", len(loadedTrackList.Media))
	}
	if loadedTrackList.Media[0].Format != "Existing Format" {
		t.Errorf("Expected existing format to be preserved, got '%s'", loadedTrackList.Media[0].Format)
	}
	if loadedTrackList.Media[0].Tracks[0].Title != "Existing Track" {
		t.Errorf("Expected existing track to be preserved, got '%s'", loadedTrackList.Media[0].Tracks[0].Title)
	}
}

// TestTrackFormatDuration tests the FormatDuration method
func TestTrackFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		expected string
	}{
		{"Zero duration", 0, ""},
		{"3 minutes", 180000, "3:00"},
		{"3 minutes 30 seconds", 210000, "3:30"},
		{"4 minutes 5 seconds", 245000, "4:05"},
		{"1 minute 1 second", 61000, "1:01"},
		{"10 minutes 59 seconds", 659000, "10:59"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track := Track{Duration: tt.duration}
			result := track.FormatDuration()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestMediaLoadTrackList tests loading track list from media directory
func TestMediaLoadTrackList(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a track list file
	trackList := &TrackList{
		Media: []Medium{
			{
				Position: 1,
				Format:   "CD",
				Tracks: []Track{
					{Position: "1", Title: "Test Track", Duration: 180000},
				},
			},
		},
	}
	err := SaveTrackList(trackList, tmpDir)
	if err != nil {
		t.Fatalf("Failed to save track list: %v", err)
	}

	// Create media and load track list
	media := &Media{
		Title: "Test Album",
		Type:  Music,
		Path:  tmpDir,
	}

	loadedTrackList := media.LoadTrackList()
	if loadedTrackList == nil {
		t.Fatal("Expected track list to be loaded")
	}

	if len(loadedTrackList.Media) != 1 {
		t.Errorf("Expected 1 medium, got %d", len(loadedTrackList.Media))
	}
	if loadedTrackList.Media[0].Tracks[0].Title != "Test Track" {
		t.Errorf("Expected track title 'Test Track', got '%s'", loadedTrackList.Media[0].Tracks[0].Title)
	}
}

// TestMediaLoadTrackListMissing tests loading when tracks.json doesn't exist
func TestMediaLoadTrackListMissing(t *testing.T) {
	tmpDir := t.TempDir()

	media := &Media{
		Title: "Test Album",
		Type:  Music,
		Path:  tmpDir,
	}

	loadedTrackList := media.LoadTrackList()
	if loadedTrackList != nil {
		t.Errorf("Expected nil track list, got %v", loadedTrackList)
	}
}
