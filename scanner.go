package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regex patterns for parsing directory names
	filmPattern = regexp.MustCompile(`^(.+) \((\d{4})\) \[Film\]$`)
	tvPattern   = regexp.MustCompile(`^(.+) \[TV\]$`)
	diskPattern = regexp.MustCompile(`^Disk \[.+\]$`)
	tvDiskPattern = regexp.MustCompile(`^Series (\d+) Disk (\d+) \[.+\]$`)
)

// Scanner scans a directory for media items
type Scanner struct {
	mediaDir   string
	tmdbClient *TMDBClient
}

// NewScanner creates a new Scanner for the given directory
func NewScanner(mediaDir string) *Scanner {
	return &Scanner{mediaDir: mediaDir}
}

// NewScannerWithTMDB creates a new Scanner with TMDB client for poster fetching
func NewScannerWithTMDB(mediaDir string, tmdbClient *TMDBClient) *Scanner {
	return &Scanner{
		mediaDir:   mediaDir,
		tmdbClient: tmdbClient,
	}
}

// Scan scans the configured directory and returns a slice of Media items
func (s *Scanner) Scan() ([]Media, error) {
	// Verify directory exists and is readable
	info, err := os.Stat(s.mediaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("media directory does not exist: %s", s.mediaDir)
		}
		return nil, fmt.Errorf("cannot access media directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("media path is not a directory: %s", s.mediaDir)
	}

	// Read directory entries
	entries, err := os.ReadDir(s.mediaDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read media directory: %w", err)
	}

	var mediaList []Media

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		dirPath := filepath.Join(s.mediaDir, dirName)

		// Try to parse as film
		if media, ok := s.parseFilm(dirName, dirPath); ok {
			mediaList = append(mediaList, media)
			continue
		}

		// Try to parse as TV show
		if media, ok := s.parseTV(dirName, dirPath); ok {
			mediaList = append(mediaList, media)
			continue
		}

		// If neither pattern matches, skip this directory
	}

	return mediaList, nil
}

// parseFilm attempts to parse a directory as a film
func (s *Scanner) parseFilm(dirName, dirPath string) (Media, bool) {
	matches := filmPattern.FindStringSubmatch(dirName)
	if matches == nil {
		return Media{}, false
	}

	title := matches[1]
	year, _ := strconv.Atoi(matches[2]) // We know it's a valid 4-digit number from regex

	media := Media{
		Title: title,
		Type:  Film,
		Year:  year,
		Path:  dirPath,
	}

	// Collect disk details
	media.Disks = s.collectFilmDisks(dirPath)
	media.DiskCount = len(media.Disks)

	// Read TMDB ID if present
	media.TMDBID = s.readTMDBID(dirPath)

	// Fetch metadata if TMDB client is configured
	if s.tmdbClient != nil && media.TMDBID != "" {
		if err := s.tmdbClient.FetchAndSaveMetadata(&media); err != nil {
			log.Printf("Warning: Failed to fetch metadata for %s: %v", media.Title, err)
		}
	}

	return media, true
}

// parseTV attempts to parse a directory as a TV show
func (s *Scanner) parseTV(dirName, dirPath string) (Media, bool) {
	matches := tvPattern.FindStringSubmatch(dirName)
	if matches == nil {
		return Media{}, false
	}

	title := matches[1]

	media := Media{
		Title: title,
		Type:  TV,
		Year:  0, // TV shows don't have years in their directory names
		Path:  dirPath,
	}

	// Collect disk details
	media.Disks = s.collectTVDisks(dirPath)
	media.DiskCount = len(media.Disks)

	// Read TMDB ID if present
	media.TMDBID = s.readTMDBID(dirPath)

	// Fetch metadata if TMDB client is configured
	if s.tmdbClient != nil && media.TMDBID != "" {
		if err := s.tmdbClient.FetchAndSaveMetadata(&media); err != nil {
			log.Printf("Warning: Failed to fetch metadata for %s: %v", media.Title, err)
		}
	}

	return media, true
}

// calculateDirSize calculates the total size of a directory in bytes
func calculateDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// loadSizeCache loads the disk size cache from sizes.json in the media directory
// Returns an empty map if the file doesn't exist or contains invalid JSON
func loadSizeCache(mediaDir string) map[string]int64 {
	cachePath := filepath.Join(mediaDir, "sizes.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		// Cache file doesn't exist or can't be read - return empty map
		return make(map[string]int64)
	}

	var cache map[string]int64
	if err := json.Unmarshal(data, &cache); err != nil {
		// Invalid JSON - return empty map
		return make(map[string]int64)
	}

	return cache
}

// saveSizeCache saves the disk size cache to sizes.json in the media directory
func saveSizeCache(mediaDir string, cache map[string]int64) error {
	cachePath := filepath.Join(mediaDir, "sizes.json")

	// Marshal with indentation for human readability
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	// Write the cache file
	return os.WriteFile(cachePath, data, 0644)
}

// extractFormat extracts the disk format from directory name (content within brackets)
func extractFormat(dirName string) string {
	// Find content within square brackets
	start := strings.Index(dirName, "[")
	end := strings.Index(dirName, "]")
	if start >= 0 && end > start {
		return dirName[start+1 : end]
	}
	return ""
}

// countFilmDisks counts the number of disk directories in a film directory
func (s *Scanner) countFilmDisks(dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if diskPattern.MatchString(entry.Name()) {
			count++
		}
	}

	return count
}

// collectFilmDisks collects detailed information about film disks
func (s *Scanner) collectFilmDisks(dirPath string) []Disk {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []Disk{}
	}

	// Load size cache
	cache := loadSizeCache(dirPath)
	cacheUpdated := false

	var disks []Disk
	diskNum := 1
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if diskPattern.MatchString(entry.Name()) {
			format := extractFormat(entry.Name())
			diskPath := filepath.Join(dirPath, entry.Name())

			// Try to get size from cache first
			var size int64
			diskDirName := entry.Name()
			if cachedSize, exists := cache[diskDirName]; exists {
				size = cachedSize
			} else {
				// Calculate size and update cache
				var err error
				size, err = calculateDirSize(diskPath)
				if err == nil {
					cache[diskDirName] = size
					cacheUpdated = true
				}
			}

			sizeGB := float64(size) / (1024 * 1024 * 1024) // Convert bytes to GB

			disks = append(disks, Disk{
				Name:   fmt.Sprintf("Disk %d", diskNum),
				Format: format,
				SizeGB: sizeGB,
				Path:   diskPath,
			})
			diskNum++
		}
	}

	// Save cache if it was updated
	if cacheUpdated {
		if err := saveSizeCache(dirPath, cache); err != nil {
			log.Printf("Warning: Failed to save size cache for %s: %v", dirPath, err)
		}
	}

	return disks
}

// countTVDisks counts the number of disk directories in a TV show directory
func (s *Scanner) countTVDisks(dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if tvDiskPattern.MatchString(entry.Name()) {
			count++
		}
	}

	return count
}

// collectTVDisks collects detailed information about TV show disks
func (s *Scanner) collectTVDisks(dirPath string) []Disk {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []Disk{}
	}

	// Load size cache
	cache := loadSizeCache(dirPath)
	cacheUpdated := false

	var disks []Disk
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matches := tvDiskPattern.FindStringSubmatch(entry.Name())
		if matches != nil {
			seriesNum := matches[1]
			diskNum := matches[2]
			format := extractFormat(entry.Name())
			diskPath := filepath.Join(dirPath, entry.Name())

			// Try to get size from cache first
			var size int64
			diskDirName := entry.Name()
			if cachedSize, exists := cache[diskDirName]; exists {
				size = cachedSize
			} else {
				// Calculate size and update cache
				var err error
				size, err = calculateDirSize(diskPath)
				if err == nil {
					cache[diskDirName] = size
					cacheUpdated = true
				}
			}

			sizeGB := float64(size) / (1024 * 1024 * 1024) // Convert bytes to GB

			disks = append(disks, Disk{
				Name:   fmt.Sprintf("Series %s Disk %s", seriesNum, diskNum),
				Format: format,
				SizeGB: sizeGB,
				Path:   diskPath,
			})
		}
	}

	// Save cache if it was updated
	if cacheUpdated {
		if err := saveSizeCache(dirPath, cache); err != nil {
			log.Printf("Warning: Failed to save size cache for %s: %v", dirPath, err)
		}
	}

	return disks
}

// readTMDBID reads the TMDB ID from tmdb.txt file if it exists
func (s *Scanner) readTMDBID(dirPath string) string {
	tmdbPath := filepath.Join(dirPath, "tmdb.txt")

	data, err := os.ReadFile(tmdbPath)
	if err != nil {
		return "" // File doesn't exist or can't be read
	}

	// Trim whitespace and return
	return strings.TrimSpace(string(data))
}
