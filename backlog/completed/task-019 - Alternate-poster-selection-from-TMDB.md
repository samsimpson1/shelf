---
id: task-019
title: Alternate poster selection from TMDB
status: Done
assignee: []
created_date: '2025-11-09 23:42'
updated_date: '2025-11-09 23:53'
labels:
  - feature
  - tmdb
  - ui
  - posters
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow users to browse and select alternate posters from TMDB's poster collection for their media items. TMDB provides multiple poster options for most movies and TV shows. Users should be able to view these alternate posters in a grid, select their preferred one, and have it saved to the media directory, replacing any existing poster file.

This feature enhances the user's ability to customize their media library with their preferred poster artwork, following the same UI patterns as the existing TMDB ID management workflow.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Detail page shows 'Select Alternate Poster' button when media has a TMDB ID
- [x] #2 Button is not displayed when TMDB API is not configured
- [x] #3 Button is not displayed when media has no TMDB ID
- [x] #4 Clicking button navigates to poster selection page showing TMDB's poster collection
- [x] #5 Poster selection page displays current poster prominently for comparison
- [x] #6 Alternate posters are shown in a responsive grid layout with thumbnails
- [x] #7 Grid shows at least the English language posters, with option to show all languages
- [x] #8 Clicking a poster thumbnail shows preview and 'Select This Poster' button
- [x] #9 Selecting a poster downloads it from TMDB in original quality
- [x] #10 New poster is saved as poster.{ext} in the media directory
- [x] #11 Existing poster file is deleted before saving new one, even if extension differs
- [x] #12 If current poster is poster.jpg and new is poster.png, the old .jpg file is removed
- [x] #13 Media detail page reflects the new poster immediately after selection
- [x] #14 Poster serving continues to work correctly with the new file
- [x] #15 Error handling displays clear message if poster download fails
- [x] #16 E2E test covers selecting an alternate poster for a film
- [x] #17 E2E test covers selecting an alternate poster for a TV show
- [x] #18 E2E test verifies old poster file with different extension is removed
<!-- AC:END -->
