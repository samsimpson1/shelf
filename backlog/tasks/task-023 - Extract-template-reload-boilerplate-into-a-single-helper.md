---
id: TASK-023
title: Extract template-reload boilerplate into a single helper
status: Done
assignee: []
created_date: '2026-04-24 20:48'
updated_date: '2026-04-24 22:02'
labels:
  - maintainability
  - refactor
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The `if app.devMode { tmpl = app.loadTemplates() } else { tmpl = app.templates }` pattern is copy-pasted at ~15 render sites across `handlers.go` and `import_handlers.go`. In parallel, the 13-entry list of template files is duplicated between `main.go` (initial load) and `handlers.go` (dev-mode reload) — adding a template today requires editing both.

Replace the boilerplate with a single `(*App).getTemplates()` method and collapse the template file list to one source of truth.

Affected lines (non-exhaustive): handlers.go:116-118, 190-195, 266-270, 375-379, 568-579; import_handlers.go:67-71 and ~10 more sites; main.go:158-172; handlers.go:91-105.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 A single `(*App).getTemplates()` helper is used at every template render site
- [x] #2 The list of template files is defined in exactly one place
- [x] #3 Dev-mode hot-reload behaviour is unchanged
- [x] #4 `go test ./...` passes
- [x] #5 Playwright E2E suite passes against the built binary
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add a top-level `parseTemplates()` helper that loads `templates/*.html` via `template.ParseGlob` — single source of truth for the template set.
2. Add `(*App).getTemplates()` method: returns `app.templates` outside dev mode; in dev mode calls `parseTemplates()` with the same fallback-to-cached behaviour as today.
3. Replace every `tmpl := app.templates; if app.devMode { tmpl = app.loadTemplates() }` site in `handlers.go` and `import_handlers.go` with `tmpl := app.getTemplates()`. Remove the old `loadTemplates()` method.
4. Update `main.go`'s initial template load to call `parseTemplates()`.
5. Run `go test ./...` and the Playwright E2E suite.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Replaced the duplicated render-time `if app.devMode { tmpl = app.loadTemplates() } else { tmpl = app.templates }` block at every render site (5 in handlers.go, 8 in import_handlers.go) with a single `tmpl := app.getTemplates()` call.

Centralised the template set in [handlers.go](handlers.go) via a `templatesGlob = "templates/*.html"` constant plus a top-level `parseTemplates()` helper that returns `template.ParseGlob(templatesGlob)`. `(*App).getTemplates()` returns `app.templates` outside dev mode and re-parses via `parseTemplates()` in dev mode (with the same fallback-to-cached behaviour as the old `loadTemplates`).

`main.go`'s startup load now calls `parseTemplates()` too, so the 13-entry `template.ParseFiles` list disappears entirely. The stale `html/template` import in `main.go` was dropped as a result.

Side benefit: `import_step3_manual.html` (used by `ImportStep3ManualHandler`) was missing from both old `ParseFiles` lists; the glob picks it up automatically. Conversely the unused `import_step5.html` template is harmlessly loaded by the glob — a follow-up could delete the file but it's out of scope here.

Verification:
- `go test ./...` → ok shelf 9.317s
- `npm test` (Playwright) → 29 passed, 1 skipped, 10 failed. The same 10 tests fail on `main` (verified by re-running with the refactor stashed): they are pre-existing failures around TMDB UI / poster modal / import workflow, unrelated to the template-loading refactor.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

- Added `parseTemplates()` (uses `template.ParseGlob("templates/*.html")`) as the single source of truth for which templates are loaded, plus `(*App).getTemplates()` to encapsulate the dev-mode reload-vs-cached choice.
- Collapsed the ~15 copy-pasted `if app.devMode { ... } else { ... }` blocks across [handlers.go](handlers.go) and [import_handlers.go](import_handlers.go) into a single `tmpl := app.getTemplates()` call per render site.
- Deleted the duplicated 13-entry `template.ParseFiles` lists in `handlers.go` and `main.go`; both now go through `parseTemplates()`.

## Behaviour

Dev-mode hot-reload is preserved (re-parse on every render, fall back to cached set on parse error). Production path is unchanged: a single `parseTemplates()` call at startup, then served from `app.templates`.

## Tests

- `go test ./...` passes.
- Playwright suite: 29 passed, 1 skipped, 10 failed — the same 10 failures occur on `main` with this refactor stashed, so they are pre-existing and unrelated.
<!-- SECTION:FINAL_SUMMARY:END -->
