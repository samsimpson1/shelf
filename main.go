package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// printHelp prints the help message showing all configuration options
func printHelp(w io.Writer) {
	fmt.Fprintf(w, `Media Backup Manager - A web application for managing media disk backups

Usage:
  ./shelf           Start the server with default or environment-configured settings
  ./shelf -help     Show this help message
  ./shelf --help    Show this help message
  ./shelf -h        Show this help message

Configuration:
  The application is configured using environment variables:

  MEDIA_DIR
      Path to media backup directory
      Default: /home/sam/Scratch/media/backup

  IMPORT_DIR
      Path to import directory for organizing raw disk backups (optional)
      If not set, import functionality will be disabled
      Default: empty

  PORT
      HTTP server port
      Default: 8080

  TMDB_API_KEY
      TMDB API key for metadata fetching (optional)
      If not set, poster and metadata fetching will be disabled
      Get your API key at: https://www.themoviedb.org/settings/api

  DEV_MODE
      Development mode - templates will be reloaded on every request (optional)
      Set to "true" to enable
      Default: false

  PLAY_URL_PREFIX
      URL prefix for VLC play commands (optional)
      Used to construct full paths for network shares or mount points
      Default: empty (assumes local paths)

Examples:
  # Start with defaults
  ./shelf

  # Start with custom media directory and port
  MEDIA_DIR=/path/to/media PORT=9000 ./shelf

  # Start with TMDB metadata fetching enabled
  TMDB_API_KEY=your_api_key_here ./shelf

  # Start in development mode
  DEV_MODE=true ./shelf

  # Start with network path prefix for VLC play commands
  PLAY_URL_PREFIX=/mnt/media ./shelf
`)
}

// shouldShowHelp checks if help flag is present in command-line arguments
func shouldShowHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-help" || arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func main() {
	// Check for help flag
	if shouldShowHelp(os.Args) {
		printHelp(os.Stdout)
		os.Exit(0)
	}

	// Read configuration from environment variables
	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "/home/sam/Scratch/media/backup"
	}

	importDir := os.Getenv("IMPORT_DIR")
	if importDir == "" {
		log.Println("Warning: IMPORT_DIR not set, import functionality will be disabled")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	if tmdbAPIKey == "" {
		log.Println("Warning: TMDB_API_KEY not set, poster fetching will be disabled")
	}

	devMode := os.Getenv("DEV_MODE") == "true"
	if devMode {
		log.Println("Development mode enabled - templates will be reloaded on every request")
	}

	playURLPrefix := os.Getenv("PLAY_URL_PREFIX")
	if playURLPrefix == "" {
		playURLPrefix = "" // Empty by default, assumes local paths
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
	var tmdbClient *TMDBClient
	if tmdbAPIKey != "" {
		log.Println("TMDB API key configured, poster fetching enabled")
		tmdbClient = NewTMDBClient(tmdbAPIKey)
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
	tmpl, err := template.ParseFiles(
		"templates/index.html",
		"templates/detail.html",
		"templates/search.html",
		"templates/confirm.html",
		"templates/import_list.html",
		"templates/import_step1.html",
		"templates/import_step2.html",
		"templates/import_step3.html",
		"templates/import_step4.html",
		"templates/import_step5.html",
		"templates/import_confirm.html",
		"templates/import_success.html",
	)
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Create app
	app := NewApp(mediaList, tmpl, mediaDir, importDir)
	app.SetDevMode(devMode)
	app.SetPlayURLPrefix(playURLPrefix)

	// Set TMDB client if available
	if tmdbClient != nil {
		app.SetTMDBClient(tmdbClient)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.IndexHandler)
	mux.HandleFunc("/posters/", app.PosterHandler)

	// Import routes
	mux.HandleFunc("/import", app.ImportListHandler)
	mux.HandleFunc("/import/start", app.ImportStartHandler)
	mux.HandleFunc("/import/step1", app.ImportStep1Handler)
	mux.HandleFunc("/import/step2", app.ImportStep2Handler)
	mux.HandleFunc("/import/step2/confirm", app.ImportStep2ConfirmHandler)
	mux.HandleFunc("/import/step3", app.ImportStep3Handler)
	mux.HandleFunc("/import/step4", app.ImportStep4Handler)
	mux.HandleFunc("/import/step5", app.ImportStep5Handler)
	mux.HandleFunc("/import/confirm", app.ImportConfirmHandler)
	mux.HandleFunc("/import/execute", app.ImportExecuteHandler)
	mux.HandleFunc("/import/success", app.ImportSuccessHandler)

	// TMDB routes (must come before the general /media/ route)
	mux.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Route to specific handlers based on path suffix
		if strings.HasSuffix(path, "/search-tmdb") {
			app.SearchTMDBHandler(w, r)
		} else if strings.HasSuffix(path, "/confirm-tmdb") {
			app.ConfirmTMDBHandler(w, r)
		} else if strings.HasSuffix(path, "/set-tmdb") {
			app.SaveTMDBHandler(w, r)
		} else {
			// Default to detail handler
			app.DetailHandler(w, r)
		}
	})

	// Serve static files (CSS, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
