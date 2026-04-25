---
id: TASK-024
title: Move importSessionStore off the package global onto App
status: Done
assignee: []
created_date: '2026-04-24 20:48'
updated_date: '2026-04-24 22:08'
labels:
  - maintainability
  - refactor
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
`import_handlers.go:57` declares `var importSessionStore = NewImportSessionStore()` as the only mutable package-level state in the handler layer. This couples every import handler to a global, is hostile to parallel tests, and muddies ownership.

Make the session store a field on `App`, initialised in `NewApp`. Update every call site in `import_handlers.go` to use `app.importSessions` (or similar).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 No package-level mutable handler state remains
- [x] #2 `importSessionStore` is a field on `App`, initialised in `NewApp`
- [x] #3 All call sites in `import_handlers.go` updated
- [x] #4 Tests can run in parallel without session-store interference
- [x] #5 `go test ./...` passes
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Moved `importSessionStore` from a package-global variable in `import_handlers.go` to a field `importSessions` on the `App` struct in `handlers.go`.

**Changes made:**
- Added `importSessions *ImportSessionStore` field to `App` struct (`handlers.go`)
- Initialized `importSessions: NewImportSessionStore()` in `NewApp` (`handlers.go`)
- Removed the global `var importSessionStore = NewImportSessionStore()` declaration (`import_handlers.go:57`)
- Replaced all 10 call sites in `import_handlers.go` from `importSessionStore.` to `app.importSessions.` (Create, Get, Delete operations across ImportStartHandler, ImportStep1Handler, ImportStep2Handler, ImportStep3Handler, ImportStep3ConfirmHandler, ImportStep3ManualHandler, ImportStep4Handler, ImportConfirmHandler, ImportExecuteHandler)

**Verification:**
- `go test ./...` passes (all ~170 tests)
- `go test -race ./...` passes with no race conditions detected
- No remaining code references to the old global variable
<!-- SECTION:FINAL_SUMMARY:END -->
