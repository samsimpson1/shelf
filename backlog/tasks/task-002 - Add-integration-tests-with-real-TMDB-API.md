---
id: task-002
title: Add integration tests with real TMDB API
status: In Progress
assignee: []
created_date: '2025-11-08 15:53'
updated_date: '2025-11-08 16:11'
labels:
  - testing
  - tmdb
  - integration
  - api
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive integration tests that make real API calls to TMDB, requiring users to provide their own API key via the TMDB_API_KEY environment variable. Tests should be skipped when the API key is not provided or when running in short mode.

## Current State
- Existing tests use mock HTTP servers ([tmdb_test.go](tmdb_test.go:25-244))
- No tests validate actual TMDB API compatibility
- TMDB_API_KEY is already used in main.go for runtime configuration
- Integration tests exist but don't test TMDB API ([integration_test.go](integration_test.go))

## Test Structure Pattern
Tests should use Go's standard testing patterns:
- Check `testing.Short()` to skip in short mode
- Check `TMDB_API_KEY` env var and skip if not set
- Provide helpful skip messages to guide users
- Clean up temporary files after execution
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 New file tmdb_integration_test.go created with real API tests
- [ ] #2 Tests skip gracefully when TMDB_API_KEY is not set
- [ ] #3 Tests skip when running go test -short
- [ ] #4 Movie metadata fetching tested with known movie ID (e.g., 550 for Fight Club)
- [ ] #5 TV metadata fetching tested with known TV ID (e.g., 60059 for Better Call Saul)
- [ ] #6 Movie search by title and year tested
- [ ] #7 TV search by title tested
- [ ] #8 Poster download to temp directory tested and verified
- [ ] #9 Full metadata save workflow tested (poster + description + genres)
- [ ] #10 TMDB ID validation tested with real API
- [ ] #11 Error handling tested with invalid IDs
- [ ] #12 All tests clean up temporary files properly
- [ ] #13 CLAUDE.md updated with integration test documentation
- [ ] #14 Instructions added for running with TMDB_API_KEY=your_key go test -v
- [ ] #15 All integration tests pass when TMDB_API_KEY is provided
<!-- AC:END -->
