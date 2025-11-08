# Media Backup Management Web App - Initial Implementation Plan

## Project Overview
A simple Go web application to display and manage media disk backups (BDMV, DVD, etc.) from `/home/sam/Scratch/media/backup/`.

## Directory Structure Analysis
Based on exploration, the media follows this structure:
- **Films**: `Title (Year) [Film]/` with single `Disk [Format]/` subdirectories
- **TV Shows**: `Title [TV]/` with multiple `Series X Disk Y [Format]/` subdirectories
- Each media item may contain a `tmdb.txt` file with a TMDB ID (single line)

## Implementation Plan

### 1. Project Structure
```
/home/sam/Documents/git/shelf/
├── main.go              # Entry point, HTTP server setup
├── models.go            # Data structures (Media, Film, TVShow)
├── models_test.go       # Unit tests for models
├── scanner.go           # Directory scanning and parsing logic
├── scanner_test.go      # Unit tests for scanner
├── handlers.go          # HTTP handlers for web UI
├── handlers_test.go     # Unit tests for handlers
├── integration_test.go  # Integration tests
├── templates/
│   └── index.html       # Simple web UI template
├── testdata/            # Test fixtures for unit tests
│   ├── Film 1 (2020) [Film]/
│   │   ├── tmdb.txt
│   │   └── Disk [Blu-Ray]/
│   ├── TV Show [TV]/
│   │   ├── tmdb.txt
│   │   ├── Series 1 Disk 1 [Blu-Ray]/
│   │   └── Series 1 Disk 2 [Blu-Ray]/
│   └── No TMDB (2021) [Film]/
│       └── Disk [DVD]/
├── go.mod               # Go module file
└── PLAN.md              # This file
```

### 2. Core Data Structures (`models.go`)
- `MediaType` enum (Film, TV)
- `Media` struct with:
  - Title (string)
  - Type (MediaType)
  - Year (int, optional for TV)
  - DiskCount (int)
  - TMDBID (string, optional)
  - Path (string)

### 3. Scanner Logic (`scanner.go`)
- Accept directory path as parameter (configurable)
- Parse folder names using regex patterns:
  - Film: `^(.+) \((\d{4})\) \[Film\]$`
  - TV: `^(.+) \[TV\]$`
- Count disks by matching subdirectories:
  - Film: `Disk [...]`
  - TV: `Series X Disk Y [...]`
- Read `tmdb.txt` if present (trim whitespace from ID)
- Store all data in-memory (slice of Media structs)
- Return errors for invalid paths or permission issues

### 4. Web Server (`main.go`)
- Read configuration from environment variables:
  - `MEDIA_DIR` - Path to media backup directory (default: `/home/sam/Scratch/media/backup`)
  - `PORT` - HTTP server port (default: `8080`)
- Initialize HTTP server on configured port
- Route handlers:
  - `GET /` - Display media list
  - `GET /static/*` - Serve static assets (if needed later)
- On startup, scan configured directory and load media into memory
- Validate that `MEDIA_DIR` exists and is readable before starting

### 5. HTTP Handlers (`handlers.go`)
- `indexHandler`: Render media list using template
- Pass sorted media list to template (Films first, then TV shows)

### 6. Web UI (`templates/index.html`)
Simple HTML template displaying:
- Page title: "Media Backup Manager"
- Table with columns:
  - **Title** (with year for films)
  - **Type** (Film/TV)
  - **Disks** (count)
  - **TMDB ID** (if available)
- Basic CSS for readability (inline or minimal)
- No JavaScript required for initial version

### 7. Testing Strategy

#### Unit Tests

**`models_test.go`**:
- Test `MediaType` string representation
- Test `Media` struct field validation
- Test helper methods (if any)

**`scanner_test.go`**:
- Test film name parsing (with/without year, edge cases)
- Test TV show name parsing
- Test disk counting for films
- Test disk counting for TV shows (multiple series)
- Test TMDB ID reading (valid, missing, malformed)
- Test directory filtering (ignore non-media folders)
- Test error handling (invalid path, permission errors, invalid structure)
- Test scanner with custom directory paths
- Use `testdata/` directory for test fixtures

**`handlers_test.go`**:
- Test HTTP handler responses (status codes)
- Test template rendering with mock data
- Test empty media list handling
- Test handler with various media combinations
- Use `httptest` package for HTTP testing

#### Integration Tests

**`integration_test.go`**:
- Test full scan of `testdata/` directory
- Test complete HTTP request/response cycle
- Test server startup and shutdown
- Test with different `MEDIA_DIR` configurations
- Test environment variable parsing and defaults
- Verify end-to-end data flow: scan → store → display

#### Test Coverage Goals
- Aim for >80% code coverage
- All public functions should have tests
- Edge cases and error paths covered
- Use `go test -cover` to measure coverage
- Use `go test -race` to detect race conditions

### 8. Implementation Steps
1. Initialize Go module (`go mod init shelf`)
2. Create `testdata/` directory with test fixtures
3. Create `models.go` with data structures
4. Create `models_test.go` with unit tests (TDD approach)
5. Create `scanner.go` with directory scanning logic
6. Create `scanner_test.go` with comprehensive tests
7. Create `handlers.go` with HTTP handlers
8. Create `handlers_test.go` with handler tests
9. Create `templates/index.html` with simple UI
10. Create `integration_test.go` for end-to-end tests
11. Create `main.go` to tie everything together
12. Run full test suite and verify coverage
13. Test manually with actual media directory

### 9. Key Technical Decisions
- **No external dependencies**: Use standard library only (`net/http`, `html/template`, `os`, `path/filepath`, `regexp`)
- **In-memory storage**: All data stored in a simple slice
- **Single scan on startup**: No live filesystem watching
- **No state persistence**: Restart to refresh data
- **Read-only**: No modifications to filesystem
- **Environment-based configuration**: Use env vars for `MEDIA_DIR` and `PORT` with sensible defaults
- **Test-Driven Development**: Write tests before/alongside implementation
- **High test coverage**: Target >80% coverage with comprehensive unit and integration tests

### 10. Future Enhancements (Not in Initial Version)
- Refresh button to rescan directory
- Filtering/sorting capabilities
- Link to TMDB pages
- Detailed view per media item
- Disk format information display
- Database persistence

## Success Criteria
- Web server runs on configurable port (default `localhost:8080`)
- Reads media directory from `MEDIA_DIR` environment variable (default `/home/sam/Scratch/media/backup`)
- Lists all 5 media items (3 films, 2 TV shows) from configured directory
- Correctly identifies media types
- Displays disk counts (Films: 1, Better Call Saul: 5, The Good Place: 4)
- Shows TMDB IDs for all items
- Clean, readable UI
- >80% test coverage with passing unit and integration tests
- Proper error handling for invalid/missing directory paths
