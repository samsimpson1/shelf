---
id: task-001.01
title: Implement TMDB Search API Integration
status: To Do
assignee: []
created_date: '2025-11-08 15:39'
labels: []
dependencies: []
parent_task_id: task-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add TMDB API search functionality to the TMDBClient for searching movies and TV shows. This provides the foundation for users to find and select the correct TMDB ID for their media items.

## Technical Scope
- Add search methods to TMDBClient in tmdb.go
- Create search result structs for movies and TV shows
- Implement unified search method based on media type
- Handle pagination and result limiting (max 20 results)

## Dependencies
This is the foundation phase - no dependencies on other subtasks.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SearchMovies(query, year) method returns movie search results from TMDB API
- [ ] #2 SearchTV(query) method returns TV show search results from TMDB API
- [ ] #3 Search results include ID, title/name, release date, overview, poster path, and popularity
- [ ] #4 Results are sorted by popularity score
- [ ] #5 API errors are handled gracefully with meaningful error messages
- [ ] #6 Unit tests cover search methods with mock HTTP responses
- [ ] #7 Search result parsing handles missing fields (e.g., no poster path)
<!-- AC:END -->
