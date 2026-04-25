---
id: TASK-033
title: 'Miscellaneous cleanup: temp dirs, dead constructors, personal default'
status: To Do
assignee: []
created_date: '2026-04-24 20:49'
labels:
  - maintainability
  - cleanup
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Three small unrelated hygiene items that are cheap enough to ship together:

1. **Convert remaining `os.MkdirTemp` usages to `t.TempDir()`** — `integration_test.go` at lines 66, 111, 174 still uses `os.MkdirTemp` with manual `defer os.RemoveAll`. Every other test uses `t.TempDir()`; bringing these in line removes cleanup debt and a foot-gun on early-return.
2. **Delete unused Scanner constructors** — `NewScanner` and `NewScannerWithTMDB` at scanner.go:31-41 are no longer called (only `NewScannerWithClients` is used). The package is `main`, so nothing external can consume them. Delete both.
3. **Remove the personal `MEDIA_DIR` default** — `main.go:92` defaults `MEDIA_DIR` to `/home/sam/Scratch/media/backup`. Either drop the default (fail fast with a clear error) or change it to something portable like `./media`. Dockerfile already sets its own default.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 No `os.MkdirTemp` calls remain in `*_test.go` files
- [ ] #2 `NewScanner` and `NewScannerWithTMDB` are deleted; `NewScannerWithClients` remains the only constructor
- [ ] #3 `MEDIA_DIR` default is either absent (fail fast) or portable
- [ ] #4 `go test ./...` passes
- [ ] #5 Docker image still builds and runs
<!-- AC:END -->
