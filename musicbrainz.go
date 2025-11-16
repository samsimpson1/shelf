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
	ID       string      `json:"id"`
	Position int         `json:"position"`
	Title    string      `json:"title"`
	Length   int         `json:"length"` // Duration in milliseconds
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
	Date  string     `json:"date"`  // Release date (YYYY-MM-DD)
	Media []MBMedium `json:"media"`
}

// MBArtistCredit represents an artist credit from MusicBrainz API
type MBArtistCredit struct {
	Name   string   `json:"name"`
	Artist MBArtist `json:"artist"`
}

// MBArtist represents an artist from MusicBrainz API
type MBArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MBReleaseSearchResult represents a single release in search results
type MBReleaseSearchResult struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Date          string             `json:"date"`  // Release date
	Status        string             `json:"status"` // Official, Promotion, etc.
	Country       string             `json:"country"`
	ArtistCredit  []MBArtistCredit   `json:"artist-credit"`
	ReleaseGroup  MBReleaseGroup     `json:"release-group"`
	Media         []MBMedium         `json:"media"`
}

// MBReleaseGroup represents a release group from MusicBrainz API
type MBReleaseGroup struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	PrimaryType  string `json:"primary-type"`  // Album, Single, EP, etc.
}

// MBReleaseSearchResponse represents the response from release search
type MBReleaseSearchResponse struct {
	Created  string                  `json:"created"`
	Count    int                     `json:"count"`
	Offset   int                     `json:"offset"`
	Releases []MBReleaseSearchResult `json:"releases"`
}

// GetArtistNames returns a comma-separated list of artist names
func (r *MBReleaseSearchResult) GetArtistNames() string {
	if len(r.ArtistCredit) == 0 {
		return "Unknown Artist"
	}

	var names []string
	for _, credit := range r.ArtistCredit {
		if credit.Name != "" {
			names = append(names, credit.Name)
		}
	}

	if len(names) == 0 {
		return "Unknown Artist"
	}

	return formatArtistNames(names)
}

// formatArtistNames formats a list of artist names with proper separators
func formatArtistNames(names []string) string {
	if len(names) == 0 {
		return ""
	}
	if len(names) == 1 {
		return names[0]
	}
	if len(names) == 2 {
		return names[0] + " & " + names[1]
	}
	// For 3+ artists, use commas and & for the last one
	result := ""
	for i, name := range names {
		if i == len(names)-1 {
			result += " & " + name
		} else if i > 0 {
			result += ", " + name
		} else {
			result += name
		}
	}
	return result
}

// GetTrackCount returns the total number of tracks across all media
func (r *MBReleaseSearchResult) GetTrackCount() int {
	total := 0
	for _, medium := range r.Media {
		total += len(medium.Tracks)
	}
	return total
}

// GetFormat returns the format of the first medium, or "Various" if multiple formats
func (r *MBReleaseSearchResult) GetFormat() string {
	if len(r.Media) == 0 {
		return "Unknown"
	}

	firstFormat := r.Media[0].Format
	if firstFormat == "" {
		return "Unknown"
	}

	// Check if all media have the same format
	for _, medium := range r.Media {
		if medium.Format != firstFormat {
			return "Various"
		}
	}

	return firstFormat
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

// SearchReleases searches for releases on MusicBrainz by artist and/or title
// Returns up to 20 results, sorted by relevance
func (c *MusicBrainzClient) SearchReleases(artist, title string) ([]MBReleaseSearchResult, error) {
	if artist == "" && title == "" {
		return nil, fmt.Errorf("at least one of artist or title must be provided")
	}

	// Build search query using Lucene syntax
	// MusicBrainz search API: https://musicbrainz.org/doc/MusicBrainz_API/Search
	queryParts := []string{}
	if artist != "" {
		queryParts = append(queryParts, fmt.Sprintf("artist:%q", artist))
	}
	if title != "" {
		queryParts = append(queryParts, fmt.Sprintf("release:%q", title))
	}

	query := ""
	if len(queryParts) > 0 {
		query = queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			query += " AND " + queryParts[i]
		}
	}

	// Build URL with search query, limit to 20 results, and include media/artist-credits
	url := fmt.Sprintf("%s/release?query=%s&limit=20&inc=media+artist-credits+release-groups&fmt=json",
		musicBrainzAPIBaseURL,
		escapeQueryString(query))

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
		return nil, fmt.Errorf("MusicBrainz API returned status %d", resp.StatusCode)
	}

	var searchResp MBReleaseSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return searchResp.Releases, nil
}

// escapeQueryString properly escapes a query string for URL encoding
func escapeQueryString(query string) string {
	return url.QueryEscape(query)
}
