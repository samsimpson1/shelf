---
id: task-020.02
title: Fetch cover art from MusicBrainz
status: To Do
assignee: []
created_date: '2025-11-12 13:28'
labels:
  - musicbrainz
  - music
  - metadata
  - cover-art
dependencies: []
parent_task_id: task-020
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement MusicBrainz Cover Art Archive API integration to automatically download cover art for Music media. When a `musicbrainz.txt` file exists with a valid release ID, fetch the cover art and save it as a poster file in the media directory.

The Cover Art Archive (coverartarchive.org) is MusicBrainz's official source for release cover art. This is analogous to TMDB poster fetching for Films/TV.

**User Value:** Users get automatic cover art for their music releases without manual downloads.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 MusicBrainz client can fetch cover art from Cover Art Archive API
- [ ] #2 Cover art is downloaded during initial scan when musicbrainz.txt exists
- [ ] #3 Cover art saved as poster.jpg/png/webp in media directory
- [ ] #4 Existing poster files are not overwritten (caching behavior)
- [ ] #5 Failed cover art fetches log warnings but don't stop scanning
- [ ] #6 Cover art display works in both grid and detail views
- [ ] #7 Tests cover cover art fetching and file saving
<!-- AC:END -->
