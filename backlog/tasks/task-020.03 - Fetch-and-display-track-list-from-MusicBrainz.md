---
id: task-020.03
title: Fetch and display track list from MusicBrainz
status: To Do
assignee: []
created_date: '2025-11-12 13:28'
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
- [ ] #1 MusicBrainz client can fetch release track list from API
- [ ] #2 Track list includes: disc/medium number, position, title, and duration
- [ ] #3 Track list stored in tracks.json in media directory
- [ ] #4 Track list fetched during initial scan when musicbrainz.txt exists
- [ ] #5 Existing tracks.json files are not overwritten (caching behavior)
- [ ] #6 Detail page displays track list in organized format with proper grouping by disc
- [ ] #7 Track list only displays for Music media type
- [ ] #8 Failed track list fetches log warnings but don't stop scanning
- [ ] #9 Tests cover track list fetching, JSON serialization, and display
<!-- AC:END -->
