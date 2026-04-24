---
id: task-020.03
title: Fetch and display track list from MusicBrainz
status: Done
assignee: []
created_date: '2025-11-12 13:28'
updated_date: '2025-11-16 15:11'
labels:
  - musicbrainz
  - music
  - metadata
  - tracks
dependencies: []
parent_task_id: task-020
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement track list fetching from MusicBrainz API and display on Music media detail pages. When a `musicbrainz.txt` file exists, fetch the release's track list (including track titles, positions, and durations) and store it in a `tracks.json` file in the media directory.

The detail page should display the complete track list in an organized, readable format showing disc/side numbers, track numbers, titles, and durations.

**User Value:** Users can see the complete track listing for their music releases directly in the web interface without needing to mount or browse the physical media.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 MusicBrainz client can fetch release track list from API
- [x] #2 Track list includes: disc/medium number, position, title, and duration
- [x] #3 Track list stored in tracks.json in media directory
- [x] #4 Track list fetched during initial scan when musicbrainz.txt exists
- [x] #5 Existing tracks.json files are not overwritten (caching behavior)
- [x] #6 Detail page displays track list in organized format with proper grouping by disc
- [x] #7 Track list only displays for Music media type
- [x] #8 Failed track list fetches log warnings but don't stop scanning
- [x] #9 Tests cover track list fetching, JSON serialization, and display
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Summary

Successfully implemented MusicBrainz track list fetching and display functionality for Music media type.

### Key Components Added:

1. **Data Structures** (models.go):
   - `Track` - Individual track with position, title, and duration (milliseconds)
   - `Medium` - Disc/medium containing tracks with position and format
   - `TrackList` - Complete track listing with all media/discs
   - `Track.FormatDuration()` - Helper method for mm:ss display format
   - `Media.LoadTrackList()` - Loads track list from tracks.json

2. **MusicBrainz Client** (musicbrainz.go):
   - `MusicBrainzClient` - API client for fetching release data
   - `FetchRelease()` - Fetches release with recordings from MusicBrainz API
   - `ConvertToTrackList()` - Converts API response to internal format
   - `SaveTrackList()` - Saves track list to tracks.json with formatting
   - `FetchAndSaveTrackList()` - Complete workflow with caching
   - Proper User-Agent header for MusicBrainz API compliance
   - Fallback from track data to recording data for title and duration

3. **Scanner Integration** (scanner.go):
   - Added `musicBrainzClient` field to Scanner
   - `NewScannerWithClients()` - Constructor supporting both TMDB and MusicBrainz clients
   - Track list fetching in `parseMusic()` when MusicBrainz ID exists
   - Warning logs for failed fetches without stopping scan

4. **Main Application** (main.go):
   - Initialize MusicBrainz client (no API key required)
   - Pass both TMDB and MusicBrainz clients to scanner

5. **Web Interface**:
   - **Template** (templates/detail.html):
     - Track list display section (Music media only)
     - Disc grouping with position and format labels
     - Track table with number, title, and formatted duration
     - Responsive CSS styling
   - **Handler** (handlers.go):
     - Load and pass track list to template in DetailHandler

6. **Tests** (musicbrainz_test.go, scanner_test.go):
   - Track list conversion tests (including fallback behavior)
   - Track list saving and loading tests
   - Caching behavior verification (existing files not overwritten)
   - Duration formatting tests
   - Scanner integration tests

### Caching Behavior:
- Checks for existing tracks.json before fetching
- Skips download if file exists
- Logs "Track list already exists" message
- Consistent with TMDB metadata caching

### API Integration:
- MusicBrainz API endpoint: `/ws/2/release/{mbid}?inc=recordings&fmt=json`
- Required User-Agent header set
- No API key required (public API)
- JSON response parsing with fallback handling

### Display Features:
- Only shown for Music media type (Type == 2)
- Disc grouping when multiple discs present
- Track position, title, and duration display
- Duration formatted as mm:ss (e.g., "3:45")
- Clean table layout with proper styling

### Error Handling:
- Failed API requests log warnings
- Missing data falls back to recording data
- Parse errors don't crash scanner
- Missing tracks.json returns nil (graceful degradation)

All acceptance criteria met and tested successfully.
<!-- SECTION:NOTES:END -->
