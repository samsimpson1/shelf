---
id: task-006
title: Run Go tests on push via GitHub Actions
status: Done
assignee: []
created_date: '2025-11-08 21:17'
updated_date: '2025-11-09 00:25'
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
- [x] #1 GitHub Actions workflow file is created in the repository
- [x] #2 Tests run automatically on every push to any branch
- [x] #3 Tests run on multiple Go versions (at minimum latest stable)
- [x] #4 Workflow runs the full test suite including integration tests when TMDB_API_KEY is available
- [x] #5 Workflow status is visible in pull requests
- [ ] #6 Test failures prevent merge when branch protection is enabled
- [x] #7 Workflow includes race condition detection (-race flag)
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Summary

Created a GitHub Actions workflow file at `.github/workflows/go-tests.yml` that implements all required features:

### Features Implemented:
1. **Workflow file created**: `.github/workflows/go-tests.yml` added to the repository
2. **Automatic execution**: Triggers on every push to any branch and on pull requests
3. **Multiple Go versions**: Tests run on Go 1.21, 1.22, and 1.23 using matrix strategy
4. **Integration tests**: Separate job runs full test suite with TMDB_API_KEY when available
5. **PR visibility**: Workflow status automatically appears on pull requests
6. **Race detection**: All test runs include `-race` flag to detect race conditions
7. **Coverage reporting**: Latest Go version generates coverage report uploaded as artifact

### Workflow Structure:
- **Main test job**: Runs on all Go versions with race detection
- **TMDB integration job**: Conditionally runs when TMDB_API_KEY is configured

### Notes:
- Criterion #6 (preventing merge on failures) requires repository branch protection rules to be configured in GitHub settings - this is a repository setting, not something controlled by the workflow file itself
- The workflow has been committed and pushed, triggering the first CI run
- Coverage reports are generated and uploaded as artifacts for the latest Go version
<!-- SECTION:NOTES:END -->
