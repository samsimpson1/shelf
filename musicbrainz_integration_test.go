package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// checkMusicBrainzAccess tests if we can reach the MusicBrainz API
func checkMusicBrainzAccess(t *testing.T) {
	t.Helper()
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "https://musicbrainz.org/ws/2/", nil)
	req.Header.Set("User-Agent", musicBrainzUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Cannot reach MusicBrainz API (no internet access): %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		t.Skip("MusicBrainz API returned Forbidden - may need different User-Agent or rate limiting")
	}
}

// TestIntegrationFetchRelease tests fetching a real release from MusicBrainz
// Note: This test makes real API calls to MusicBrainz and requires internet access
// Run with: go test -v -run TestIntegrationFetchRelease
func TestIntegrationFetchRelease(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	client := NewMusicBrainzClient()

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	// Test with a well-known release: Pink Floyd - Dark Side of the Moon
	// MusicBrainz ID: c712a2bc-1e57-46c1-9432-8b9a33dc4159 (1994 UK CD)
	mbid := "c712a2bc-1e57-46c1-9432-8b9a33dc4159"

	release, err := client.FetchRelease(mbid)
	if err != nil {
		t.Fatalf("Failed to fetch release: %v", err)
	}

	// Verify release data
	if release.ID != mbid {
		t.Errorf("Expected release ID %s, got %s", mbid, release.ID)
	}

	if release.Title != "Dark Side of the Moon" {
		t.Errorf("Expected title 'Dark Side of the Moon', got '%s'", release.Title)
	}

	// Should have media (discs)
	if len(release.Media) == 0 {
		t.Fatal("Expected at least one medium, got none")
	}

	// First medium should have tracks
	if len(release.Media[0].Tracks) == 0 {
		t.Fatal("Expected tracks on first medium, got none")
	}

	// Verify some track data exists
	firstTrack := release.Media[0].Tracks[0]
	if firstTrack.Title == "" && firstTrack.Recording.Title == "" {
		t.Error("Expected track to have a title")
	}

	t.Logf("Successfully fetched release: %s", release.Title)
	t.Logf("  Media count: %d", len(release.Media))
	t.Logf("  Track count on first medium: %d", len(release.Media[0].Tracks))
}

// TestIntegrationConvertAndSaveTrackList tests the complete workflow
// Note: This test makes real API calls to MusicBrainz and requires internet access
func TestIntegrationConvertAndSaveTrackList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	tmpDir := t.TempDir()
	client := NewMusicBrainzClient()

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	// Test with Pink Floyd - Dark Side of the Moon
	mbid := "c712a2bc-1e57-46c1-9432-8b9a33dc4159"

	// Fetch release
	release, err := client.FetchRelease(mbid)
	if err != nil {
		t.Fatalf("Failed to fetch release: %v", err)
	}

	// Convert to track list
	trackList := ConvertToTrackList(release)
	if trackList == nil {
		t.Fatal("Failed to convert to track list")
	}

	// Save track list
	err = SaveTrackList(trackList, tmpDir)
	if err != nil {
		t.Fatalf("Failed to save track list: %v", err)
	}

	// Verify file was created
	tracksPath := filepath.Join(tmpDir, "tracks.json")
	if _, err := os.Stat(tracksPath); os.IsNotExist(err) {
		t.Fatal("tracks.json was not created")
	}

	// Load and verify
	media := &Media{
		Title: "Dark Side of the Moon",
		Type:  Music,
		Path:  tmpDir,
	}

	loadedTrackList, _ := NewMetadataLoader().LoadTrackList(media.Path)
	if loadedTrackList == nil {
		t.Fatal("Failed to load track list")
	}

	if len(loadedTrackList.Media) != len(trackList.Media) {
		t.Errorf("Expected %d media, got %d", len(trackList.Media), len(loadedTrackList.Media))
	}

	t.Logf("Successfully converted and saved track list")
	t.Logf("  Media count: %d", len(loadedTrackList.Media))
	if len(loadedTrackList.Media) > 0 {
		t.Logf("  First medium track count: %d", len(loadedTrackList.Media[0].Tracks))
	}
}

// TestIntegrationFetchAndSaveTrackList tests the complete Media workflow
// Note: This test makes real API calls to MusicBrainz and requires internet access
func TestIntegrationFetchAndSaveTrackList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	tmpDir := t.TempDir()
	client := NewMusicBrainzClient()

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	// Create media with MusicBrainz ID
	media := &Media{
		Title:         "Dark Side of the Moon",
		Type:          Music,
		MusicBrainzID: "c712a2bc-1e57-46c1-9432-8b9a33dc4159",
		Path:          tmpDir,
	}

	// Fetch and save
	err := client.FetchAndSaveTrackList(media)
	if err != nil {
		t.Fatalf("FetchAndSaveTrackList failed: %v", err)
	}

	// Verify file was created
	tracksPath := filepath.Join(tmpDir, "tracks.json")
	if _, err := os.Stat(tracksPath); os.IsNotExist(err) {
		t.Fatal("tracks.json was not created")
	}

	// Load and verify
	trackList, _ := NewMetadataLoader().LoadTrackList(media.Path)
	if trackList == nil {
		t.Fatal("Failed to load track list")
	}

	if len(trackList.Media) == 0 {
		t.Fatal("Expected at least one medium")
	}

	if len(trackList.Media[0].Tracks) == 0 {
		t.Fatal("Expected tracks on first medium")
	}

	// Verify track has proper data
	firstTrack := trackList.Media[0].Tracks[0]
	if firstTrack.Title == "" {
		t.Error("Expected first track to have a title")
	}
	if firstTrack.Position == "" {
		t.Error("Expected first track to have a position")
	}

	t.Logf("Successfully fetched and saved track list for %s", media.Title)
	t.Logf("  Media count: %d", len(trackList.Media))
	t.Logf("  First track: %s - %s (%s)", firstTrack.Position, firstTrack.Title, firstTrack.FormatDuration())
}

// TestIntegrationScannerWithMusicBrainz tests scanner integration
// Note: This test makes real API calls to MusicBrainz and requires internet access
func TestIntegrationScannerWithMusicBrainz(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	tmpDir := t.TempDir()

	// Create a music directory structure
	musicDir := filepath.Join(tmpDir, "Pink Floyd - The Dark Side of the Moon [Music]")
	err := os.Mkdir(musicDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create music directory: %v", err)
	}

	// Create a disk directory
	diskDir := filepath.Join(musicDir, "Disk [CD]")
	err = os.Mkdir(diskDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create disk directory: %v", err)
	}

	// Write MusicBrainz ID file
	mbidPath := filepath.Join(musicDir, "musicbrainz.txt")
	err = os.WriteFile(mbidPath, []byte("c712a2bc-1e57-46c1-9432-8b9a33dc4159"), 0644)
	if err != nil {
		t.Fatalf("Failed to write musicbrainz.txt: %v", err)
	}

	// Create scanner with MusicBrainz client
	musicBrainzClient := NewMusicBrainzClient()
	scanner := NewScannerWithClients(tmpDir, nil, musicBrainzClient)

	// Scan directory
	mediaList, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scanner.Scan() failed: %v", err)
	}

	// Verify we found the music
	if len(mediaList) != 1 {
		t.Fatalf("Expected 1 media item, got %d", len(mediaList))
	}

	media := mediaList[0]
	if media.Type != Music {
		t.Errorf("Expected Music type, got %v", media.Type)
	}

	if media.MusicBrainzID != "c712a2bc-1e57-46c1-9432-8b9a33dc4159" {
		t.Errorf("Expected MusicBrainz ID to be set, got '%s'", media.MusicBrainzID)
	}

	// Verify tracks.json was created
	tracksPath := filepath.Join(musicDir, "tracks.json")
	if _, err := os.Stat(tracksPath); os.IsNotExist(err) {
		t.Fatal("tracks.json was not created during scan")
	}

	// Load and verify track list
	trackList, _ := NewMetadataLoader().LoadTrackList(media.Path)
	if trackList == nil {
		t.Fatal("Failed to load track list after scan")
	}

	if len(trackList.Media) == 0 {
		t.Fatal("Expected at least one medium in track list")
	}

	t.Logf("Successfully scanned music directory with MusicBrainz integration")
	t.Logf("  Media title: %s", media.Title)
	t.Logf("  Media count: %d", len(trackList.Media))
	t.Logf("  Track count: %d", len(trackList.Media[0].Tracks))
}

// TestIntegrationMultiDiscRelease tests multi-disc handling with a simulated release
// Note: This test makes real API calls to MusicBrainz and requires internet access
func TestIntegrationMultiDiscRelease(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	client := NewMusicBrainzClient()

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	// Test with Pink Floyd - The Wall (1994 2CD remaster)
	// MusicBrainz ID: 2437980f-513a-44c7-bfbe-3425f1c19d68
	mbid := "2437980f-513a-44c7-bfbe-3425f1c19d68"

	release, err := client.FetchRelease(mbid)
	if err != nil {
		// If this specific release is not found, skip the test
		// Multi-disc handling is tested in unit tests anyway
		t.Skipf("Multi-disc release not available: %v", err)
	}

	// Should have at least one medium
	if len(release.Media) == 0 {
		t.Fatal("Expected at least one medium")
	}

	// Convert to track list
	trackList := ConvertToTrackList(release)
	if trackList == nil {
		t.Fatal("Failed to convert to track list")
	}

	// Verify media are properly numbered
	for i, medium := range trackList.Media {
		expectedPosition := i + 1
		if medium.Position != expectedPosition {
			t.Errorf("Expected medium position %d, got %d", expectedPosition, medium.Position)
		}
		if len(medium.Tracks) == 0 {
			t.Errorf("Medium %d has no tracks", medium.Position)
		}
	}

	t.Logf("Successfully fetched release: %s", release.Title)
	t.Logf("  Disc count: %d", len(trackList.Media))
	for i, medium := range trackList.Media {
		t.Logf("  Disc %d: %d tracks (%s)", i+1, len(medium.Tracks), medium.Format)
	}
}

// TestIntegrationCachingBehavior tests that cached files are not re-downloaded
// Note: This test makes real API calls to MusicBrainz and requires internet access
func TestIntegrationCachingBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checkMusicBrainzAccess(t)

	// Rate limit: MusicBrainz recommends 1 request per second
	time.Sleep(1 * time.Second)

	tmpDir := t.TempDir()
	client := NewMusicBrainzClient()

	media := &Media{
		Title:         "Dark Side of the Moon",
		Type:          Music,
		MusicBrainzID: "c712a2bc-1e57-46c1-9432-8b9a33dc4159",
		Path:          tmpDir,
	}

	// First fetch - should download
	err := client.FetchAndSaveTrackList(media)
	if err != nil {
		t.Fatalf("First fetch failed: %v", err)
	}

	// Get file modification time
	tracksPath := filepath.Join(tmpDir, "tracks.json")
	info1, err := os.Stat(tracksPath)
	if err != nil {
		t.Fatalf("Failed to stat tracks.json: %v", err)
	}

	// Second fetch - should skip (cached)
	err = client.FetchAndSaveTrackList(media)
	if err != nil {
		t.Fatalf("Second fetch failed: %v", err)
	}

	// Verify file was not modified
	info2, err := os.Stat(tracksPath)
	if err != nil {
		t.Fatalf("Failed to stat tracks.json after second fetch: %v", err)
	}

	if !info1.ModTime().Equal(info2.ModTime()) {
		t.Error("File was modified on second fetch - caching not working")
	}

	t.Log("Caching behavior verified - existing files not overwritten")
}
