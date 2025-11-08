---
id: task-006
title: Run Go tests on push via GitHub Actions
status: To Do
assignee: []
created_date: '2025-11-08 21:17'
labels:
  - ci/cd
  - testing
  - github-actions
  - automation
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Set up continuous integration to automatically run Go tests on every push to ensure code quality and catch regressions early. This will provide automated feedback on pull requests and prevent broken code from being merged.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GitHub Actions workflow file is created in the repository
- [ ] #2 Tests run automatically on every push to any branch
- [ ] #3 Tests run on multiple Go versions (at minimum latest stable)
- [ ] #4 Workflow runs the full test suite including integration tests when TMDB_API_KEY is available
- [ ] #5 Workflow status is visible in pull requests
- [ ] #6 Test failures prevent merge when branch protection is enabled
- [ ] #7 Workflow includes race condition detection (-race flag)
<!-- AC:END -->
