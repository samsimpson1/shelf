---
id: task-001.05
title: Add Routing and Integration for TMDB Feature
status: To Do
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:39'
labels: []
dependencies:
  - task-001.04
parent_task_id: task-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Wire up the HTTP handlers with routes in main.go and test the complete end-to-end workflow. This brings all components together into a working feature.

## Technical Scope
- Add routes in main.go for search, confirm, and save endpoints
- Integration testing of full workflow
- Verify TMDB API key configuration handling
- Test graceful degradation when API key is missing

## Dependencies
Requires all other subtasks to be complete.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Routes added: /media/{slug}/search, /media/{slug}/confirm, /media/{slug}/tmdb
- [ ] #2 Integration test covers full search-to-save workflow
- [ ] #3 Integration test verifies tmdb.txt file creation and content
- [ ] #4 Integration test verifies metadata fetch trigger
- [ ] #5 Integration test covers changing existing TMDB ID
- [ ] #6 Feature disabled gracefully when TMDB_API_KEY not set
- [ ] #7 All routes properly handle GET and POST methods
- [ ] #8 Manual testing checklist completed (see PLAN.md)
<!-- AC:END -->
