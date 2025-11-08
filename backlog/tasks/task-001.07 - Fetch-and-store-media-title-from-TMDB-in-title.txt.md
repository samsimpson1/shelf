---
id: task-001.07
title: Fetch and store media title from TMDB in title.txt
status: To Do
assignee: []
created_date: '2025-11-08 16:31'
labels:
  - enhancement
  - tmdb
  - metadata
dependencies: []
parent_task_id: task-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add support for fetching the official title from TMDB API and saving it to `title.txt` file in each media directory. This extends the existing TMDB metadata fetching functionality to include titles alongside description.txt and genre.txt.

The title should be fetched from the TMDB API response (the "title" field for movies, "name" field for TV shows) and saved as plain text to `title.txt` in the media's directory.

This will help preserve the official/canonical title from TMDB, which may differ from the directory name parsed from the filesystem.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TMDB client fetches the official title (movie.title or tv.name) from API response
- [ ] #2 Title is saved to title.txt file in the media directory
- [ ] #3 Existing title.txt files are not overwritten (consistent with caching behavior)
- [ ] #4 Unit tests added to tmdb_test.go covering title fetching and saving
- [ ] #5 Integration test added to tmdb_integration_test.go for title metadata workflow
- [ ] #6 CLAUDE.md documentation updated to document title.txt file in directory structure section
<!-- AC:END -->
