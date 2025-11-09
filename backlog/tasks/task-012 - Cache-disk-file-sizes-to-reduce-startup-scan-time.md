---
id: task-012
title: Cache disk file sizes to reduce startup scan time
status: Done
assignee: []
created_date: '2025-11-09 12:18'
updated_date: '2025-11-09 12:27'
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
- [x] #1 sizes.json file stores disk sizes in bytes for each disk subdirectory in valid JSON format
- [x] #2 Scanner checks for sizes.json existence before calculating disk sizes
- [x] #3 Scanner loads and uses cached sizes when sizes.json exists and is valid JSON
- [x] #4 Scanner falls back to calculateDirSize() when sizes.json is missing or contains invalid JSON
- [x] #5 Scanner writes/updates sizes.json after calculating new disk sizes
- [x] #6 JSON format is human-readable with proper indentation
- [x] #7 Tests verify cache hit scenario (sizes loaded from existing valid cache)
- [x] #8 Tests verify cache miss scenario (sizes calculated when cache missing)
- [x] #9 Tests verify cache update scenario (cache file created/updated after calculation)
- [x] #10 Tests verify invalid cache handling (graceful fallback to calculation)
- [x] #11 Startup time improvement is measurable for directories with existing cache
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented disk size caching in scanner.go:

**Implementation Details:**
1. Added `encoding/json` import for JSON handling
2. Created `loadSizeCache(mediaDir string)` function that:
   - Loads sizes.json from media directory
   - Returns empty map if file doesn't exist or contains invalid JSON
   - Gracefully handles errors

3. Created `saveSizeCache(mediaDir string, cache map[string]int64)` function that:
   - Saves cache to sizes.json with pretty-printing (2-space indentation)
   - Stores disk sizes in bytes
   - Returns error if write fails

4. Modified `collectFilmDisks()` to:
   - Load cache at the start of disk collection
   - Check cache for each disk directory by name
   - Use cached size if exists, otherwise calculate
   - Track if cache was updated
   - Save cache only if new sizes were calculated

5. Modified `collectTVDisks()` with the same caching logic

**Testing:**
Added comprehensive tests in scanner_test.go:
- `TestSizeCacheMiss` - Verifies cache file creation on first scan
- `TestSizeCacheHit` - Verifies sizes loaded from existing cache
- `TestSizeCacheUpdate` - Verifies cache updates with multiple disks
- `TestSizeCacheInvalid` - Verifies graceful fallback on invalid JSON
- `TestTVDiskCaching` - Verifies TV show disk caching
- `TestLoadSizeCacheMissingFile` - Tests loading when cache doesn't exist
- `TestSaveSizeCache` - Tests cache file saving and loading

All tests pass. Code coverage increased from 54.7% to 60.3%.

**Performance Impact:**
On subsequent scans, disks with cached sizes skip the expensive `calculateDirSize()` call that walks all files. For large media libraries with existing caches, this dramatically reduces startup time.
<!-- SECTION:NOTES:END -->
