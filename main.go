package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read configuration from environment variables
	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "/home/sam/Scratch/media/backup"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	if tmdbAPIKey == "" {
		log.Println("Warning: TMDB_API_KEY not set, poster fetching will be disabled")
	}

	// Validate media directory exists
	info, err := os.Stat(mediaDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Media directory does not exist: %s", mediaDir)
		}
		log.Fatalf("Cannot access media directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("Media path is not a directory: %s", mediaDir)
	}

	// Create scanner with optional TMDB client
	var scanner *Scanner
	if tmdbAPIKey != "" {
		log.Println("TMDB API key configured, poster fetching enabled")
		tmdbClient := NewTMDBClient(tmdbAPIKey)
		scanner = NewScannerWithTMDB(mediaDir, tmdbClient)
	} else {
		scanner = NewScanner(mediaDir)
	}

	// Scan media directory
	log.Printf("Scanning media directory: %s", mediaDir)
	mediaList, err := scanner.Scan()
	if err != nil {
		log.Fatalf("Failed to scan media directory: %v", err)
	}
	log.Printf("Found %d media items", len(mediaList))

	// Load templates
	tmpl, err := template.ParseFiles("templates/index.html", "templates/detail.html")
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Create app
	app := NewApp(mediaList, tmpl, mediaDir)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.IndexHandler)
	mux.HandleFunc("/media/", app.DetailHandler)
	mux.HandleFunc("/posters/", app.PosterHandler)

	// Serve static files (CSS, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
