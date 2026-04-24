---
id: task-010
title: Add import functionality for organizing raw disk backups into MEDIA_DIR
status: Done
assignee: []
created_date: '2025-11-09 12:09'
updated_date: '2025-11-09 16:51'
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
- [x] #1 New IMPORT_DIR environment variable is configurable
- [x] #2 Web UI provides import interface accessible from main page
- [x] #3 User can select a directory from IMPORT_DIR to import
- [x] #4 User is prompted to choose media kind (TV or Film)
- [x] #5 Optional TMDB search integration allows finding and linking media
- [x] #6 If TMDB match is selected, title and year are fetched from TMDB automatically
- [x] #7 If TMDB search is skipped, user is prompted to enter title and year manually
- [x] #8 For TV shows: user can specify series number and disk number
- [x] #9 User can choose to add disk to existing media or create new media
- [x] #10 Disk type is auto-detected from directory structure if possible
- [x] #11 User can manually select disk type: Blu-Ray, Blu-Ray UHD, DVD, or custom text
- [x] #12 Directory is moved/renamed to correct location in MEDIA_DIR with proper naming convention
- [x] #13 TMDB ID is saved to tmdb.txt if TMDB match is selected
- [x] #14 Success/error feedback is shown to user after import
- [x] #15 Existing media library scanning is not affected
- [x] #16 Import process validates directory structure before moving files

- [x] #17 Tests cover import workflow including validation and file operations

- [x] #18 Directory and file names are sanitized to handle filesystem limitations (e.g., `:` replaced with `_`, other invalid characters handled)
- [x] #19 Sanitization is consistent with any existing naming conventions in the codebase
- [x] #20 Sanitized names are still readable and maintain media identification
<!-- AC:END -->
