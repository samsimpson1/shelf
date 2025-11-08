---
id: task-001.03
title: Create Web UI Templates for TMDB Search
status: To Do
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:39'
labels: []
dependencies:
  - task-001.01
parent_task_id: task-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Build the user-facing templates for searching, selecting, and confirming TMDB IDs. This includes updates to the detail page and new search/confirmation pages.

## Technical Scope
- Update detail.html to show TMDB ID status and action buttons
- Create search.html for displaying search form and results
- Create confirm.html for preview before saving
- Ensure responsive design and accessibility
- Add loading states and error message displays

## Dependencies
Requires TMDB Search API for result display structure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Detail page shows 'Search for TMDB ID' button when TMDB ID is missing
- [ ] #2 Detail page shows 'Change TMDB ID' button with warning when ID exists
- [ ] #3 Search page displays pre-filled search form with media title and year
- [ ] #4 Search results show poster thumbnails, titles, dates, and truncated overviews
- [ ] #5 Search results include 'Select This' button for each result
- [ ] #6 Manual entry section allows direct TMDB ID input (collapsed by default)
- [ ] #7 Confirmation page displays both media details and selected TMDB match
- [ ] #8 Confirmation page includes 'Download metadata now' checkbox option
- [ ] #9 All templates are mobile-responsive
- [ ] #10 Empty state message displays when no search results found
<!-- AC:END -->
