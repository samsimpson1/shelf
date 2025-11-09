# Media Backup Manager - Claude Code Documentation

This project was built by Claude Code following a test-driven development approach.

## Project Overview

A simple Go web application that scans and displays media disk backups (BDMV, DVD, etc.) from a configured directory. The application parses directory structures to identify films and TV shows, counts their disks, and displays them in a clean web interface.

## Architecture

### Core Components

1. **Models** ([models.go](models.go))
   - `MediaType` enum for Film/TV classification
   - `Media` struct containing title, type, year, disk count, TMDB ID, and path
   - Helper methods for display formatting
   - `Slug()` - Generates URL-friendly slugs
   - `LoadDescription()` / `LoadGenres()` - On-demand metadata loading
   - `FindPosterFile()` / `PosterURL()` - Poster file management

2. **Scanner** ([scanner.go](scanner.go))
   - Directory scanning and parsing logic
   - Regex-based pattern matching for film and TV show directories
   - Disk counting for both media types
   - TMDB ID reading from `tmdb.txt` files
   - Optional TMDB poster fetching during scan

3. **TMDB Client** ([tmdb.go](tmdb.go))
   - TMDB API integration for fetching metadata
   - Movie and TV show metadata retrieval (overview, genres, poster, title)
   - Poster image downloading and saving
   - Description (overview), genre list, and title saving
   - Automatic detection of existing metadata files
   - Smart caching to avoid re-downloading

4. **Handlers** ([handlers.go](handlers.go))
   - HTTP request handlers
   - Media list sorting (films first, then TV, alphabetically within each type)
   - Template rendering
   - Poster image serving with security validation
   - Detail page routing with slug-based URLs

5. **Import System** ([import.go](import.go), [import_handlers.go](import_handlers.go))
   - Import directory scanning for raw disk backups
   - Multi-step import workflow with web UI
   - TMDB search integration for automatic metadata
   - Manual title/year entry option
   - Disk type detection (Blu-ray, DVD, custom)
   - Name sanitization for filesystem compatibility
   - Add to existing media or create new entries
   - File moving and validation

6. **Main** ([main.go](main.go))
   - Application entry point
   - Environment-based configuration
   - HTTP server setup
   - Optional TMDB client initialization
   - Optional import scanner initialization

### Directory Structure Parsing

The application expects media to follow this naming convention:

**Films:**
```
Title (Year) [Film]/
  ├── Disk [Format]/
  ├── tmdb.txt (optional - contains TMDB ID)
  ├── poster.jpg (optional - auto-downloaded from TMDB)
  ├── description.txt (optional - auto-downloaded from TMDB)
  ├── genre.txt (optional - auto-downloaded from TMDB)
  └── title.txt (optional - auto-downloaded from TMDB)
```

**TV Shows:**
```
Title [TV]/
  ├── Series X Disk Y [Format]/
  ├── Series X Disk Z [Format]/
  ├── tmdb.txt (optional - contains TMDB ID)
  ├── poster.jpg (optional - auto-downloaded from TMDB)
  ├── description.txt (optional - auto-downloaded from TMDB)
  ├── genre.txt (optional - auto-downloaded from TMDB)
  └── title.txt (optional - auto-downloaded from TMDB)
```

### Regex Patterns

- Film: `^(.+) \((\d{4})\) \[Film\]$`
- TV: `^(.+) \[TV\]$`
- Film Disk: `^Disk \[.+\]$`
- TV Disk: `^Series (\d+) Disk (\d+) \[.+\]$`

## Configuration

The application uses environment variables for configuration:

- `MEDIA_DIR` - Path to media backup directory (default: `/home/sam/Scratch/media/backup`)
- `IMPORT_DIR` - Path to import directory for raw disk backups (optional, import functionality disabled if not set)
- `PORT` - HTTP server port (default: `8080`)
- `TMDB_API_KEY` - TMDB API key for metadata fetching (optional, metadata fetching disabled if not set)

### Getting a TMDB API Key

To enable poster fetching, you need a TMDB API key:

1. Create a free account at [themoviedb.org](https://www.themoviedb.org/)
2. Go to Settings → API
3. Request an API key (select "Developer" option)
4. Copy your API Key (v3 auth)
5. Set the environment variable: `export TMDB_API_KEY=your_api_key_here`

### TMDB Metadata Fetching

When `TMDB_API_KEY` is configured:

- Metadata is automatically fetched during the initial directory scan
- Only media items with a `tmdb.txt` file will have metadata fetched
- Four files are saved to each media's directory:
  - `poster.jpg` (or .png, .webp) - Movie/TV poster in "original" size
  - `description.txt` - Overview/synopsis of the movie or TV show
  - `genre.txt` - Comma-separated list of genres (e.g., "Action, Drama, Thriller")
  - `title.txt` - Official title from TMDB (movie.title or tv.name)
- If all four files already exist in the directory, they will not be re-downloaded
- Missing files are fetched individually (e.g., if only poster exists, description, genres, and title will be fetched)
- Failed metadata fetches log warnings but don't stop the scan
- Single API call fetches all metadata efficiently

## Testing

The project achieves 54.7% code coverage with comprehensive tests:

### Unit Tests

- **[models_test.go](models_test.go)** - Tests for MediaType string representation and Media display methods
- **[scanner_test.go](scanner_test.go)** - Tests for directory parsing, disk counting, TMDB ID reading, scanner with TMDB client, and error handling
- **[handlers_test.go](handlers_test.go)** - Tests for HTTP handlers, template rendering, and sorting logic
- **[tmdb_test.go](tmdb_test.go)** - Tests for TMDB client, metadata fetching (movies/TV), poster/description/genre/title saving, and file caching
- **[import_test.go](import_test.go)** - Tests for import scanner, disk type detection, name sanitization, directory name generation, and import execution

### Integration Tests

- **[integration_test.go](integration_test.go)** - End-to-end tests covering full scan-and-serve workflows
- **[tmdb_integration_test.go](tmdb_integration_test.go)** - Real TMDB API integration tests (requires TMDB_API_KEY)

### Test Fixtures

- **[test_helpers.go](test_helpers.go)** - Test utility functions for creating isolated test fixtures
  - `setupTestData(t)` - Creates temporary test directory with standard media structure

Test fixtures are created programmatically in each test using `t.TempDir()` for:
  - War of the Worlds (2025) [Film] - 1 disk, TMDB ID 755898
  - Better Call Saul [TV] - 2 disks (Series 1), TMDB ID 60059
  - No TMDB (2021) [Film] - 1 disk, no TMDB ID

This ensures test isolation and consistent behavior across all environments.

### Running Tests

```bash
# Run all tests with coverage
go test -v -cover ./...

# Run with race detection
go test -race ./...

# Skip integration tests (fast mode)
go test -short -v ./...

# Run TMDB integration tests with real API
TMDB_API_KEY=your_api_key_here go test -v -run TestIntegration

# Run all tests including TMDB integration tests
TMDB_API_KEY=your_api_key_here go test -v ./...
```

### TMDB Integration Tests

The [tmdb_integration_test.go](tmdb_integration_test.go) file contains integration tests that make real API calls to TMDB:

**Test Coverage:**
- Movie metadata fetching (Fight Club - ID 550)
- TV metadata fetching (Better Call Saul - ID 60059)
- Movie search with and without year filtering
- TV show search
- Poster download and validation
- Description, genre, and title saving
- Full metadata workflow (all four files)
- TMDB ID validation
- Error handling with invalid IDs
- File caching behavior (skipping existing metadata)
- Result limiting (max 20 results)

**Running TMDB Integration Tests:**

These tests require a valid TMDB API key and are automatically skipped in two cases:
1. When running with `go test -short` (short mode)
2. When the `TMDB_API_KEY` environment variable is not set

To run these tests:

```bash
# Set your TMDB API key (get one at https://www.themoviedb.org/settings/api)
export TMDB_API_KEY=your_api_key_here

# Run all TMDB integration tests
go test -v -run TestIntegration

# Run a specific TMDB integration test
go test -v -run TestIntegrationFetchMovieMetadata
```

**Note:** These tests make real HTTP requests to the TMDB API and will count against your API rate limits. They are designed to be comprehensive but respectful of API usage.

### End-to-End Tests

The project uses [Playwright](https://playwright.dev/) for comprehensive end-to-end testing of the web interface. These tests verify complete user workflows by running a real browser and interacting with the application as a user would.

**Test Coverage:**
- Media details page viewing (posters, descriptions, genres, metadata)
- Copying VLC and MPV play commands from disk lists
- Setting TMDB IDs for media without one
- Changing TMDB IDs for media that already have one
- Importing films with Blu-Ray disk type
- Importing films with DVD disk type
- Importing TV shows with Blu-Ray disk type
- Importing TV shows with DVD disk type
- Importing media with custom/other disk types

**Test Files:**
- **[e2e/media-details.spec.ts](e2e/media-details.spec.ts)** - Media detail page viewing tests
- **[e2e/copy-play-command.spec.ts](e2e/copy-play-command.spec.ts)** - Copy command functionality tests
- **[e2e/tmdb-management.spec.ts](e2e/tmdb-management.spec.ts)** - TMDB ID management tests
- **[e2e/import-workflow.spec.ts](e2e/import-workflow.spec.ts)** - Import workflow tests

**Running E2E Tests:**

Prerequisites:
- Node.js 20 or later
- npm

```bash
# Install dependencies
npm install

# Install Playwright browsers
npx playwright install

# Build the Go application
go build -o shelf .

# Run all E2E tests
npm test

# Run tests in headed mode (see browser)
npm run test:headed

# Run tests in UI mode (interactive debugging)
npm run test:ui

# Run tests in debug mode
npm run test:debug

# View test report
npm run test:report
```

**Test Fixtures:**

Test fixtures are automatically generated before tests run - they are NOT stored in the repository:

- **[e2e/setup-fixtures.ts](e2e/setup-fixtures.ts)** - Generates all test data (disk files, metadata, posters)
- Fixtures are generated as part of the `webServer` command in [playwright.config.ts](playwright.config.ts)
- Fixtures include 3 media items (The Matrix, Breaking Bad, No TMDB Film) and 3 import directories
- See [e2e/README.md](e2e/README.md) for details on how fixtures work

**Adding E2E Tests for New Features:**

When adding a new feature to the web interface, you **must** add corresponding E2E tests:

1. Identify the user workflows your feature enables
2. Create a new test file in `e2e/` (e.g., `e2e/my-feature.spec.ts`)
3. Write tests covering happy paths and edge cases
4. Add any necessary test fixtures by editing [e2e/setup-fixtures.ts](e2e/setup-fixtures.ts)
5. Run tests locally to verify they pass
6. Ensure tests pass in CI before merging

See [E2E_TESTING.md](E2E_TESTING.md) for detailed documentation on writing and debugging E2E tests.

**CI Integration:**

E2E tests run automatically in GitHub Actions on every push and pull request via the `.github/workflows/e2e-tests.yml` workflow. Fixtures are generated automatically in CI. Test reports are uploaded as artifacts for failed runs.

## Building and Running

### Build

```bash
go build -o shelf .
```

### Run

```bash
# With defaults (no poster fetching)
./shelf

# With TMDB poster fetching
TMDB_API_KEY=your_api_key_here ./shelf

# With all custom configuration
MEDIA_DIR=/path/to/media PORT=9000 TMDB_API_KEY=your_api_key_here ./shelf
```

### Access

Open http://localhost:8080 in your browser to view the media backup manager.

## Web Interface

The web UI features a modern poster-based design:

### Index Page ([templates/index.html](templates/index.html))
- **Poster Grid Layout** - Responsive CSS Grid (4 columns → 2 columns on mobile)
- Clickable poster cards with hover effects
- Each card displays:
  - Poster image (2:3 aspect ratio) with fallback emoji placeholders
  - Title with year (for films)
  - Type badge (color-coded: blue for films, purple for TV)
  - Disk count
- Automatic sorting (films first alphabetically, then TV shows alphabetically)
- Empty state handling when no media is found
- Fully responsive design

### Detail Page ([templates/detail.html](templates/detail.html))
- **Two-column layout** (poster + metadata)
- Large poster display (sticky on desktop)
- Full description/overview from TMDB
- Genre tags displayed as chips
- Metadata grid showing: Type, Year, Disks, TMDB ID (with link to themoviedb.org)
- "Back to Library" navigation
- Responsive layout (stacks on mobile)

### Routing
- `/` - Poster grid view (index)
- `/media/{slug}` - Individual media detail page (e.g., `/media/the-thing-1982`)
- `/posters/{slug}` - Poster image serving (e.g., `/posters/the-thing-1982`)
- `/media/{slug}/search-tmdb` - TMDB ID search page
- `/media/{slug}/confirm-tmdb` - TMDB match confirmation page
- `/media/{slug}/set-tmdb` - TMDB ID save endpoint (POST only)
- `/import` - Import media workflow (see Import Media section for full routing)

## TMDB ID Management

The application provides a complete workflow for managing TMDB IDs through the web interface. This feature requires the `TMDB_API_KEY` environment variable to be set.

### Setting a TMDB ID for the First Time

When a media item doesn't have a TMDB ID:

1. **Navigate to the media detail page** - Click on any media item from the library view
2. **Click "Search for TMDB ID"** - A blue button appears when no TMDB ID is set
3. **Search TMDB** - The search form is pre-filled with the media title:
   - For films: You can optionally specify a year to narrow results
   - For TV shows: Search by title only
   - Click "Search TMDB" to find matches
4. **Review search results** - Each result shows:
   - Poster thumbnail
   - Title and release/air date
   - Overview/synopsis
   - Popularity score
5. **Select the correct match** - Click "Select This" on the matching result
6. **Confirm the match** - Review the comparison page showing:
   - Your current media metadata (left side)
   - The TMDB match details (right side)
   - Option to download metadata immediately (checked by default)
7. **Save the TMDB ID** - Click "Confirm and Save TMDB ID"

After saving, you'll be redirected back to the detail page with the new TMDB ID set.

### Manual TMDB ID Entry

If you already know the TMDB ID (e.g., from themoviedb.org):

1. **Navigate to the search page** - Click "Search for TMDB ID" from the detail page
2. **Click "Or enter TMDB ID manually"** - Expands a manual entry form
3. **Enter the numeric TMDB ID** - Type or paste the ID (numbers only)
4. **Click "Set TMDB ID"** - The ID is saved immediately without confirmation

**Finding TMDB IDs manually:**
- Visit [themoviedb.org](https://www.themoviedb.org/)
- Search for your movie or TV show
- The ID is in the URL: `https://www.themoviedb.org/movie/550` → ID is `550`

### Changing an Existing TMDB ID

If a media item already has a TMDB ID but it's incorrect:

1. **Navigate to the media detail page** - View the media with the incorrect ID
2. **Click "Change TMDB ID"** - A gray button appears when a TMDB ID exists
3. **Warning displayed** - You'll see: "⚠️ Changing the TMDB ID will replace existing metadata"
4. **Follow the same search/confirm workflow** - As described above
5. **Review the replacement warning** - The confirmation page shows:
   - Yellow warning box explaining metadata will be replaced
   - Current TMDB ID displayed
   - Side-by-side comparison of current vs. new metadata
6. **Save the new ID** - Existing metadata files are replaced when you confirm

### Metadata Download Options

When confirming a TMDB ID, you have two options for metadata:

**Immediate Download (Default - Recommended)**
- Checkbox: "Download metadata now (poster, description, genres)" - **Checked**
- Metadata is fetched and saved immediately after setting the ID
- You'll see updated poster, description, and genres right away
- No server restart required

**Deferred Download**
- Checkbox: "Download metadata now (poster, description, genres)" - **Unchecked**
- Only the `tmdb.txt` file is created with the ID
- Metadata will be downloaded on the next server restart
- Useful if you want to batch-set IDs and download metadata later

### Metadata Files

When a TMDB ID is set and metadata is downloaded, four files are saved to the media directory:

1. **`tmdb.txt`** - Contains the numeric TMDB ID
2. **`poster.jpg` (or `.png`, `.webp`)** - Movie/TV poster in original quality
3. **`description.txt`** - Overview/synopsis of the movie or TV show
4. **`genre.txt`** - Comma-separated list of genres (e.g., "Action, Drama, Thriller")

These files are stored directly in the media directory alongside the disk folders.

### Refreshing Metadata After ID Change

After changing a TMDB ID, metadata is automatically refreshed if you chose "Download metadata now":

**What gets updated:**
- Poster image (replaces old poster file)
- Description text (replaces old description)
- Genre list (replaces old genres)
- TMDB ID reference (updates `tmdb.txt`)

**What stays the same:**
- Media directory name
- Disk folders and content
- File organization

**Manual refresh (if needed):**
- Delete the metadata files (`poster.*`, `description.txt`, `genre.txt`) from the media directory
- Restart the server - metadata will be re-downloaded from TMDB
- OR use the web interface to change the TMDB ID again with "Download metadata now" checked

### Error Handling

The TMDB ID workflow includes comprehensive error handling:

**When TMDB API is not configured:**
- Setting/changing TMDB IDs is disabled
- "Search for TMDB ID" button is not displayed
- Direct access to search/confirm pages shows: "TMDB API is not configured"

**When search returns no results:**
- "No Results Found" message with suggestion to try different search terms
- Manual entry option remains available

**When TMDB ID validation fails:**
- Invalid IDs (non-existent or wrong media type) are rejected
- Error message displayed: "Invalid TMDB ID: [reason]"
- User can go back and search again

**When metadata download fails:**
- TMDB ID is still saved successfully
- Warning logged to server console
- User can trigger re-download by changing the ID again or restarting the server

### Security Considerations

The TMDB ID management feature includes security measures:

- **Path validation** - All file operations validate paths are within the media directory
- **Input validation** - TMDB IDs are validated against the TMDB API before saving
- **Type checking** - Movie IDs cannot be set for TV shows and vice versa
- **POST-only saves** - TMDB ID changes require POST requests (not GET)
- **No directory traversal** - File paths are sanitized to prevent accessing files outside the media directory

## Import Media

The import functionality streamlines the process of organizing raw disk backups (from tools like MakeMKV) into the proper MEDIA_DIR structure. This feature requires the `IMPORT_DIR` environment variable to be set.

### Setting Up Import

1. **Configure IMPORT_DIR** - Set the environment variable to point to your raw disk backup directory:
   ```bash
   export IMPORT_DIR=/path/to/raw/backups
   ```

2. **Access Import UI** - When IMPORT_DIR is configured, an "Import Media" button appears on the main library page

### Import Workflow

The import process is a guided 5-step workflow:

**Step 1: Select Directory**
- View all directories in IMPORT_DIR
- See directory sizes to identify which backup to import
- Click "Import" to begin

**Step 2: Choose Media Type**
- Select whether the disk is a Film or TV Show
- Disk type is auto-detected when possible (Blu-ray BDMV, DVD VIDEO_TS)

**Step 3: Search TMDB or Manual Entry**
- **Option A: TMDB Search** (if TMDB_API_KEY is configured)
  - Search for the media by title
  - For films, optionally filter by year
  - Select the correct match from search results
  - Title, year, and metadata are fetched automatically
- **Option B: Manual Entry**
  - Enter title manually
  - For films, enter the release year
  - Skip TMDB integration entirely

**Step 4: Disk Details**
- **For TV Shows:**
  - Enter series (season) number
  - Enter disk number within that series
- **For Films:**
  - Disk number defaults to 1
- **Disk Type:**
  - Auto-detected types: Blu-Ray, DVD
  - Manual options: Blu-Ray, Blu-Ray UHD, DVD, or custom text
  - Custom types allow any format description (e.g., "4K HDR", "BDMV")

**Step 5: Add to Existing or Create New**
- **Create New Media Entry** - Creates a new media directory
- **Add to Existing Media** - Adds the disk to an existing media item (for multi-disk releases)
  - Shows compatible existing media of the same type
  - Useful for adding additional disks to a film or TV series

**Step 6: Confirm and Execute**
- Review all import details
- Preview the destination path
- Confirm to execute the import
- Source directory is moved to the destination with proper naming

### Import Features

**Automatic Disk Type Detection:**
- Detects Blu-ray disks (BDMV directory structure)
- Detects DVD disks (VIDEO_TS directory structure)
- Allows manual override for any disk type

**Name Sanitization:**
- Problematic characters are automatically sanitized for filesystem compatibility:
  - Colons (`:`) → underscores (`_`)
  - Slashes (`/`, `\`) → underscores
  - Quotes (`"`) → apostrophes (`'`)
  - Other invalid characters (`<`, `>`, `|`, `?`, `*`) → underscores
  - Multiple consecutive underscores are collapsed to single underscores
- Maintains readability while ensuring cross-platform compatibility

**TMDB Integration:**
- When a TMDB match is selected, the official title and year are used
- TMDB ID is saved to `tmdb.txt` in the media directory
- Metadata (poster, description, genres) can be downloaded immediately
- Or defer metadata download until next server restart

**Directory Organization:**
- Films: `Title (Year) [Film]/Disk [Format]/`
- TV Shows: `Title [TV]/Series X Disk Y [Format]/`
- Format is preserved from disk type selection

**Validation:**
- Prevents importing to paths that already exist
- Ensures IMPORT_DIR and MEDIA_DIR are accessible
- Validates all user input before moving files

### Import Routing

- `/import` - List of directories available for import
- `/import/start` - Start import workflow for selected directory
- `/import/step1` - Choose media type (Film/TV)
- `/import/step2` - TMDB search or skip
- `/import/step3` - Manual title/year entry
- `/import/step4` - Disk details and type selection
- `/import/step5` - Add to existing or create new
- `/import/confirm` - Final confirmation with preview
- `/import/execute` - Execute the import (POST only)
- `/import/success` - Success confirmation page

### Import Security

The import feature includes security measures:

- **Path validation** - All file operations validate paths are within IMPORT_DIR and MEDIA_DIR
- **No directory traversal** - Sanitized paths prevent accessing files outside allowed directories
- **POST-only execution** - Import execution requires POST requests
- **Session isolation** - Each import workflow uses a unique session ID
- **Input validation** - All user input is validated before file operations

### Example Import Scenarios

**Scenario 1: New Film with TMDB**
1. Place raw Blu-ray backup in IMPORT_DIR (e.g., "Matrix_BDMV")
2. Navigate to Import page
3. Select "Matrix_BDMV" and click Import
4. Choose "Film"
5. Search TMDB for "The Matrix", select 1999 release
6. Confirm Blu-Ray disk type (auto-detected)
7. Choose "Create new media entry"
8. Confirm import

Result: `The Matrix (1999) [Film]/Disk [Blu-Ray]/` with TMDB metadata

**Scenario 2: Adding TV Show Disk to Existing Series**
1. Place raw DVD backup in IMPORT_DIR
2. Navigate to Import page and click Import
3. Choose "TV Show"
4. Skip TMDB or search for the show
5. Enter Series 2, Disk 3
6. Confirm DVD disk type (auto-detected)
7. Choose "Add to existing media" and select your TV show
8. Confirm import

Result: Disk added to existing TV show directory as `Series 2 Disk 3 [DVD]/`

## Technical Decisions

1. **Standard Library Only** - No external dependencies; uses only Go's standard library (including for TMDB API calls)
2. **In-Memory Storage** - All data stored in a simple slice; no database required
3. **Single Scan on Startup** - Scans directory once when the server starts, fetching metadata if configured
4. **Poster Storage in Media Directories** - Posters saved directly to each media's directory for easy management
5. **Optional TMDB Integration** - Metadata fetching is optional and doesn't affect core functionality
6. **Smart Metadata Caching** - Existing metadata files (poster, description, genres) are never re-downloaded
7. **Slug-Based Routing** - Clean, SEO-friendly URLs for detail pages (e.g., `/media/the-thing-1982`)
8. **Lazy Loading Metadata** - Description and genres loaded on-demand for detail pages
9. **Security-First File Serving** - Path validation prevents directory traversal attacks
10. **Environment-Based Config** - Configuration via environment variables with sensible defaults
11. **Test-Driven Development** - Tests written alongside implementation

## Code Quality

- All public functions have tests
- Edge cases and error paths covered
- No race conditions detected
- Clear separation of concerns (models, scanning, handlers, main)
- Comprehensive error handling

## Future Enhancements

Potential improvements not included in the initial version:

- Refresh button to rescan directory without restart
- Filtering and sorting controls in the UI
- Direct links to TMDB pages
- Detailed view per media item showing disk formats
- Disk format information in the main table
- Database persistence for faster startup
- Watch mode for live filesystem monitoring
- Search functionality
- Export to CSV/JSON

## Development Process

This project was built following the implementation plan in [PLAN.md](PLAN.md), which outlined:

1. Project structure and file organization
2. Data structures and models
3. Scanner logic with regex patterns
4. HTTP handlers and routing
5. Web UI template design
6. Comprehensive testing strategy
7. Integration and deployment steps

The implementation followed a test-driven approach, with unit tests written before or alongside the implementation code, ensuring high code quality and coverage from the start.

<!-- BACKLOG.MD MCP GUIDELINES START -->

<CRITICAL_INSTRUCTION>

## BACKLOG WORKFLOW INSTRUCTIONS

This project uses Backlog.md MCP for all task and project management activities.

**CRITICAL GUIDANCE**

- If your client supports MCP resources, read `backlog://workflow/overview` to understand when and how to use Backlog for this project.
- If your client only supports tools or the above request fails, call `backlog.get_workflow_overview()` tool to load the tool-oriented overview (it lists the matching guide tools).

- **First time working here?** Read the overview resource IMMEDIATELY to learn the workflow
- **Already familiar?** You should have the overview cached ("## Backlog.md Overview (MCP)")
- **When to read it**: BEFORE creating tasks, or when you're unsure whether to track work

These guides cover:
- Decision framework for when to create tasks
- Search-first workflow to avoid duplicates
- Links to detailed guides for task creation, execution, and completion
- MCP tools reference

You MUST read the overview resource to understand the complete workflow. The information is NOT summarized here.

</CRITICAL_INSTRUCTION>

<!-- BACKLOG.MD MCP GUIDELINES END -->
