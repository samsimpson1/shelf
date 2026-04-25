---
id: TASK-028
title: Deduplicate shared patterns between TMDB and MusicBrainz clients
status: To Do
assignee: []
created_date: '2026-04-24 20:49'
labels:
  - maintainability
  - refactor
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
`tmdb.go` and `musicbrainz.go` are structurally near-clones: each constructs its own `http.Client`, parses search responses, caps at 20 results, writes metadata files, and caches by checking file existence. See especially `FetchAndSaveMetadata` in tmdb.go:395-499 vs `FetchAndSaveTrackList` in musicbrainz.go:225-256.

Extract small shared helpers — not a forced unified API. Likely worth sharing:
- a `writeTextFile(dir, name, content string) error` helper
- poster/asset download to file
- injection of a shared `*http.Client` so both clients can be constructed with a single transport in tests

Goal is reducing duplication, not producing a `MetadataProvider` abstraction unless it falls out naturally. Stop when the remaining differences are load-bearing domain differences.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Both clients share a single implementation of the metadata-file-writing primitives
- [ ] #2 Both clients accept an injected `*http.Client`
- [ ] #3 Tests use a shared `*http.Client` / `httptest` transport where helpful
- [ ] #4 No behaviour changes in search, fetch, or caching semantics
- [ ] #5 `go test ./...` passes
<!-- AC:END -->
