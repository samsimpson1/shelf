package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrMetadataNotFound is returned by MetadataLoader methods when the requested
// metadata file does not exist on disk. Callers can use errors.Is to
// distinguish absence from unexpected I/O or parse failures.
var ErrMetadataNotFound = errors.New("metadata file not found")

// MetadataLoader reads supplemental metadata (description, genres, track list,
// poster) from a media directory. It exists so the Media domain struct stays
// free of disk I/O and so handlers can decide how to degrade when reads fail.
type MetadataLoader struct{}

// NewMetadataLoader returns a new MetadataLoader.
func NewMetadataLoader() *MetadataLoader {
	return &MetadataLoader{}
}

// LoadDescription reads description.txt from mediaPath. Returns
// ErrMetadataNotFound if the file does not exist.
func (l *MetadataLoader) LoadDescription(mediaPath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(mediaPath, "description.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrMetadataNotFound
		}
		return "", fmt.Errorf("read description: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

// LoadGenres reads genre.txt from mediaPath and splits it on commas. Returns
// ErrMetadataNotFound if the file does not exist.
func (l *MetadataLoader) LoadGenres(mediaPath string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(mediaPath, "genre.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrMetadataNotFound
		}
		return nil, fmt.Errorf("read genres: %w", err)
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return []string{}, nil
	}
	parts := strings.Split(text, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts, nil
}

// LoadTrackList reads tracks.json from mediaPath. Returns ErrMetadataNotFound
// if the file does not exist.
func (l *MetadataLoader) LoadTrackList(mediaPath string) (*TrackList, error) {
	data, err := os.ReadFile(filepath.Join(mediaPath, "tracks.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrMetadataNotFound
		}
		return nil, fmt.Errorf("read track list: %w", err)
	}
	var tl TrackList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, fmt.Errorf("parse track list: %w", err)
	}
	return &tl, nil
}

// FindPoster locates the poster file in mediaPath. The bool indicates whether
// a poster was found. A non-nil error is only returned for unexpected I/O
// failures (not for "file does not exist").
func (l *MetadataLoader) FindPoster(mediaPath string) (string, bool, error) {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		p := filepath.Join(mediaPath, "poster"+ext)
		info, err := os.Stat(p)
		if err == nil {
			if info.IsDir() {
				continue
			}
			return p, true, nil
		}
		if !os.IsNotExist(err) {
			return "", false, fmt.Errorf("stat poster: %w", err)
		}
	}
	return "", false, nil
}
