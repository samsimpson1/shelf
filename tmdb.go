package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

// MovieSearchResult represents a movie search result from TMDB
type MovieSearchResult struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"poster_path"`
	Popularity  float64 `json:"popularity"`
}

// GetTitle returns the title for template rendering
func (m MovieSearchResult) GetTitle() string {
	return m.Title
}

// GetDate returns the release date for template rendering
func (m MovieSearchResult) GetDate() string {
	return m.ReleaseDate
}

// TVSearchResult represents a TV show search result from TMDB
type TVSearchResult struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	Popularity   float64 `json:"popularity"`
}

// GetTitle returns the name for template rendering (matches MovieSearchResult)
func (t TVSearchResult) GetTitle() string {
	return t.Name
}

// GetDate returns the first air date for template rendering (matches MovieSearchResult)
func (t TVSearchResult) GetDate() string {
	return t.FirstAirDate
}

// MovieSearchResponse represents the TMDB API response for movie search
type MovieSearchResponse struct {
	Results []MovieSearchResult `json:"results"`
}

// TVSearchResponse represents the TMDB API response for TV search
type TVSearchResponse struct {
	Results []TVSearchResult `json:"results"`
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

// SearchMovies searches for movies on TMDB by title and optional year
// Returns up to 20 results sorted by popularity
func (c *TMDBClient) SearchMovies(query string, year int) ([]MovieSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Build URL with query parameter (URL-encoded)
	searchURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", tmdbAPIBaseURL, c.apiKey, url.QueryEscape(query))

	// Add year parameter if provided
	if year > 0 {
		searchURL = fmt.Sprintf("%s&year=%d", searchURL, year)
	}

	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d for movie search", resp.StatusCode)
	}

	var searchResp MovieSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode movie search response: %w", err)
	}

	// Limit to 20 results
	results := searchResp.Results
	if len(results) > 20 {
		results = results[:20]
	}

	return results, nil
}

// SearchTV searches for TV shows on TMDB by name
// Returns up to 20 results sorted by popularity
func (c *TMDBClient) SearchTV(query string) ([]TVSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	searchURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s", tmdbAPIBaseURL, c.apiKey, url.QueryEscape(query))

	resp, err := c.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API returned status %d for TV search", resp.StatusCode)
	}

	var searchResp TVSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode TV search response: %w", err)
	}

	// Limit to 20 results
	results := searchResp.Results
	if len(results) > 20 {
		results = results[:20]
	}

	return results, nil
}

// WriteTMDBID writes a TMDB ID to a tmdb.txt file in the specified directory
// Creates or overwrites the file with permissions 0644
func WriteTMDBID(tmdbID, mediaPath string) error {
	if tmdbID == "" {
		return fmt.Errorf("TMDB ID cannot be empty")
	}

	// Validate path to prevent directory traversal
	cleanPath := filepath.Clean(mediaPath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid path: directory traversal detected")
	}

	// Ensure the directory exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", cleanPath)
	}

	// Create the tmdb.txt file path
	tmdbPath := filepath.Join(cleanPath, "tmdb.txt")

	// Write the TMDB ID to the file
	err := os.WriteFile(tmdbPath, []byte(tmdbID), 0644)
	if err != nil {
		return fmt.Errorf("failed to write tmdb.txt: %w", err)
	}

	log.Printf("Wrote TMDB ID %s to %s", tmdbID, tmdbPath)
	return nil
}

// ValidateTMDBID verifies that a TMDB ID exists and matches the expected media type
// Returns an error if the ID doesn't exist or the type mismatches
func (c *TMDBClient) ValidateTMDBID(tmdbID string, mediaType MediaType) error {
	if tmdbID == "" {
		return fmt.Errorf("TMDB ID cannot be empty")
	}

	// Attempt to fetch metadata based on media type
	if mediaType == Film {
		_, err := c.FetchMovieMetadata(tmdbID)
		if err != nil {
			return fmt.Errorf("invalid movie ID or API error: %w", err)
		}
		return nil
	} else if mediaType == TV {
		_, err := c.FetchTVMetadata(tmdbID)
		if err != nil {
			return fmt.Errorf("invalid TV show ID or API error: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unknown media type: %v", mediaType)
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

// saveTitle saves the title text to title.txt
func (c *TMDBClient) saveTitle(title, destDir string) error {
	if title == "" {
		return fmt.Errorf("title is empty")
	}

	destPath := filepath.Join(destDir, "title.txt")

	// Write the title to file
	err := os.WriteFile(destPath, []byte(title), 0644)
	if err != nil {
		return fmt.Errorf("failed to write title file: %w", err)
	}

	log.Printf("Saved title to %s", destPath)
	return nil
}

// FetchAndSaveMetadata fetches metadata and downloads poster, description, genres, and title for a media item
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

	titlePath := filepath.Join(media.Path, "title.txt")
	_, titleErr := os.Stat(titlePath)
	titleExists := titleErr == nil

	// If all files exist, skip fetching
	if posterExists && descriptionExists && genreExists && titleExists {
		log.Printf("All metadata files already exist for %s, skipping download", media.Title)
		return nil
	}

	var posterPath string
	var overview string
	var genres []Genre
	var title string
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
		title = movie.Title
	} else if media.Type == TV {
		tv, err := c.FetchTVMetadata(media.TMDBID)
		if err != nil {
			return fmt.Errorf("failed to fetch TV metadata: %w", err)
		}
		posterPath = tv.PosterPath
		overview = tv.Overview
		genres = tv.Genres
		title = tv.Name
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

	// Save title if it doesn't exist
	if !titleExists {
		if title == "" {
			log.Printf("Warning: No title available for %s", media.Title)
		} else {
			if err = c.saveTitle(title, media.Path); err != nil {
				log.Printf("Warning: Failed to save title for %s: %v", media.Title, err)
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
