---
id: TASK-026
title: Replace /media/* suffix-matching router with Go 1.22 pattern routes
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
`main.go:211-228` dispatches all `/media/*` sub-routes by calling `strings.HasSuffix` inside a closure. This hides the route table, scales poorly, and duplicates slug extraction across handlers.

The module is on `go 1.24.7`, which supports Go 1.22+ pattern routes on `http.ServeMux` (e.g. `mux.HandleFunc("GET /media/{slug}/search-tmdb", ...)`). Register each endpoint as its own route with method constraints where appropriate (POST-only for `set-tmdb`, `save-poster`).

Handlers can read the slug via `r.PathValue("slug")`, which reduces the need for the helper in Task 3 at these sites.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Each `/media/*` sub-route is a separately registered route with an explicit method constraint where applicable
- [ ] #2 The `strings.HasSuffix` dispatch closure in `main.go` is removed
- [ ] #3 Handlers read slug via `r.PathValue("slug")` where applicable
- [ ] #4 `go test ./...` passes
- [ ] #5 Playwright E2E suite passes
<!-- AC:END -->
