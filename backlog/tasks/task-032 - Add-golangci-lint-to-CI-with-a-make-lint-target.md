---
id: TASK-032
title: Add golangci-lint to CI with a make lint target
status: To Do
assignee: []
created_date: '2026-04-24 20:49'
labels:
  - maintainability
  - tooling
  - ci
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The repo has no linter configuration, no `go vet` step in CI, and the `Makefile` only exposes `test`, `shelf`, and `run`. This lets quality drift accrue — e.g. the silently-swallowed errors called out in Task 8 would be caught by `errcheck`.

Add:
- a `.golangci.yml` enabling at least `errcheck`, `govet`, `staticcheck`, `ineffassign`, `gofmt`
- a `make lint` target
- a GitHub Actions step (new or appended to existing workflow under `.github/workflows/`) that runs the linter and fails the build on violations

Expect to fix whatever the linter surfaces in the same PR. If Task 8 has landed first, this should be a small fix set.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `.golangci.yml` committed with at least `errcheck`, `govet`, `staticcheck`, `ineffassign`, `gofmt` enabled
- [ ] #2 `make lint` runs cleanly locally
- [ ] #3 CI runs `make lint` (or equivalent) and fails on lint violations
- [ ] #4 All existing lint violations fixed
- [ ] #5 `go test ./...` and the Playwright E2E suite remain green
<!-- AC:END -->
