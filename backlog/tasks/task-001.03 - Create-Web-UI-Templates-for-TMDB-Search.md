---
id: task-001.03
title: Create Web UI Templates for TMDB Search
status: Done
assignee: []
created_date: '2025-11-08 15:39'
updated_date: '2025-11-08 15:56'
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
- [x] #1 Detail page shows 'Search for TMDB ID' button when TMDB ID is missing
- [x] #2 Detail page shows 'Change TMDB ID' button with warning when ID exists
- [x] #3 Search page displays pre-filled search form with media title and year
- [x] #4 Search results show poster thumbnails, titles, dates, and truncated overviews
- [x] #5 Search results include 'Select This' button for each result
- [x] #6 Manual entry section allows direct TMDB ID input (collapsed by default)
- [x] #7 Confirmation page displays both media details and selected TMDB match
- [x] #8 Confirmation page includes 'Download metadata now' checkbox option
- [x] #9 All templates are mobile-responsive
- [x] #10 Empty state message displays when no search results found
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully:

## Updated detail.html
- Added CSS styles for buttons and TMDB action section
- Added conditional display of TMDB ID buttons:
  - Shows "Search for TMDB ID" button (primary blue) when no TMDB ID exists
  - Shows "Change TMDB ID" button (secondary gray) when TMDB ID exists
  - Displays warning message when changing existing ID
- Buttons link to `/media/{slug}/search-tmdb` route
- Fully mobile-responsive with stacked buttons on small screens

## Created search.html
- Pre-filled search form with media title in query field
- Year input field for films (optional, shows media's year as placeholder)
- Primary "Search TMDB" button
- Collapsible "Manual Entry" section (collapsed by default via JavaScript)
- Manual entry allows direct TMDB ID input via POST to `/media/{slug}/set-tmdb`
- Search results displayed in a responsive grid:
  - Poster thumbnails (100px width, loads from TMDB w200 size)
  - Fallback emoji placeholders when no poster
  - Result title, date, and popularity score
  - Truncated overview text with CSS ellipsis
  - "Select This" button (green) linking to confirmation page
- Empty state message when no results found
- Error message display support
- Fully mobile-responsive (switches to single column on mobile)

## Created confirm.html
- Side-by-side comparison layout (2 columns)
- Left side: Current media details (title, year, type, disks, existing poster/description)
- Right side: TMDB match details (highlighted with green border):
  - TMDB ID, title, release/air date, popularity
  - Poster from TMDB (w300 size)
  - Full overview text
- Warning box displayed when overwriting existing TMDB ID
- Confirmation form with:
  - "Download metadata now" checkbox (checked by default)
  - Helper text explaining checkbox behavior
  - "Confirm and Save TMDB ID" button (green, bold)
  - "Cancel" button returning to search
- Hidden form fields to preserve TMDB ID and search query
- Fully mobile-responsive (stacks to single column)

All 10 acceptance criteria met. Templates follow existing design patterns, use consistent styling, and are fully responsive.
<!-- SECTION:NOTES:END -->
