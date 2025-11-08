---
id: task-001.04
title: Implement HTTP Handlers for TMDB Workflow
status: To Do
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:39'
labels: []
dependencies:
  - task-001.01
  - task-001.02
  - task-001.03
parent_task_id: task-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create HTTP handlers to connect the UI templates with the TMDB API and file management functionality. This implements the complete user workflow from search to save.

## Technical Scope
- SearchTMDBHandler (GET/POST) for search form and results
- ConfirmTMDBHandler (GET) for confirmation page
- SaveTMDBHandler (POST) for writing tmdb.txt and triggering metadata fetch
- Error handling and user feedback
- Input validation and sanitization

## Dependencies
Requires TMDB Search API, File Management, and Web UI Templates.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SearchTMDBHandler displays search form with pre-filled media title
- [ ] #2 SearchTMDBHandler processes search queries and renders results
- [ ] #3 ConfirmTMDBHandler validates TMDB ID and shows preview
- [ ] #4 SaveTMDBHandler writes tmdb.txt file to correct directory
- [ ] #5 SaveTMDBHandler triggers metadata fetch when checkbox is selected
- [ ] #6 SaveTMDBHandler redirects to detail page with success message
- [ ] #7 All handlers validate slug to prevent path traversal
- [ ] #8 Search queries are sanitized before passing to TMDB API
- [ ] #9 Network errors and API errors display user-friendly messages
- [ ] #10 Unit tests cover all handlers and error paths
<!-- AC:END -->
