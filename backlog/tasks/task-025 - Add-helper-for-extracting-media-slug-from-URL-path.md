---
id: TASK-025
title: Add helper for extracting media slug from URL path
status: Done
assignee: []
created_date: '2026-04-24 20:48'
updated_date: '2026-04-24 22:30'
labels:
  - maintainability
  - refactor
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The pattern `strings.TrimPrefix(r.URL.Path, "/media/")` followed by `strings.Split(..., "/")[0]` to extract a slug from the URL path is repeated at least 8 times across `handlers.go` — see handlers.go:155-156, 197-199, 272-279, 381-388 and elsewhere.

Introduce a single helper (e.g. `mediaSlugFromPath(r *http.Request) string`) and replace every occurrence. Note: if Task 4 (Go 1.22 pattern routes) lands first, many of these will disappear — this task should be coordinated with or follow that one, but can stand alone if needed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A single slug-extraction helper exists
- [x] #2 All duplicated slug-parsing call sites replaced
- [x] #3 Behaviour unchanged (same slug returned for same input)
- [x] #4 `go test ./...` passes
<!-- AC:END -->

## Implementation Notes

Added `mediaSlugFromPath(r *http.Request) string` near the top of [handlers.go](handlers.go). It trims the `/media/` prefix, drops a trailing `/`, then takes the first segment via `strings.Cut` — covering both `/media/{slug}` and `/media/{slug}/{action}` URL shapes with a single implementation.

Replaced six call sites: `DetailHandler`, `SearchTMDBHandler`, `ConfirmTMDBHandler`, `SaveTMDBHandler`, `SelectPosterHandler`, `SavePosterHandler`. Each now reads `slug := mediaSlugFromPath(r)` followed by the existing `if slug == "" { http.NotFound(...) ; return }` check.

`PosterHandler` was left alone — it lives at `/posters/`, not `/media/`, so it's outside the scope of a media-slug-specific helper. A more general `pathSlug(path, prefix)` helper would have folded it in too, but the task explicitly suggested the `mediaSlugFromPath(r *http.Request)` signature, so I stuck with that.

Behaviour notes:
- For valid URLs like `/media/foo` and `/media/foo/search-tmdb`, the returned slug is identical to the old code.
- The previous multi-segment handlers had a `len(parts) < 2` guard that 404'd URLs like `/media/search-tmdb` (no slug). The new helper would return `"search-tmdb"` for that path, but the subsequent `findMediaBySlug` lookup fails and the handler still returns 404 — same user-visible behaviour.

Verification: `go test ./...` → `ok shelf 9.355s`.
