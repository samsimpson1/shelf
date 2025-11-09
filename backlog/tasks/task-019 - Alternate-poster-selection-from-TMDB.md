---
id: task-019
title: Alternate poster selection from TMDB
status: To Do
assignee: []
created_date: '2025-11-09 23:42'
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
- [ ] #1 Detail page shows 'Select Alternate Poster' button when media has a TMDB ID
- [ ] #2 Button is not displayed when TMDB API is not configured
- [ ] #3 Button is not displayed when media has no TMDB ID
- [ ] #4 Clicking button navigates to poster selection page showing TMDB's poster collection
- [ ] #5 Poster selection page displays current poster prominently for comparison
- [ ] #6 Alternate posters are shown in a responsive grid layout with thumbnails
- [ ] #7 Grid shows at least the English language posters, with option to show all languages
- [ ] #8 Clicking a poster thumbnail shows preview and 'Select This Poster' button
- [ ] #9 Selecting a poster downloads it from TMDB in original quality
- [ ] #10 New poster is saved as poster.{ext} in the media directory
- [ ] #11 Existing poster file is deleted before saving new one, even if extension differs
- [ ] #12 If current poster is poster.jpg and new is poster.png, the old .jpg file is removed
- [ ] #13 Media detail page reflects the new poster immediately after selection
- [ ] #14 Poster serving continues to work correctly with the new file
- [ ] #15 Error handling displays clear message if poster download fails
- [ ] #16 E2E test covers selecting an alternate poster for a film
- [ ] #17 E2E test covers selecting an alternate poster for a TV show
- [ ] #18 E2E test verifies old poster file with different extension is removed
<!-- AC:END -->
