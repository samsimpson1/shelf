---
id: task-020.01
title: Read musicbrainz.txt file and link to MusicBrainz release
status: To Do
assignee: []
created_date: '2025-11-12 13:28'
labels:
  - musicbrainz
  - music
  - metadata
dependencies: []
parent_task_id: task-020
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement support for reading a `musicbrainz.txt` file from Music media directories. The file should contain a MusicBrainz release ID (UUID format). When present, display a link to the MusicBrainz release page on the media detail page.

Similar to how `tmdb.txt` works for Films/TV, this provides a simple way to associate a Music media item with its MusicBrainz release entry.

**User Value:** Users can quickly navigate from their local media to the official MusicBrainz release page for additional information.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Scanner reads musicbrainz.txt file from Music media directories
- [ ] #2 Media struct stores MusicBrainz release ID
- [ ] #3 Detail page displays MusicBrainz ID in metadata section
- [ ] #4 MusicBrainz ID links to musicbrainz.org release page (format: https://musicbrainz.org/release/{id})
- [ ] #5 MusicBrainz link only appears for Music media type
- [ ] #6 Invalid or missing musicbrainz.txt is handled gracefully
- [ ] #7 Tests cover reading musicbrainz.txt and ID storage
<!-- AC:END -->
