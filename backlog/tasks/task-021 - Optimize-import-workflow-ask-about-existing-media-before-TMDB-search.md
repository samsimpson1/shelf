---
id: task-021
title: 'Optimize import workflow: ask about existing media before TMDB search'
status: To Do
assignee: []
created_date: '2025-12-04 21:33'
labels:
  - enhancement
  - ux
  - import
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The current import workflow has a redundant step when adding disks to existing media. It prompts users to search TMDB for metadata before asking if the disk belongs to existing media. This is inefficient because:

1. User searches TMDB and enters metadata
2. Later, user selects "Add to existing media"
3. The TMDB search was unnecessary since the existing media already has that information

**Current Flow:**
Step 1: Select directory
Step 2: Choose media type (Film/TV)
Step 3: Search TMDB or manual entry ← unnecessary if adding to existing
Step 4: Disk details
Step 5: Add to existing or create new ← decision made too late
Step 6: Confirm and execute

**Proposed Flow:**
Step 1: Select directory
Step 2: Choose media type (Film/TV)
Step 3: Add to existing or create new ← decision made earlier
Step 4a (if new): Search TMDB or manual entry
Step 4b (if existing): Select which media
Step 5: Disk details
Step 6: Confirm and execute

This saves time and cognitive load when adding additional disks to multi-disk releases or TV series.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When importing, users are asked "Add to existing or create new?" immediately after choosing media type
- [ ] #2 If adding to existing media, TMDB search step is skipped entirely
- [ ] #3 If creating new media, TMDB search/manual entry works as before
- [ ] #4 All existing import functionality remains working
- [ ] #5 E2E tests updated to cover both import paths
<!-- AC:END -->
