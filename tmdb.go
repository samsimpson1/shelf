package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	tmdbAPIBaseURL = "https://api.themoviedb.org/3"
	tmdbImageBaseURL = "https://image.tmdb.org/t/p/original"
)

// TMDBClient handles interactions with the TMDB API
type TMDBClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewTMDBClient creates a new TMDB API client
func NewTMDBClient(apiKey string) *TMDBClient {
	return &TMDBClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Genre represents a genre from TMDB
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MovieResponse represents the TMDB API response for a movie
type MovieResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Overview    string  `json:"overview"`
	Genres      []Genre `json:"genres"`
}

// TVResponse represents the TMDB API response for a TV show
type TVResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	PosterPath   string  `json:"poster_path"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	Genres       []Genre `json:"genres"`
}

// FetchMovieMetadata fetches metadata for a movie from TMDB
func (c *TMDBClient) FetchMovieMetadata(movieID string) (*MovieResponse, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s", tmdbAPIBaseURL, movieID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d for movie %s", resp.StatusCode, movieID)
	}

	var movie MovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("failed to decode movie response: %w", err)
	}

	return &movie, nil
}

// FetchTVMetadata fetches metadata for a TV show from TMDB
func (c *TMDBClient) FetchTVMetadata(tvID string) (*TVResponse, error) {
	url := fmt.Sprintf("%s/tv/%s?api_key=%s", tmdbAPIBaseURL, tvID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TV metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d for TV show %s", resp.StatusCode, tvID)
	}

	var tv TVResponse
	if err := json.NewDecoder(resp.Body).Decode(&tv); err != nil {
		return nil, fmt.Errorf("failed to decode TV response: %w", err)
	}

	return &tv, nil
}

// DownloadPoster downloads a poster image from TMDB and saves it to the specified directory
func (c *TMDBClient) DownloadPoster(posterPath, destDir string) error {
	if posterPath == "" {
		return fmt.Errorf("poster path is empty")
	}

	// Remove leading slash if present
	posterPath = strings.TrimPrefix(posterPath, "/")

	// Construct the full image URL
	imageURL := fmt.Sprintf("%s/%s", tmdbImageBaseURL, posterPath)

	// Download the image
	resp, err := c.httpClient.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download poster: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download poster: status %d", resp.StatusCode)
	}

	// Determine file extension from the poster path
	ext := filepath.Ext(posterPath)
	if ext == "" {
		ext = ".jpg" // Default to jpg if no extension
	}

	// Create the destination file path
	destPath := filepath.Join(destDir, "poster"+ext)

	// Create the file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create poster file: %w", err)
	}
	defer file.Close()

	// Write the image data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write poster file: %w", err)
	}

	log.Printf("Downloaded poster to %s", destPath)
	return nil
}

// saveDescription saves the overview text to description.txt
func (c *TMDBClient) saveDescription(overview, destDir string) error {
	if overview == "" {
		return fmt.Errorf("overview is empty")
	}

	destPath := filepath.Join(destDir, "description.txt")

	// Write the description to file
	err := os.WriteFile(destPath, []byte(overview), 0644)
	if err != nil {
		return fmt.Errorf("failed to write description file: %w", err)
	}

	log.Printf("Saved description to %s", destPath)
	return nil
}

// saveGenres saves the genres as comma-separated values to genre.txt
func (c *TMDBClient) saveGenres(genres []Genre, destDir string) error {
	if len(genres) == 0 {
		return fmt.Errorf("no genres available")
	}

	// Extract genre names
	genreNames := make([]string, len(genres))
	for i, genre := range genres {
		genreNames[i] = genre.Name
	}

	// Join with commas
	genresText := strings.Join(genreNames, ", ")

	destPath := filepath.Join(destDir, "genre.txt")

	// Write the genres to file
	err := os.WriteFile(destPath, []byte(genresText), 0644)
	if err != nil {
		return fmt.Errorf("failed to write genre file: %w", err)
	}

	log.Printf("Saved genres to %s", destPath)
	return nil
}

// FetchAndSaveMetadata fetches metadata and downloads poster, description, and genres for a media item
func (c *TMDBClient) FetchAndSaveMetadata(media *Media) error {
	if media.TMDBID == "" {
		return fmt.Errorf("no TMDB ID for media: %s", media.Title)
	}

	// Check if all files already exist
	posterExists := false
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		posterPath := filepath.Join(media.Path, "poster"+ext)
		if _, err := os.Stat(posterPath); err == nil {
			posterExists = true
			break
		}
	}

	descriptionPath := filepath.Join(media.Path, "description.txt")
	_, descErr := os.Stat(descriptionPath)
	descriptionExists := descErr == nil

	genrePath := filepath.Join(media.Path, "genre.txt")
	_, genreErr := os.Stat(genrePath)
	genreExists := genreErr == nil

	// If all files exist, skip fetching
	if posterExists && descriptionExists && genreExists {
		log.Printf("All metadata files already exist for %s, skipping download", media.Title)
		return nil
	}

	var posterPath string
	var overview string
	var genres []Genre
	var err error

	// Fetch metadata based on media type
	if media.Type == Film {
		movie, err := c.FetchMovieMetadata(media.TMDBID)
		if err != nil {
			return fmt.Errorf("failed to fetch movie metadata: %w", err)
		}
		posterPath = movie.PosterPath
		overview = movie.Overview
		genres = movie.Genres
	} else if media.Type == TV {
		tv, err := c.FetchTVMetadata(media.TMDBID)
		if err != nil {
			return fmt.Errorf("failed to fetch TV metadata: %w", err)
		}
		posterPath = tv.PosterPath
		overview = tv.Overview
		genres = tv.Genres
	}

	// Download the poster if it doesn't exist
	if !posterExists {
		if posterPath == "" {
			log.Printf("Warning: No poster available for %s", media.Title)
		} else {
			if err = c.DownloadPoster(posterPath, media.Path); err != nil {
				log.Printf("Warning: Failed to download poster for %s: %v", media.Title, err)
			}
		}
	}

	// Save description if it doesn't exist
	if !descriptionExists {
		if overview == "" {
			log.Printf("Warning: No overview available for %s", media.Title)
		} else {
			if err = c.saveDescription(overview, media.Path); err != nil {
				log.Printf("Warning: Failed to save description for %s: %v", media.Title, err)
			}
		}
	}

	// Save genres if they don't exist
	if !genreExists {
		if len(genres) == 0 {
			log.Printf("Warning: No genres available for %s", media.Title)
		} else {
			if err = c.saveGenres(genres, media.Path); err != nil {
				log.Printf("Warning: Failed to save genres for %s: %v", media.Title, err)
			}
		}
	}

	return nil
}

// FetchAndSavePoster is deprecated, use FetchAndSaveMetadata instead
// Kept for backward compatibility
func (c *TMDBClient) FetchAndSavePoster(media *Media) error {
	return c.FetchAndSaveMetadata(media)
}
