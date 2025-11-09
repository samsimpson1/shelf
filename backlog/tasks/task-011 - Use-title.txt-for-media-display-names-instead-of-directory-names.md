---
id: task-011
title: Use title.txt for media display names instead of directory names
status: In Progress
assignee: []
created_date: '2025-11-09 12:13'
updated_date: '2025-11-09 12:22'
labels:
  - enhancement
  - tmdb
  - display
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, media titles are parsed from directory names (e.g., "The Thing (1982) [Film]" becomes "The Thing"). However, TMDB provides official titles that are saved to `title.txt` files. We should prefer these official titles when they exist, falling back to directory parsing only when `title.txt` is not present.

This will ensure that media is displayed with their correct, official titles from TMDB rather than potentially incorrect or non-standard directory names.

**Current Behavior:**
- Title is always parsed from directory name using regex
- `title.txt` files are saved during TMDB metadata fetch but never used

**Expected Behavior:**
- Check for `title.txt` in media directory during scan
- Use content from `title.txt` as the media title if file exists
- Fall back to directory name parsing if `title.txt` doesn't exist
- Maintain backward compatibility for media without TMDB metadata

**Files to Modify:**
- [scanner.go](scanner.go) - Update `ScanDirectory()` to read `title.txt` after parsing directory structure
- [scanner_test.go](scanner_test.go) - Add tests for title.txt reading and fallback behavior
- [test_helpers.go](test_helpers.go) - Update test fixtures to include title.txt files for testing
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Media with title.txt files display the TMDB official title
- [ ] #2 Media without title.txt files still display correctly using directory name parsing
- [ ] #3 Unit tests cover title.txt reading and fallback logic
- [ ] #4 Integration tests verify end-to-end behavior with both scenarios
- [ ] #5 No breaking changes to existing functionality
<!-- AC:END -->
