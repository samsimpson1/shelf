---
id: task-017
title: Refresh media list after import completes
status: To Do
assignee: []
created_date: '2025-11-09 23:25'
labels: []
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, after importing media through the web UI, the new media does not appear in the library until the server is restarted. This is because the media list is scanned once on startup and stored in memory.

Implement automatic media list refresh after a successful import so that:
- Newly imported media appears immediately in the library
- Users don't need to restart the server to see their imports
- The refresh is triggered automatically after the import execution completes

This could be implemented by:
- Re-running the scanner after successful import
- Adding the new media item directly to the in-memory slice
- Or providing a manual refresh button/endpoint

The solution should maintain thread-safety if multiple imports could happen concurrently.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 After completing an import, the new media appears in the library list without restarting the server
- [ ] #2 The import success page shows the newly imported media with correct metadata
- [ ] #3 No race conditions or data corruption when refreshing the media list
- [ ] #4 Solution handles both 'create new' and 'add to existing' import scenarios
- [ ] #5 Performance is acceptable (refresh doesn't cause noticeable delay)
<!-- AC:END -->
