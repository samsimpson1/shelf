---
id: task-004
title: Create comprehensive template tests for all HTML templates
status: In Progress
assignee: []
created_date: '2025-11-08 16:22'
updated_date: '2025-11-08 21:12'
labels:
  - testing
  - templates
  - quality
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
## Overview
Create comprehensive tests for Go HTML templates to catch template errors before runtime. The project currently has 4 templates (index.html, detail.html, search.html, confirm.html) that are tested indirectly through handler tests, but template-specific errors (syntax, nil pointers, type assertions) are not caught until runtime.

## Current State
- **Templates:** index.html, detail.html, search.html, confirm.html
- **Current Testing:** handlers_test.go has basic handler tests but uses simplified templates, not the actual template files
- **Problem:** Template errors (syntax errors, field access errors, nil pointer dereferences) are discovered at runtime, not in tests

## Templates and Their Data Structures

### index.html
- Data: `struct { MediaList []Media }`
- Fields used: `.Type`, `.Slug`, `.PosterURL`, `.DisplayTitle`, `.DiskCount`
- Edge cases: empty list, items without posters

### detail.html
- Data: `struct { Media *Media; Description string; Genres []string; HasPoster bool }`
- Fields used: `.Media.DisplayTitle`, `.Media.Type`, `.Media.Year`, `.Media.DiskCount`, `.Media.TMDBID`, `.Media.PosterURL`, `.Media.Slug`
- Edge cases: missing description, empty genres, no poster, no TMDB ID

### search.html
- Data: `struct { Media *Media; Query string; Year int; Results interface{}; Error string }`
- Results type: `[]MovieSearchResult` or `[]TVSearchResult`
- Fields used: `.PosterPath`, `.GetTitle`, `.GetDate`, `.Popularity`, `.Overview`, `.ID`
- Edge cases: no results, search error, movie vs TV results

### confirm.html
- Data: `struct { Media *Media; TMDBID string; TMDBMatch interface{}; Query string; Description string; HasPoster bool; Error string }`
- TMDBMatch type: `MovieSearchResult` or `TVSearchResult`
- Fields used: `.TMDBMatch.Title`, `.TMDBMatch.Name`, `.TMDBMatch.ID`, `.TMDBMatch.ReleaseDate`, `.TMDBMatch.FirstAirDate`, `.TMDBMatch.Popularity`, `.TMDBMatch.PosterPath`, `.TMDBMatch.Overview`
- Edge cases: existing vs new TMDB ID, fetch error, movie vs TV match

## Implementation Plan

1. Create `templates_test.go` file
2. Add helper functions to build test data for each template
3. Test template parsing (catch syntax errors)
4. Test each template with various data scenarios:
   - Valid data (happy path)
   - Empty/nil optional fields
   - Edge cases (zero values, empty slices)
   - Type variations (Film vs TV, MovieSearchResult vs TVSearchResult)
5. Verify template output contains expected elements

## Benefits
- Catch template syntax errors at test time
- Prevent nil pointer panics from optional fields
- Verify type assertions work correctly
- Document expected data structures
- Regression protection for template changes
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 New templates_test.go file created with comprehensive template tests
- [ ] #2 All 4 templates (index, detail, search, confirm) have dedicated tests
- [ ] #3 Tests cover template parsing validation
- [ ] #4 Tests cover execution with valid data for each template
- [ ] #5 Tests cover edge cases: empty lists, nil pointers, missing optional fields
- [ ] #6 Tests cover type variations: Film vs TV, MovieSearchResult vs TVSearchResult
- [ ] #7 All tests pass successfully
- [ ] #8 Code coverage for template rendering paths improved
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create `templates_test.go` file
2. Add helper functions to build test data for each template
3. Test template parsing (catch syntax errors)
4. Test each template with various data scenarios:
   - Valid data (happy path)
   - Empty/nil optional fields
   - Edge cases (zero values, empty slices)
   - Type variations (Film vs TV, MovieSearchResult vs TVSearchResult)
5. Verify template output contains expected elements

## Benefits
- Catch template syntax errors at test time
- Prevent nil pointer panics from optional fields
- Verify type assertions work correctly
- Document expected data structures
- Regression protection for template changes
<!-- SECTION:DESCRIPTION:END -->
<!-- SECTION:PLAN:END -->
