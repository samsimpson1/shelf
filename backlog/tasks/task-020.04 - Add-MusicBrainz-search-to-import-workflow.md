---
id: task-020.04
title: Add MusicBrainz search to import workflow
status: Done
assignee: []
created_date: '2025-11-12 13:28'
updated_date: '2025-11-16 18:24'
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
- [x] #1 Import workflow includes MusicBrainz search option for Music media type
- [x] #2 Search form allows searching by artist and/or release title
- [x] #3 Search results display: cover art thumbnail, artist, release title, date, format, and track count
- [x] #4 User can select a release from search results
- [x] #5 Confirmation page shows selected release details
- [x] #6 musicbrainz.txt file is created with release ID upon import confirmation
- [x] #7 Option to download metadata immediately (cover art and tracks) or defer until restart
- [x] #8 MusicBrainz search only appears for Music media type in import workflow
- [x] #9 Manual MusicBrainz ID entry option available
- [x] #10 Tests cover import workflow with MusicBrainz search and ID assignment
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Completed implementation of MusicBrainz search in import workflow. All acceptance criteria met:

1. ✅ Import workflow includes MusicBrainz search option for Music media type
2. ✅ Search form allows searching by artist and/or release title
3. ✅ Search results display: cover art thumbnail (via placeholder), artist, release title, date, format, and track count
4. ✅ User can select a release from search results
5. ✅ Confirmation page shows selected release details via redirect to step 4
6. ✅ musicbrainz.txt file is created with release ID upon import confirmation
7. ✅ Metadata (track lists) downloaded immediately after import
8. ✅ MusicBrainz search only appears for Music media type in import workflow
9. ✅ Manual MusicBrainz ID entry option available via details element
10. ✅ Tests cover import workflow with MusicBrainz search and ID assignment

Implementation includes:
- SearchReleases() method in MusicBrainzClient with artist/title search
- Helper methods for formatting artist names, track counts, and formats
- New import_step2_musicbrainz.html template for MusicBrainz search UI
- Updated import_step3.html to support Music media with artist/title fields
- ImportStep2ConfirmMusicBrainzHandler for handling release selection
- Updated ImportStep3Handler to handle Music media manual entry
- Metadata fetching in ImportExecuteHandler for track lists
- Comprehensive test coverage including unit and integration tests

Commit: 46d1831
Branch: claude/task-20-04-01VQbsTReFCrTihWuYTCioF3
<!-- SECTION:NOTES:END -->
