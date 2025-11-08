---
id: task-001.02
title: Implement TMDB ID File Management
status: To Do
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:39'
labels: []
dependencies:
  - task-001.01
parent_task_id: task-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create functions to write, validate, and manage tmdb.txt files in media directories. This enables persistence of TMDB ID selections and ensures file operations are secure.

## Technical Scope
- WriteTMDBID function to create/overwrite tmdb.txt files
- ValidateTMDBID function to verify IDs exist and match media type
- File permission handling (0644)
- Path validation to prevent directory traversal

## Dependencies
Depends on TMDB Search API for validation functionality.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 WriteTMDBID creates tmdb.txt file with correct permissions (0644)
- [ ] #2 WriteTMDBID overwrites existing tmdb.txt files safely
- [ ] #3 ValidateTMDBID verifies TMDB ID exists via API
- [ ] #4 ValidateTMDBID prevents type mismatches (film ID for TV show, etc.)
- [ ] #5 Path validation prevents directory traversal attacks
- [ ] #6 File write errors are handled gracefully with clear messages
- [ ] #7 Unit tests cover file creation, validation, and error paths
<!-- AC:END -->
