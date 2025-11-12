---
id: task-019
title: Add Music media type without TMDB integration
status: To Do
assignee: []
created_date: '2025-11-12 13:21'
labels:
  - feature
  - media-types
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a third media type "Music" to complement existing Film and TV types. Music media should not integrate with TMDB API, but should still support local metadata files (poster, title, description, genre.txt) stored in the media directory.

Music content typically represents concert films, music documentaries, or artist performances on physical media (Blu-ray/DVD).

**User Value:** Users can organize and display music-related physical media in their backup library alongside films and TV shows, with proper categorization and metadata display.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 MediaType enum includes a Music type alongside Film and TV
- [ ] #2 Scanner recognizes directories with [Music] suffix (e.g., 'Artist Name [Music]/')
- [ ] #3 Music media displays in the web UI with appropriate type badge and color
- [ ] #4 Music detail pages show poster, title, description, and genres from local files
- [ ] #5 TMDB search/set functionality is hidden for Music media type
- [ ] #6 Import workflow supports importing Music media type
- [ ] #7 Music media appears in library sorted alphabetically after TV shows
- [ ] #8 Existing tests updated to cover Music media type where applicable
- [ ] #9 New tests added for Music-specific behavior (scanning, display, metadata loading)
<!-- AC:END -->
