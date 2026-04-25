---
id: TASK-029
title: Introduce constants for metadata filenames
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
Metadata filenames — `description.txt`, `genre.txt`, `title.txt`, `tmdb.txt`, `musicbrainz.txt`, `tracks.json` — appear as string literals in at least five files. See models.go:125, 135, 153; scanner.go:441; tmdb.go:400-420; import.go.

Renaming a file today requires grepping and editing every site. Introduce a single `const (...)` block (in `models.go` or a new `files.go`) and replace every literal.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 One `const (...)` block defines every metadata filename
- [ ] #2 No string literals for those filenames remain outside the const block
- [ ] #3 `go test ./...` passes
<!-- AC:END -->
