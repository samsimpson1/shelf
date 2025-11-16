---
id: task-019
title: Add Music media type without TMDB integration
status: Done
assignee: []
created_date: '2025-11-12 13:21'
updated_date: '2025-11-16 14:43'
labels:
  - feature
  - media-types
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a third media type "Music" to complement existing Film and TV types. Music media should not integrate with TMDB API, but should still support local metadata files (poster, title, description, genre.txt) stored in the media directory.

Music content can include concert films, music documentaries, artist performances, soundtracks, and albums on various physical media formats including Blu-ray (BDMV), DVD, CD images, and other disk types.

**User Value:** Users can organize and display music-related physical media in their backup library alongside films and TV shows, with proper categorization and metadata display across all disk formats.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 MediaType enum includes a Music type alongside Film and TV
- [x] #2 Scanner recognizes directories with [Music] suffix (e.g., 'Artist Name [Music]/')
- [x] #3 Music media supports all disk types: BDMV (Blu-ray), DVD, CD images, and custom formats
- [x] #4 Music media disk counting works correctly across all supported disk types
- [x] #5 Music media displays in the web UI with appropriate type badge and color
- [x] #6 Music detail pages show poster, title, description, and genres from local files
- [x] #7 TMDB search/set functionality is hidden for Music media type
- [x] #8 Import workflow supports importing Music media type with all disk format options
- [x] #9 Music media appears in library sorted alphabetically after TV shows

- [x] #10 Existing tests updated to cover Music media type where applicable
- [x] #11 New tests added for Music-specific behavior (scanning, display, metadata loading, multiple disk types)
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Summary

Successfully implemented Music media type support without TMDB integration.

### Changes Made:

**Models (models.go)**
- Added Music to MediaType enum (Film=0, TV=1, Music=2)
- Updated comments to reflect Music doesn't have years or TMDB integration
- Updated DisplayTitle() to handle Music (no year display)

**Scanner (scanner.go)**
- Added musicPattern regex: `^(.+) \[Music\]$`
- Implemented parseMusic() function (similar to parseTV but no TMDB fetching)
- Music uses collectFilmDisks() for simple "Disk [Format]" pattern
- Updated Scan() to try parsing as Music after Film and TV

**Handlers (handlers.go)**
- Updated sorting to use `sorted[i].Type < sorted[j].Type` (Film < TV < Music)
- Added Music type checks in SearchTMDBHandler, ConfirmTMDBHandler, and SaveTMDBHandler to block TMDB functionality

**Templates**
- index.html: Added Music emoji (🎵) for placeholders
- detail.html: Added Music emoji and hidden TMDB actions for Music type
- import_step1.html: Added Music radio button option

**Import System**
- import.go: Updated GenerateMediaDirName() to generate "Title [Music]" directories
- import.go: Updated GenerateDiskDirName() to treat Music like Film (simple disk pattern)
- import_handlers.go: Added Music case to ImportStep1Handler
- import_handlers.go: Updated ImportStep4Handler to treat Music like Film (disk number defaults to 1)

**Tests**
- models_test.go: Added Music test cases to TestMediaTypeString and TestMediaDisplayTitle
- scanner_test.go: Added TestParseMusicName with comprehensive test cases
- import_test.go: Added Music test cases to TestGenerateMediaDirName and TestGenerateDiskDirName
- Created and ran end-to-end Music test to verify full workflow

### All Acceptance Criteria Met:
✓ MediaType enum includes Music type
✓ Scanner recognizes [Music] suffix directories
✓ Music supports all disk types (BDMV, DVD, CD, custom)
✓ Disk counting works correctly
✓ Web UI displays Music with type badge (🎵) and color
✓ Detail pages show local metadata (poster, title, description, genres)
✓ TMDB functionality explicitly disabled for Music
✓ Import workflow fully supports Music media type
✓ Music sorted alphabetically after TV shows
✓ Existing tests updated for Music compatibility
✓ New Music-specific tests added and passing

### Testing:
- All existing tests pass (2 pre-existing failures unrelated to Music changes)
- New Music-specific tests added and passing
- End-to-end test verified full Music workflow
- Code compiles successfully

### Commit:
Changes committed and pushed to branch: claude/work-on-ta-0114JzyDr6RWhWKK5eK5d8b2
<!-- SECTION:NOTES:END -->
