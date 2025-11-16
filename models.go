package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MediaType represents the type of media (Film, TV, or Music)
type MediaType int

const (
	Film MediaType = iota
	TV
	Music
)

// String returns the string representation of MediaType
func (m MediaType) String() string {
	switch m {
	case Film:
		return "Film"
	case TV:
		return "TV"
	case Music:
		return "Music"
	default:
		return "Unknown"
	}
}

// Disk represents an individual disk in a media backup
type Disk struct {
	Name   string  // Disk name/identifier (e.g., "Disk 1", "Series 1 Disk 2")
	Format string  // Disk format (e.g., "Blu-Ray", "DVD", "Blu-Ray UHD")
	SizeGB float64 // Disk size in gigabytes
	Path   string  // Absolute path to the disk directory
}

// Track represents an individual track on a music release
type Track struct {
	Position string `json:"position"` // Track position (e.g., "1", "A1", "2.3")
	Title    string `json:"title"`    // Track title
	Duration int    `json:"duration"` // Duration in milliseconds
}

// FormatDuration returns the track duration formatted as mm:ss
func (t *Track) FormatDuration() string {
	if t.Duration == 0 {
		return ""
	}
	totalSeconds := t.Duration / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// Medium represents a disc/medium in a music release
type Medium struct {
	Position int     `json:"position"` // Medium position (disc number)
	Format   string  `json:"format"`   // Format (e.g., "CD", "Vinyl", "Digital Media")
	Tracks   []Track `json:"tracks"`   // Tracks on this medium
}

// TrackList represents the complete track listing for a music release
type TrackList struct {
	Media []Medium `json:"media"` // All media/discs in the release
}

// Media represents a media item from the backup directory
type Media struct {
	Title          string    // Title of the media
	Type           MediaType // Film, TV, or Music
	Year           int       // Year (for films, 0 for TV shows and Music)
	DiskCount      int       // Number of disks
	Disks          []Disk    // Individual disk information
	TMDBID         string    // TMDB ID (optional, empty string if not present; Music does not use TMDB)
	MusicBrainzID  string    // MusicBrainz release ID (optional, only for Music media type)
	Path           string    // Absolute path to the media directory
}

// DisplayTitle returns the title with year for films, just title for TV and Music
func (m *Media) DisplayTitle() string {
	if m.Type == Film && m.Year > 0 {
		return fmt.Sprintf("%s (%d)", m.Title, m.Year)
	}
	return m.Title
}

// Slug generates a URL-friendly slug from the media title and year
func (m *Media) Slug() string {
	slug := strings.ToLower(m.Title)
	// Replace non-alphanumeric characters with hyphens
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	// Add year for films to make slugs unique
	if m.Type == Film && m.Year > 0 {
		slug = fmt.Sprintf("%s-%d", slug, m.Year)
	}
	return slug
}

// FindPosterFile returns the path and extension of the poster file if it exists
func (m *Media) FindPosterFile() (string, bool) {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		posterPath := filepath.Join(m.Path, "poster"+ext)
		if _, err := os.Stat(posterPath); err == nil {
			return posterPath, true
		}
	}
	return "", false
}

// PosterURL returns the relative URL for the poster image
func (m *Media) PosterURL() string {
	return fmt.Sprintf("/posters/%s", m.Slug())
}

// LoadDescription reads and returns the description from description.txt
func (m *Media) LoadDescription() string {
	descPath := filepath.Join(m.Path, "description.txt")
	data, err := os.ReadFile(descPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// LoadGenres reads and returns the genres from genre.txt as a slice
func (m *Media) LoadGenres() []string {
	genrePath := filepath.Join(m.Path, "genre.txt")
	data, err := os.ReadFile(genrePath)
	if err != nil {
		return []string{}
	}
	genresText := strings.TrimSpace(string(data))
	if genresText == "" {
		return []string{}
	}
	// Split by comma and trim whitespace
	genres := strings.Split(genresText, ",")
	for i := range genres {
		genres[i] = strings.TrimSpace(genres[i])
	}
	return genres
}

// LoadTrackList reads and returns the track list from tracks.json
func (m *Media) LoadTrackList() *TrackList {
	tracksPath := filepath.Join(m.Path, "tracks.json")
	data, err := os.ReadFile(tracksPath)
	if err != nil {
		return nil
	}

	var trackList TrackList
	if err := json.Unmarshal(data, &trackList); err != nil {
		return nil
	}

	return &trackList
}

// PlayCommand generates a VLC play command for the disk
func (d *Disk) PlayCommand(prefix string) string {
	// Determine protocol based on disk format
	var protocol string
	formatLower := strings.ToLower(d.Format)

	if strings.Contains(formatLower, "blu-ray") || strings.Contains(formatLower, "bluray") {
		protocol = "bluray://"
	} else if strings.Contains(formatLower, "dvd") {
		protocol = "dvd://"
	} else {
		protocol = "file://"
	}

	// Construct the full path
	fullPath := d.Path
	if prefix != "" {
		fullPath = prefix + d.Path
	}

	return fmt.Sprintf("vlc \"%s%s\"", protocol, fullPath)
}

// MPVPlayCommand generates an MPV play command for the disk
func (d *Disk) MPVPlayCommand(prefix string) string {
	// Construct the full path
	fullPath := d.Path
	if prefix != "" {
		fullPath = prefix + d.Path
	}

	// Determine command based on disk format
	formatLower := strings.ToLower(d.Format)

	if strings.Contains(formatLower, "blu-ray") || strings.Contains(formatLower, "bluray") {
		return fmt.Sprintf("mpv bd:// --bluray-device=\"%s\"", fullPath)
	} else if strings.Contains(formatLower, "dvd") {
		return fmt.Sprintf("mpv dvd:// --dvd-device=\"%s\"", fullPath)
	} else {
		return fmt.Sprintf("mpv \"%s\"", fullPath)
	}
}
