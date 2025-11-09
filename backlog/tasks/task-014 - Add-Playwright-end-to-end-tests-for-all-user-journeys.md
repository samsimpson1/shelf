---
id: task-014
title: Add Playwright end-to-end tests for all user journeys
status: To Do
assignee: []
created_date: '2025-11-09 15:03'
labels: []
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement comprehensive end-to-end tests using Playwright to validate all critical user workflows in the media backup manager application. This will provide automated testing coverage for the web interface, ensuring user journeys work correctly across updates.

The tests should cover:
- Viewing media details and metadata
- Copying play commands from the disk list
- Managing TMDB IDs (both setting new IDs and changing existing ones)
- Import workflows for different media types (Film/TV) and disk formats (Blu-Ray, DVD, custom)

This will improve confidence in deployments and catch regressions in the user interface.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Playwright is installed and configured with appropriate test setup
- [ ] #2 Test exists for viewing media details page with poster, description, genres, and metadata
- [ ] #3 Test exists for copying play commands from the disk list on detail page
- [ ] #4 Test exists for changing TMDB ID of media that already has a TMDB ID set
- [ ] #5 Test exists for setting TMDB ID for media that does not currently have one
- [ ] #6 Test exists for importing a Film with Blu-Ray disk type
- [ ] #7 Test exists for importing a Film with DVD disk type
- [ ] #8 Test exists for importing a TV show with Blu-Ray disk type
- [ ] #9 Test exists for importing a TV show with DVD disk type
- [ ] #10 Test exists for importing media with custom/other disk type
- [ ] #11 Playwright tests can be run locally with a single command
- [ ] #12 Playwright tests are integrated into CI pipeline (GitHub Actions)
- [ ] #13 All tests pass consistently in both local and CI environments
- [ ] #14 Test documentation is added explaining how to run tests and write new ones
<!-- AC:END -->
