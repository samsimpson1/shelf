package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	musicBrainzAPIBaseURL = "https://musicbrainz.org/ws/2"
	// MusicBrainz requires a User-Agent header with app name, version, and contact
	// Format: AppName/Version ( contact-url )
	musicBrainzUserAgent = "Shelf-MediaBackupManager/1.0 ( https://github.com/samsimpson1/shelf )"
)

// MusicBrainzClient handles interactions with the MusicBrainz API
type MusicBrainzClient struct {
	httpClient *http.Client
}

// NewMusicBrainzClient creates a new MusicBrainz API client
func NewMusicBrainzClient() *MusicBrainzClient {
	return &MusicBrainzClient{
		httpClient: &http.Client{},
	}
}

// MBTrack represents a track from MusicBrainz API
type MBTrack struct {
	ID        string      `json:"id"`
	Position  int         `json:"position"`
	Title     string      `json:"title"`
	Length    int         `json:"length"` // Duration in milliseconds
	Recording MBRecording `json:"recording"`
}

// MBRecording represents a recording from MusicBrainz API
type MBRecording struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Length int    `json:"length"` // Duration in milliseconds
}

// MBMedium represents a medium/disc from MusicBrainz API
type MBMedium struct {
	Position int       `json:"position"`
	Format   string    `json:"format"`
	Tracks   []MBTrack `json:"tracks"`
}

// MBRelease represents a release from MusicBrainz API
type MBRelease struct {
	ID    string     `json:"id"`
	Title string     `json:"title"`
	Media []MBMedium `json:"media"`
}

// MBReleaseSearchResult represents a release search result from MusicBrainz
type MBReleaseSearchResult struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"artist-credit-phrase"` // Flattened artist string
	Date       string `json:"date"`
	Country    string `json:"country"`
	TrackCount int    `json:"track-count"`
	Score      int    `json:"score"` // Search relevance score (0-100)
}

// MBSearchResponse represents the MusicBrainz search response
type MBSearchResponse struct {
	Created  string                  `json:"created"`
	Count    int                     `json:"count"`
	Offset   int                     `json:"offset"`
	Releases []MBReleaseSearchResult `json:"releases"`
}

// FetchRelease fetches release information including track list from MusicBrainz
func (c *MusicBrainzClient) FetchRelease(mbid string) (*MBRelease, error) {
	if mbid == "" {
		return nil, fmt.Errorf("MusicBrainz ID cannot be empty")
	}

	// Build URL with recordings included
	url := fmt.Sprintf("%s/release/%s?inc=recordings&fmt=json", musicBrainzAPIBaseURL, mbid)

	// Create request with required User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", musicBrainzUserAgent)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MusicBrainz API returned status %d for release %s", resp.StatusCode, mbid)
	}

	var release MBRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release response: %w", err)
	}

	return &release, nil
}

// SearchReleases searches for releases on MusicBrainz by query
func (c *MusicBrainzClient) SearchReleases(query string) ([]MBReleaseSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Build URL with query parameter (URL-encoded)
	// We search for releases and include artist credits
	url := fmt.Sprintf("%s/release?query=%s&fmt=json", musicBrainzAPIBaseURL, url.QueryEscape(query))

	// Create request with required User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", musicBrainzUserAgent)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MusicBrainz API returned status %d for search", resp.StatusCode)
	}

	var searchResp MBSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Limit to 20 results
	results := searchResp.Releases
	if len(results) > 20 {
		results = results[:20]
	}

	return results, nil
}

// ConvertToTrackList converts MusicBrainz release data to our TrackList format
func ConvertToTrackList(mbRelease *MBRelease) *TrackList {
	if mbRelease == nil {
		return nil
	}

	trackList := &TrackList{
		Media: make([]Medium, len(mbRelease.Media)),
	}

	for i, mbMedium := range mbRelease.Media {
		medium := Medium{
			Position: mbMedium.Position,
			Format:   mbMedium.Format,
			Tracks:   make([]Track, len(mbMedium.Tracks)),
		}

		for j, mbTrack := range mbMedium.Tracks {
			// Use track length, fall back to recording length if not available
			duration := mbTrack.Length
			if duration == 0 && mbTrack.Recording.Length > 0 {
				duration = mbTrack.Recording.Length
			}

			// Use recording title if track title is empty
			title := mbTrack.Title
			if title == "" {
				title = mbTrack.Recording.Title
			}

			medium.Tracks[j] = Track{
				Position: fmt.Sprintf("%d", mbTrack.Position),
				Title:    title,
				Duration: duration,
			}
		}

		trackList.Media[i] = medium
	}

	return trackList
}

// SaveTrackList saves a track list to tracks.json in the specified directory
func SaveTrackList(trackList *TrackList, destDir string) error {
	if trackList == nil {
		return fmt.Errorf("track list is nil")
	}

	destPath := filepath.Join(destDir, "tracks.json")

	// Marshal with indentation for human readability
	data, err := json.MarshalIndent(trackList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal track list: %w", err)
	}

	// Write the track list to file
	err = os.WriteFile(destPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write tracks.json: %w", err)
	}

	log.Printf("Saved track list to %s", destPath)
	return nil
}

// FetchAndSaveTrackList fetches track list from MusicBrainz and saves it to tracks.json
func (c *MusicBrainzClient) FetchAndSaveTrackList(media *Media) error {
	if media.MusicBrainzID == "" {
		return fmt.Errorf("no MusicBrainz ID for media: %s", media.Title)
	}

	// Check if tracks.json already exists (caching behavior)
	tracksPath := filepath.Join(media.Path, "tracks.json")
	if _, err := os.Stat(tracksPath); err == nil {
		log.Printf("Track list already exists for %s, skipping download", media.Title)
		return nil
	}

	// Fetch release data from MusicBrainz
	release, err := c.FetchRelease(media.MusicBrainzID)
	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}

	// Convert to our track list format
	trackList := ConvertToTrackList(release)
	if trackList == nil {
		return fmt.Errorf("failed to convert track list")
	}

	// Save to tracks.json
	if err := SaveTrackList(trackList, media.Path); err != nil {
		return fmt.Errorf("failed to save track list: %w", err)
	}

	return nil
}
