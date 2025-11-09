---
id: task-018
title: Refresh media list when files in MEDIA_DIR change
status: To Do
assignee: []
created_date: '2025-11-09 23:26'
labels: []
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, the media list is scanned once on server startup and stored in memory. If files are added, removed, or modified in MEDIA_DIR outside of the import workflow (e.g., manual file operations, external scripts, network file sync), the changes are not reflected until the server is restarted.

Implement filesystem watching to automatically refresh the media list when changes are detected in MEDIA_DIR:
- Watch for new directories being created (new media added)
- Watch for directories being deleted (media removed)
- Watch for metadata file changes (tmdb.txt, poster.*, description.txt, genre.txt, title.txt)
- Watch for disk subdirectories being added/removed

Implementation approaches:
- Use filesystem watcher (fsnotify library or similar)
- Debounce rapid changes to avoid excessive rescans
- Optionally make this feature configurable (enable/disable via env var)
- Consider performance impact on large media libraries

This complements task-017 (refresh after import) by handling changes made outside the web UI.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Filesystem watcher is implemented for MEDIA_DIR
- [ ] #2 New media directories are automatically detected and added to the library
- [ ] #3 Deleted media directories are automatically removed from the library
- [ ] #4 Metadata file changes trigger appropriate updates (e.g., new poster appears)
- [ ] #5 Changes are debounced to prevent excessive rescanning
- [ ] #6 No performance degradation on large media libraries
- [ ] #7 Feature can be enabled/disabled via configuration if needed
<!-- AC:END -->
