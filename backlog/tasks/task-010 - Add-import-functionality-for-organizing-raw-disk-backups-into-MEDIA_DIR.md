---
id: task-010
title: Add import functionality for organizing raw disk backups into MEDIA_DIR
status: To Do
assignee: []
created_date: '2025-11-09 12:09'
updated_date: '2025-11-09 12:11'
labels: []
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a new import feature that allows users to take directories in IMPORT_DIR (raw disk backups from MakeMKV) and organize them into the correct structure in MEDIA_DIR.

The feature should guide users through:
1. Selecting a directory from IMPORT_DIR to import
2. Choosing media kind (TV or Film)
3. Optional TMDB search to match the media item
   - If TMDB match selected: use title and year from TMDB
   - If TMDB search skipped: prompt user for title and year manually
4. For TV: prompting for series and disk numbers
5. Choosing to add to existing media or create new media
6. Disk type detection with manual override options (Blu-Ray, Blu-Ray UHD, DVD, or custom text)
7. Moving/organizing the directory into the correct MEDIA_DIR structure

This will streamline the workflow of importing new disk backups into the media library.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 New IMPORT_DIR environment variable is configurable
- [ ] #2 Web UI provides import interface accessible from main page
- [ ] #3 User can select a directory from IMPORT_DIR to import
- [ ] #4 User is prompted to choose media kind (TV or Film)
- [ ] #5 Optional TMDB search integration allows finding and linking media
- [ ] #6 If TMDB match is selected, title and year are fetched from TMDB automatically
- [ ] #7 If TMDB search is skipped, user is prompted to enter title and year manually
- [ ] #8 For TV shows: user can specify series number and disk number
- [ ] #9 User can choose to add disk to existing media or create new media
- [ ] #10 Disk type is auto-detected from directory structure if possible
- [ ] #11 User can manually select disk type: Blu-Ray, Blu-Ray UHD, DVD, or custom text
- [ ] #12 Directory is moved/renamed to correct location in MEDIA_DIR with proper naming convention
- [ ] #13 TMDB ID is saved to tmdb.txt if TMDB match is selected
- [ ] #14 Success/error feedback is shown to user after import
- [ ] #15 Existing media library scanning is not affected
- [ ] #16 Import process validates directory structure before moving files

- [ ] #17 Tests cover import workflow including validation and file operations
<!-- AC:END -->
