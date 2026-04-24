---
id: task-022
title: Display import directory name throughout import journey
status: Done
assignee: []
created_date: '2025-12-04 21:34'
updated_date: '2025-12-04 21:46'
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
- [x] #1 Directory name is displayed on all import workflow pages (step1, step2, step3, step4, step5, confirm, success)
- [x] #2 Display is consistent and prominent across all pages
- [x] #3 Directory name is clearly visible without cluttering the interface
- [x] #4 E2E tests verify directory name appears on workflow pages
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Summary

Successfully added directory name display throughout the import workflow.

### Changes Made:

1. **Templates Updated** - Added consistent directory name display to all import workflow templates (step1-5, confirm, success)

2. **Handler Updates** - Modified import_handlers.go to pass directory name through workflow

3. **Display Style** - Consistent gray info box: "Importing: **[directory_name]**"

4. **E2E Tests** - Added comprehensive test verifying directory name appears on all workflow steps

### Files Modified:
- 7 template files (import_step2-5, step3_manual, confirm, success)
- import_handlers.go
- e2e/import-workflow.spec.ts

### User Experience:
Users now have clear context throughout the entire import process.
<!-- SECTION:NOTES:END -->
