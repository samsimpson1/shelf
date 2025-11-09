package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ImportDirectory represents a directory available for import
type ImportDirectory struct {
	Name   string  // Directory name
	Path   string  // Absolute path
	SizeGB float64 // Total size in gigabytes
}

// DiskType represents the type of disk backup
type DiskType string

const (
	DiskTypeBluRay    DiskType = "Blu-Ray"
	DiskTypeBluRayUHD DiskType = "Blu-Ray UHD"
	DiskTypeDVD       DiskType = "DVD"
	DiskTypeCustom    DiskType = "" // Custom text entered by user
)

// String returns the string representation of DiskType
func (d DiskType) String() string {
	return string(d)
}

// ImportSession represents an ongoing import workflow
type ImportSession struct {
	// Source information
	SourceDir   *ImportDirectory // Directory being imported
	DetectedType DiskType        // Auto-detected disk type

	// User selections
	MediaKind   MediaType // Film or TV
	Title       string    // Media title
	Year        int       // Year (for films)
	TMDBID      string    // TMDB ID (optional)
	SeriesNum   int       // Series number (for TV)
	DiskNum     int       // Disk number (for TV) or film disk number
	DiskType    DiskType  // Selected disk type
	DiskTypeCustom string // Custom disk type text
	AddToExisting bool    // Add to existing media vs create new
	ExistingMediaPath string // Path to existing media (if adding)

	// Metadata from TMDB (if selected)
	TMDBTitle   string   // Official title from TMDB
	TMDBYear    int      // Year from TMDB (for films)
	TMDBOverview string  // Description
	TMDBGenres  []string // Genres
}

// ImportScanner scans the import directory for available imports
type ImportScanner struct {
	importDir string
}

// NewImportScanner creates a new ImportScanner for the given directory
func NewImportScanner(importDir string) *ImportScanner {
	return &ImportScanner{importDir: importDir}
}

// Scan scans the import directory and returns a list of ImportDirectory items
func (s *ImportScanner) Scan() ([]ImportDirectory, error) {
	// Verify directory exists and is readable
	info, err := os.Stat(s.importDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("import directory does not exist: %s", s.importDir)
		}
		return nil, fmt.Errorf("cannot access import directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("import path is not a directory: %s", s.importDir)
	}

	// Read directory entries
	entries, err := os.ReadDir(s.importDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read import directory: %w", err)
	}

	var imports []ImportDirectory

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(s.importDir, entry.Name())

		// Calculate directory size
		size, err := calculateDirSize(dirPath)
		if err != nil {
			size = 0 // If we can't calculate size, just use 0
		}
		sizeGB := float64(size) / (1024 * 1024 * 1024)

		imports = append(imports, ImportDirectory{
			Name:   entry.Name(),
			Path:   dirPath,
			SizeGB: sizeGB,
		})
	}

	return imports, nil
}

// DetectDiskType attempts to detect the disk type from directory structure
// Returns the detected type and a confidence boolean
func DetectDiskType(dirPath string) (DiskType, bool) {
	// Check for BDMV directory (Blu-ray)
	bdmvPath := filepath.Join(dirPath, "BDMV")
	if info, err := os.Stat(bdmvPath); err == nil && info.IsDir() {
		// Check for UHD indicators
		// UHD Blu-rays often have larger STREAM files or specific resolution indicators
		// For now, we'll just return regular Blu-Ray
		// TODO: Add UHD detection logic if needed
		return DiskTypeBluRay, true
	}

	// Check for VIDEO_TS directory (DVD)
	videoTSPath := filepath.Join(dirPath, "VIDEO_TS")
	if info, err := os.Stat(videoTSPath); err == nil && info.IsDir() {
		return DiskTypeDVD, true
	}

	// No recognized structure
	return "", false
}

// SanitizeName sanitizes a filename or directory name for filesystem compatibility
// Replaces problematic characters while maintaining readability
func SanitizeName(name string) string {
	// Replace colons with underscores (problematic on many filesystems)
	name = strings.ReplaceAll(name, ":", "_")

	// Replace forward slashes and backslashes
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")

	// Replace other problematic characters
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "\"", "'")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "*", "_")

	// Remove any control characters (ASCII 0-31)
	result := strings.Map(func(r rune) rune {
		if r < 32 {
			return -1 // Remove character
		}
		return r
	}, name)

	// Trim whitespace
	result = strings.TrimSpace(result)

	// Replace multiple consecutive underscores with single underscore
	multiUnderscore := regexp.MustCompile(`_+`)
	result = multiUnderscore.ReplaceAllString(result, "_")

	return result
}

// GenerateMediaDirName generates the directory name for a media item
// Follows the naming convention: "Title (Year) [Film]" or "Title [TV]"
func GenerateMediaDirName(title string, year int, mediaType MediaType) string {
	sanitizedTitle := SanitizeName(title)

	if mediaType == Film {
		return fmt.Sprintf("%s (%d) [Film]", sanitizedTitle, year)
	}
	return fmt.Sprintf("%s [TV]", sanitizedTitle)
}

// GenerateDiskDirName generates the directory name for a disk
// Follows the naming convention: "Disk [Format]" or "Series X Disk Y [Format]"
func GenerateDiskDirName(diskType string, seriesNum, diskNum int, mediaType MediaType) string {
	sanitizedType := SanitizeName(diskType)

	if mediaType == Film {
		return fmt.Sprintf("Disk [%s]", sanitizedType)
	}
	return fmt.Sprintf("Series %d Disk %d [%s]", seriesNum, diskNum, sanitizedType)
}

// ExecuteImport performs the actual import operation
// Moves the source directory to the destination with validation
func ExecuteImport(session *ImportSession, mediaDir string) error {
	if session == nil {
		return fmt.Errorf("import session is nil")
	}
	if session.SourceDir == nil {
		return fmt.Errorf("source directory is nil")
	}

	// Determine the final title (prefer TMDB title if available)
	finalTitle := session.Title
	if session.TMDBTitle != "" {
		finalTitle = session.TMDBTitle
	}

	// Determine the final year (prefer TMDB year if available)
	finalYear := session.Year
	if session.TMDBYear > 0 {
		finalYear = session.TMDBYear
	}

	// Determine disk type text
	diskTypeText := session.DiskType.String()
	if session.DiskType == DiskTypeCustom && session.DiskTypeCustom != "" {
		diskTypeText = session.DiskTypeCustom
	}

	var destMediaPath string
	var destDiskPath string

	if session.AddToExisting {
		// Add to existing media
		if session.ExistingMediaPath == "" {
			return fmt.Errorf("existing media path is required when adding to existing media")
		}
		destMediaPath = session.ExistingMediaPath

		// Generate disk directory name
		diskDirName := GenerateDiskDirName(diskTypeText, session.SeriesNum, session.DiskNum, session.MediaKind)
		destDiskPath = filepath.Join(destMediaPath, diskDirName)
	} else {
		// Create new media
		mediaDirName := GenerateMediaDirName(finalTitle, finalYear, session.MediaKind)
		destMediaPath = filepath.Join(mediaDir, mediaDirName)

		// Generate disk directory name
		diskDirName := GenerateDiskDirName(diskTypeText, session.SeriesNum, session.DiskNum, session.MediaKind)
		destDiskPath = filepath.Join(destMediaPath, diskDirName)
	}

	// Validate that destination doesn't already exist
	if _, err := os.Stat(destDiskPath); err == nil {
		return fmt.Errorf("destination already exists: %s", destDiskPath)
	}

	// Create media directory if it doesn't exist
	if _, err := os.Stat(destMediaPath); os.IsNotExist(err) {
		if err := os.MkdirAll(destMediaPath, 0755); err != nil {
			return fmt.Errorf("failed to create media directory: %w", err)
		}
	}

	// Move the source directory to the destination
	if err := os.Rename(session.SourceDir.Path, destDiskPath); err != nil {
		return fmt.Errorf("failed to move directory: %w", err)
	}

	// Write TMDB ID if provided
	if session.TMDBID != "" {
		tmdbPath := filepath.Join(destMediaPath, "tmdb.txt")
		// Only write if it doesn't exist (don't overwrite existing TMDB ID)
		if _, err := os.Stat(tmdbPath); os.IsNotExist(err) {
			if err := os.WriteFile(tmdbPath, []byte(session.TMDBID), 0644); err != nil {
				// Log warning but don't fail the import
				fmt.Printf("Warning: Failed to write TMDB ID: %v\n", err)
			}
		}
	}

	return nil
}
