---
id: task-007
title: Fix tests failing when testdata directory is not populated correctly
status: To Do
assignee: []
created_date: '2025-11-09 00:34'
labels:
  - bug
  - testing
  - test-infrastructure
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The tests behave differently depending on the state of the testdata directory. This directory should be set up before tests are run and torn down again once testing is complete to ensure consistent state.

## Problem
- Tests have inconsistent behavior based on pre-existing testdata state
- No automated setup/teardown of test fixtures
- Tests may pass or fail depending on whether testdata exists or what state it's in

## Solution
Implement proper test setup and teardown:
- Create testdata directory structure programmatically in test setup
- Clean up testdata after tests complete
- Ensure each test run starts with a known, consistent state
- Consider using t.TempDir() or similar for isolated test fixtures
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Tests create their own testdata directory structure in setup phase
- [ ] #2 Tests clean up testdata directory after completion
- [ ] #3 Tests pass consistently regardless of pre-existing testdata state
- [ ] #4 CI/CD pipeline runs tests successfully with clean state
- [ ] #5 All existing tests continue to pass with new setup/teardown logic
<!-- AC:END -->
