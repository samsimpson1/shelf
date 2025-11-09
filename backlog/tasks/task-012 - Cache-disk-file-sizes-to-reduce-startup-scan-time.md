---
id: task-012
title: Cache disk file sizes to reduce startup scan time
status: To Do
assignee: []
created_date: '2025-11-09 12:18'
labels: []
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, the scanner calculates disk directory sizes on every startup by walking all files in each disk directory using the `calculateDirSize()` function (scanner.go:156). Both `collectFilmDisks()` and `collectTVDisks()` call this function for each disk, which can significantly slow down startup times for large disk backups containing many gigabytes of media files.

The performance bottleneck occurs because `filepath.Walk()` must traverse every file in every disk directory on each application startup, even though disk sizes rarely change after initial backup.

**Solution**: Cache the calculated disk sizes in a `sizes.json` file stored in each media directory. The scanner should:
1. Check for `sizes.json` on startup
2. Load cached sizes if valid JSON exists
3. Fall back to calculating sizes only when cache is missing or invalid
4. Update the cache file after calculating new sizes

This will dramatically reduce startup time for media libraries with existing caches, while maintaining accuracy by recalculating when needed.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 sizes.json file stores disk sizes in bytes for each disk subdirectory in valid JSON format
- [ ] #2 Scanner checks for sizes.json existence before calculating disk sizes
- [ ] #3 Scanner loads and uses cached sizes when sizes.json exists and is valid JSON
- [ ] #4 Scanner falls back to calculateDirSize() when sizes.json is missing or contains invalid JSON
- [ ] #5 Scanner writes/updates sizes.json after calculating new disk sizes
- [ ] #6 JSON format is human-readable with proper indentation
- [ ] #7 Tests verify cache hit scenario (sizes loaded from existing valid cache)
- [ ] #8 Tests verify cache miss scenario (sizes calculated when cache missing)
- [ ] #9 Tests verify cache update scenario (cache file created/updated after calculation)
- [ ] #10 Tests verify invalid cache handling (graceful fallback to calculation)
- [ ] #11 Startup time improvement is measurable for directories with existing cache
<!-- AC:END -->
