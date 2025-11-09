---
id: task-010
title: Add import functionality for organizing raw disk backups into MEDIA_DIR
status: To Do
assignee: []
created_date: '2025-11-09 12:09'
labels: []
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a new import feature that allows users to take directories in IMPORT_DIR (raw disk backups from MakeMKV) and organize them into the correct structure in MEDIA_DIR.

The feature should guide users through:
1. Selecting a directory from IMPORT_DIR
2. Choosing media kind (TV or Film)
3. For TV: prompting for series and disk numbers
4. For Films: prompting for year
5. Choosing to add to existing media or create new media
6. Optional TMDB search to match the media item
7. Disk type detection with manual override options (Blu-Ray, Blu-Ray UHD, DVD, or custom text)
8. Moving/organizing the directory into the correct MEDIA_DIR structure

This will streamline the workflow of importing new disk backups into the media library.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 New IMPORT_DIR environment variable is configurable
- [ ] #2 Web UI provides import interface accessible from main page
- [ ] #3 User can select a directory from IMPORT_DIR to import
- [ ] #4 User is prompted to choose media kind (TV or Film)
- [ ] #5 For TV shows: user can specify series number and disk number
- [ ] #6 For Films: user can specify year
- [ ] #7 User can choose to add disk to existing media or create new media
- [ ] #8 Optional TMDB search integration allows finding and linking media
- [ ] #9 Disk type is auto-detected from directory structure if possible
- [ ] #10 User can manually select disk type: Blu-Ray, Blu-Ray UHD, DVD, or custom text
- [ ] #11 Directory is moved/renamed to correct location in MEDIA_DIR with proper naming convention
- [ ] #12 TMDB ID is saved to tmdb.txt if TMDB match is selected
- [ ] #13 Success/error feedback is shown to user after import
- [ ] #14 Existing media library scanning is not affected
- [ ] #15 Import process validates directory structure before moving files
- [ ] #16 Tests cover import workflow including validation and file operations
<!-- AC:END -->
