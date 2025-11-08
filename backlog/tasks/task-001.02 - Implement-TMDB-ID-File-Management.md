---
id: task-001.02
title: Implement TMDB ID File Management
status: Done
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:50'
labels: []
dependencies:
  - task-001.01
parent_task_id: task-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create functions to write, validate, and manage tmdb.txt files in media directories. This enables persistence of TMDB ID selections and ensures file operations are secure.

## Technical Scope
- WriteTMDBID function to create/overwrite tmdb.txt files
- ValidateTMDBID function to verify IDs exist and match media type
- File permission handling (0644)
- Path validation to prevent directory traversal

## Dependencies
Depends on TMDB Search API for validation functionality.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 WriteTMDBID creates tmdb.txt file with correct permissions (0644)
- [x] #2 WriteTMDBID overwrites existing tmdb.txt files safely
- [x] #3 ValidateTMDBID verifies TMDB ID exists via API
- [x] #4 ValidateTMDBID prevents type mismatches (film ID for TV show, etc.)
- [x] #5 Path validation prevents directory traversal attacks
- [x] #6 File write errors are handled gracefully with clear messages
- [x] #7 Unit tests cover file creation, validation, and error paths
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully:

- Implemented `WriteTMDBID(tmdbID, mediaPath)` function that:
  - Creates tmdb.txt file with 0644 permissions
  - Safely overwrites existing tmdb.txt files
  - Validates input (rejects empty TMDB IDs)
  - Prevents directory traversal attacks with path validation
  - Checks directory exists before writing
  - Returns clear error messages on failures
  
- Implemented `ValidateTMDBID(tmdbID, mediaType)` method on TMDBClient that:
  - Verifies TMDB ID exists via API
  - Prevents type mismatches (validates movie IDs as movies, TV IDs as TV shows)
  - Uses existing FetchMovieMetadata/FetchTVMetadata for validation
  - Returns descriptive error messages
  
- Comprehensive unit tests cover:
  - File creation with correct permissions
  - Overwriting existing files
  - Empty ID validation
  - Directory traversal prevention
  - Non-existent directory handling
  - Write permission errors
  - Invalid movie/TV ID validation
  - Type mismatch detection
  
All 7 acceptance criteria met. Tests pass successfully.
<!-- SECTION:NOTES:END -->
