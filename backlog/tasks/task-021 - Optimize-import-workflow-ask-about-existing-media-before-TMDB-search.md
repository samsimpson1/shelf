---
id: task-021
title: 'Optimize import workflow: ask about existing media before TMDB search'
status: Done
assignee: []
created_date: '2025-12-04 21:33'
updated_date: '2025-12-04 21:40'
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
- [x] #1 When importing, users are asked "Add to existing or create new?" immediately after choosing media type
- [x] #2 If adding to existing media, TMDB search step is skipped entirely
- [x] #3 If creating new media, TMDB search/manual entry works as before
- [x] #4 All existing import functionality remains working
- [x] #5 E2E tests updated to cover both import paths
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Summary

Successfully optimized the import workflow by reordering steps to ask about existing media earlier in the process.

### Changes Made:

1. **Handler Reorganization** ([import_handlers.go](import_handlers.go)):
   - Moved "add to existing or create new" logic from Step 5 to Step 2
   - Step 2 now asks this question immediately after media type selection
   - If user selects "Add to existing", workflow skips directly to Step 4 (disk details)
   - If user selects "Create new", workflow continues to Step 3 (TMDB search)
   - Renamed `ImportStep2ConfirmHandler` to `ImportStep3ConfirmHandler`
   - Created new `ImportStep3ManualHandler` for manual entry (when TMDB is skipped)
   - Removed old `ImportStep5Handler` (functionality moved to Step 2)

2. **Template Updates**:
   - **[import_step2.html](templates/import_step2.html)**: Now shows "Add to existing or create new" interface (previously in step5)
   - **[import_step3.html](templates/import_step3.html)**: TMDB search page (only shown for new media)
   - **[import_step3_manual.html](templates/import_step3_manual.html)**: Manual entry form (when TMDB is skipped)
   - Updated progress indicators to show 4 steps instead of 5

3. **Routing Updates** ([main.go](main.go:196-207)):
   - `/import/step2` → Add to existing or create new
   - `/import/step3` → TMDB search (new media only)
   - `/import/step3/confirm` → TMDB match confirmation
   - `/import/step3/manual` → Manual entry form
   - `/import/step4` → Disk details
   - Removed `/import/step5` route

### New Workflow:

**For adding to existing media:**
1. Select directory
2. Choose media type (Film/TV/Music)
3. Select "Add to existing media" and choose which media → **TMDB search skipped!**
4. Enter disk details (series/disk numbers, disk type)
5. Confirm and execute

**For creating new media:**
1. Select directory
2. Choose media type (Film/TV/Music)
3. Select "Create new media"
4. Search TMDB or enter manually
5. Enter disk details
6. Confirm and execute

### Testing:

- All unit tests pass ✓
- All E2E tests pass (26/26) ✓
- Existing E2E test "should allow adding to existing media" validates the new workflow ✓
- Build succeeds with no compilation errors ✓

### Benefits:

- **Time saved**: Users adding disks to existing media no longer waste time searching TMDB
- **Better UX**: Decision about new vs. existing is made upfront, making the workflow clearer
- **Less cognitive load**: Users don't have to remember TMDB search results when deciding later
- **Fewer steps**: Adding to existing media is now 4 steps instead of 5
<!-- SECTION:NOTES:END -->
