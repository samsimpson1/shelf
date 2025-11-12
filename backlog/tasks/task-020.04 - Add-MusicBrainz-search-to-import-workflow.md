---
id: task-020.04
title: Add MusicBrainz search to import workflow
status: To Do
assignee: []
created_date: '2025-11-12 13:28'
labels:
  - musicbrainz
  - music
  - import
  - workflow
dependencies: []
parent_task_id: task-020
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Extend the import workflow to support MusicBrainz release search for Music media, similar to how TMDB search works for Films/TV. When importing Music media, users should be able to search MusicBrainz by artist and release title, select the correct release, and have the `musicbrainz.txt` file created automatically.

This completes the import experience for Music media by enabling automatic metadata association during the import process.

**User Value:** Users can search and identify music releases during import, ensuring correct metadata is fetched without needing to manually look up MusicBrainz IDs.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Import workflow includes MusicBrainz search option for Music media type
- [ ] #2 Search form allows searching by artist and/or release title
- [ ] #3 Search results display: cover art thumbnail, artist, release title, date, format, and track count
- [ ] #4 User can select a release from search results
- [ ] #5 Confirmation page shows selected release details
- [ ] #6 musicbrainz.txt file is created with release ID upon import confirmation
- [ ] #7 Option to download metadata immediately (cover art and tracks) or defer until restart
- [ ] #8 MusicBrainz search only appears for Music media type in import workflow
- [ ] #9 Manual MusicBrainz ID entry option available
- [ ] #10 Tests cover import workflow with MusicBrainz search and ID assignment
<!-- AC:END -->
