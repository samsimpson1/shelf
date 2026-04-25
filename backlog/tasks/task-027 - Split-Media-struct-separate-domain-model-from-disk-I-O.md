---
id: TASK-027
title: 'Split Media struct: separate domain model from disk-I/O'
status: To Do
assignee: []
created_date: '2026-04-24 20:49'
labels:
  - maintainability
  - refactor
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
`Media` in `models.go:74-166` mixes domain data with filesystem I/O. `LoadDescription`, `LoadGenres`, `LoadTrackList`, and `FindPosterFile` all read files using a path stored on the struct itself and silently return zero values on error (models.go:128, 138, 157). Consequences:

- every domain test becomes a filesystem test
- failures are hidden from handlers and users
- templates are coupled to disk layout

Move file-backed loading into a `MediaMetadataLoader` (or Scanner methods). Have handlers build view structs with metadata already resolved, and propagate errors so handlers can decide how to degrade.

Ties in with Task 8 (error-handling standardisation) — the new loader should return wrapped errors, not swallow them.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `Media` struct contains no `os.ReadFile` / `os.Stat` calls
- [ ] #2 Metadata loading lives in a loader type (or Scanner) that is testable without needing a real directory tree
- [ ] #3 Errors from metadata loading are propagated rather than silently returning zero values
- [ ] #4 Templates render from pre-populated view structs
- [ ] #5 Unit tests for the loader do not depend on fixtures created by other tests
- [ ] #6 `go test ./...` passes
- [ ] #7 Playwright E2E suite passes
<!-- AC:END -->
