---
id: TASK-031
title: Split overlong handler and fetch functions
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
Four functions each do several jobs in one body, making them hard to read and test in isolation:

- `ImportStep3ConfirmHandler` ~79 lines (import_handlers.go:396-474) — branches on Film/TV/Music, calls metadata APIs inline
- `ImportExecuteHandler` ~75 lines (import_handlers.go:703-778) — moves files, fetches metadata, refreshes list, cleans up session
- `FetchAndSaveMetadata` ~105 lines (tmdb.go:395-499) — four file-existence checks, API call, poster download, and three file writes
- `ExecuteImport` ~94 lines (import.go:204-297)

Extract per-type metadata fetching and discrete phases into named helpers. Aim for roughly ~40-line handlers; don't over-fragment.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Each of the four named functions is refactored into an orchestrator plus named helpers
- [ ] #2 No handler body exceeds ~40 lines
- [ ] #3 Behaviour is unchanged (same side effects, same responses, same error cases)
- [ ] #4 Helper functions are covered by tests where practical
- [ ] #5 `go test ./...` passes
- [ ] #6 Playwright E2E suite passes
<!-- AC:END -->
