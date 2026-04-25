---
id: TASK-030
title: Standardise error handling at the handler/library boundary
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
Error handling is inconsistent today:
- Scanner logs-and-swallows at scanner.go:138, 175, 214
- `ImportExecuteHandler` leaks raw errors into HTTP responses via `fmt.Sprintf("Import failed: %v", err)` at import_handlers.go:726
- `Media.LoadDescription` / `LoadGenres` / `LoadTrackList` return zero values on error (models.go:128, 138, 157)

Establish a per-layer rule: library code returns wrapped errors; handlers decide what the user sees versus what is logged. Introduce a small `writeError(w http.ResponseWriter, status int, userMsg string, err error)` helper for handlers.

Coordinates with Task 5 (Media loader) and Task 10 (golangci-lint `errcheck`).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Handlers never render raw `err.Error()` strings into HTTP response bodies
- [ ] #2 Scanner error paths either return wrapped errors or log with enough context that the caller understands the partial-failure state
- [ ] #3 Silently-zero-returning loader methods (if still present after Task 5) return wrapped errors instead
- [ ] #4 A `writeError` helper is used at every handler error site
- [ ] #5 Handler tests assert that user-facing error bodies do not leak internal paths or error chain detail
- [ ] #6 `go test ./...` passes
<!-- AC:END -->
