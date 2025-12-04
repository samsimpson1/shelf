---
id: task-022
title: Display import directory name throughout import journey
status: To Do
assignee: []
created_date: '2025-12-04 21:34'
labels:
  - enhancement
  - ui
  - import
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, the import workflow doesn't consistently show which directory is being imported, making it easy for users to lose context during the multi-step process. The directory name should be displayed prominently on all import workflow pages (steps 1-5, confirm, and success pages) to help users track what they're importing.

This will improve user experience by providing clear context throughout the entire import process.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Directory name is displayed on all import workflow pages (step1, step2, step3, step4, step5, confirm, success)
- [ ] #2 Display is consistent and prominent across all pages
- [ ] #3 Directory name is clearly visible without cluttering the interface
- [ ] #4 E2E tests verify directory name appears on workflow pages
<!-- AC:END -->
