---
id: task-001.04
title: Implement HTTP Handlers for TMDB Workflow
status: Done
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:59'
labels: []
dependencies:
  - task-001.01
  - task-001.02
  - task-001.03
parent_task_id: task-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create HTTP handlers to connect the UI templates with the TMDB API and file management functionality. This implements the complete user workflow from search to save.

## Technical Scope
- SearchTMDBHandler (GET/POST) for search form and results
- ConfirmTMDBHandler (GET) for confirmation page
- SaveTMDBHandler (POST) for writing tmdb.txt and triggering metadata fetch
- Error handling and user feedback
- Input validation and sanitization

## Dependencies
Requires TMDB Search API, File Management, and Web UI Templates.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 SearchTMDBHandler displays search form with pre-filled media title
- [x] #2 SearchTMDBHandler processes search queries and renders results
- [x] #3 ConfirmTMDBHandler validates TMDB ID and shows preview
- [x] #4 SaveTMDBHandler writes tmdb.txt file to correct directory
- [x] #5 SaveTMDBHandler triggers metadata fetch when checkbox is selected
- [x] #6 SaveTMDBHandler redirects to detail page with success message
- [x] #7 All handlers validate slug to prevent path traversal
- [x] #8 Search queries are sanitized before passing to TMDB API
- [x] #9 Network errors and API errors display user-friendly messages
- [x] #10 Unit tests cover all handlers and error paths
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully:

## Updated App Structure
- Added `tmdbClient *TMDBClient` field to App struct
- Added `SetTMDBClient(client *TMDBClient)` method
- Updated `loadTemplates()` to include search.html and confirm.html

## SearchTMDBHandler (handlers.go:190-291)
- Validates TMDB client is configured (returns 503 if not)
- Extracts slug from URL `/media/{slug}/search-tmdb`
- Validates slug and finds media by slug (security check)
- Pre-fills query field with media title on initial GET
- Parses query and year parameters from form
- Calls appropriate search API (SearchMovies or SearchTV) based on media type
- Uses media year as default for film searches
- Sanitizes all inputs through URL query parsing
- Displays user-friendly error messages on API failures
- Renders search.html template with results

## ConfirmTMDBHandler (handlers.go:293-401)
- Validates TMDB client is configured
- Extracts slug and validates media
- Gets TMDB ID from query parameter
- Fetches full metadata from TMDB API to validate ID
- Converts MovieResponse/TVResponse to SearchResult format for template
- Loads current media description and poster status
- Displays error message if fetch fails
- Renders confirm.html with side-by-side comparison

## SaveTMDBHandler (handlers.go:403-480)
- Only accepts POST requests (405 for other methods)
- Validates TMDB client is configured
- Extracts slug and validates media
- Parses form data securely
- Validates TMDB ID matches media type via API call
- Writes tmdb.txt file using WriteTMDBID() with security checks
- Updates in-memory media object with new TMDB ID
- Checks "download_metadata" checkbox value
- Triggers FetchAndSaveMetadata if checkbox selected
- Logs warnings for metadata fetch failures (non-fatal)
- Redirects to detail page with URL-escaped slug

## Security & Error Handling
- All handlers validate slug through findMediaBySlug
- Path traversal prevented by existing slug validation
- Query parameters sanitized through Go's URL parsing
- Form data parsed with r.ParseForm()
- TMDB API errors handled gracefully with user messages
- Network errors display friendly messages
- All handlers check if TMDB client is nil

All 10 acceptance criteria met. Comprehensive error handling and security validation implemented.
<!-- SECTION:NOTES:END -->
