---
id: task-019
title: Add Music media type without TMDB integration
status: To Do
assignee: []
created_date: '2025-11-12 13:21'
updated_date: '2025-11-12 13:22'
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
- [ ] #1 MediaType enum includes a Music type alongside Film and TV
- [ ] #2 Scanner recognizes directories with [Music] suffix (e.g., 'Artist Name [Music]/')
- [ ] #3 Music media supports all disk types: BDMV (Blu-ray), DVD, CD images, and custom formats
- [ ] #4 Music media disk counting works correctly across all supported disk types
- [ ] #5 Music media displays in the web UI with appropriate type badge and color
- [ ] #6 Music detail pages show poster, title, description, and genres from local files
- [ ] #7 TMDB search/set functionality is hidden for Music media type
- [ ] #8 Import workflow supports importing Music media type with all disk format options
- [ ] #9 Music media appears in library sorted alphabetically after TV shows

- [ ] #10 Existing tests updated to cover Music media type where applicable
- [ ] #11 New tests added for Music-specific behavior (scanning, display, metadata loading, multiple disk types)
<!-- AC:END -->
