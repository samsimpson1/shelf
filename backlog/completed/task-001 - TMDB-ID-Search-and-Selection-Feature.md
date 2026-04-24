---
id: task-001
title: TMDB ID Search and Selection Feature
status: Done
assignee: []
created_date: '2025-11-08 15:38'
updated_date: '2025-11-09 12:25'
labels: []
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add functionality for users to search for and select TMDB IDs for media items that don't have one set, or when they want to change an existing ID. This feature will integrate with the TMDB API search endpoints to present options to users, then create/update the `tmdb.txt` file in the media directory.

## Implementation Approach
This functionality will be implemented as a **separate page** accessible from the media details page via a button (e.g., "Set TMDB ID" or "Change TMDB ID").

## Use Cases
1. **Missing TMDB ID**: Media items scanned without a `tmdb.txt` file need to be linked to TMDB
2. **Incorrect TMDB ID**: Wrong ID was manually set and needs correction
3. **Initial Setup**: Bulk assignment of TMDB IDs to an existing media library

## User Value
- Users can easily connect their media to TMDB without manually finding and entering IDs
- Prevents incorrect metadata by allowing search and preview before saving
- Enables users to fix mismatched IDs discovered after initial scan
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Separate TMDB ID search/selection page is created with its own route (e.g., /media/{slug}/set-tmdb-id)
- [x] #2 Media detail page has a button to navigate to the TMDB ID search page
- [x] #3 Users can search TMDB by title (and year for films) from the dedicated search page
- [x] #4 Search results display poster thumbnails, titles, release dates, and overview
- [x] #5 Users can select from search results to set TMDB ID
- [x] #6 Users can manually enter a TMDB ID if they know it
- [x] #7 Confirmation page shows preview of selected TMDB match before saving
- [x] #8 tmdb.txt file is created/updated correctly in media directory
- [x] #9 Users can change existing TMDB IDs with warning about metadata replacement
- [x] #10 Optional metadata fetch can be triggered after ID is set
- [x] #11 Feature gracefully degrades when TMDB_API_KEY is not configured
- [x] #12 All security considerations are addressed (path traversal, input validation, TMDB ID validation)

- [x] #13 Test coverage remains above 80%
- [x] #14 Documentation is updated with TMDB ID management instructions
<!-- AC:END -->
