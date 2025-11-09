---
id: task-009
title: Add 'Copy Play Command' button to disk list
status: To Do
assignee: []
created_date: '2025-11-09 11:09'
updated_date: '2025-11-09 11:15'
labels: []
dependencies:
  - task-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a button next to each disk in the disk list that copies a VLC play command to the clipboard. This feature requires the disk list from task-008 to be implemented first.

## Context
Once task-008 is complete, each disk will be displayed in the details page with its name, format, and size. This task adds an actionable button to quickly generate and copy a VLC command for playing the disk.

## Command Format
The button should copy a command in this format:
```
vlc bluray://${PLAY_URL_PREFIX}/${MEDIA_DIR}/${DISK_DIR}
```

Where:
- `PLAY_URL_PREFIX` - Environment variable for the URL prefix (e.g., network path or mount point)
- `MEDIA_DIR` - Relative path to the media directory
- `DISK_DIR` - Full path to the specific disk directory

Example output:
```
vlc bluray:///mnt/media/War of the Worlds (2025) [Film]/Disk [Blu-Ray UHD]
```

## Required Changes

1. **Add environment variable** in [main.go](main.go):
   - Add `PLAY_URL_PREFIX` environment variable (optional, defaults to empty string or local path)
   - Pass to handlers/templates for URL construction
   - Update `printHelp()` function ([main.go:14](main.go#L14)) to document the new PLAY_URL_PREFIX environment variable

2. **Update Disk struct** in [models.go](models.go):
   - Add `Path` field to store full disk directory path
   - Add `PlayCommand(prefix string)` method to generate VLC command string

3. **Update Scanner** in [scanner.go](scanner.go):
   - Store full disk directory path in `Disk.Path` when collecting disk information
   - Ensure paths are absolute or relative to a known base

4. **Update detail template** in [templates/detail.html](templates/detail.html):
   - Add "Copy Play Command" button for each disk in the disk list
   - Use JavaScript to copy command to clipboard on click
   - Show visual feedback (e.g., "Copied!" toast or button text change)
   - Handle different protocols based on disk format (bluray:// for Blu-Ray, dvd:// for DVD)

5. **Add JavaScript functionality**:
   - Implement clipboard copy using `navigator.clipboard.writeText()`
   - Add fallback for older browsers (execCommand method)
   - Show temporary success message after copy
   - Handle errors gracefully

6. **Update handlers** in [handlers.go](handlers.go):
   - Pass `PLAY_URL_PREFIX` to detail page template
   - Ensure disk paths are available in template data

7. **Update tests**:
   - Add tests for `PlayCommand()` method
   - Test environment variable handling
   - Test different disk formats generate correct protocol (bluray:// vs dvd://)
   - Update template tests in [templates_test.go](templates_test.go) to verify button presence
   - Add test for help text including PLAY_URL_PREFIX
   - **CRITICAL**: Ensure template test coverage does not decrease - current coverage includes comprehensive tests for all templates (index, detail, search, confirm) with various data states

## Protocol Selection
- Blu-Ray formats → `bluray://` protocol
- Blu-Ray UHD → `bluray://` protocol  
- DVD formats → `dvd://` protocol
- Other formats → `file://` protocol (fallback)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 PLAY_URL_PREFIX environment variable added and configurable
- [ ] #2 Help text updated to document PLAY_URL_PREFIX environment variable
- [ ] #3 Disk struct includes Path field storing disk directory path
- [ ] #4 Disk.PlayCommand() method generates correct VLC command format
- [ ] #5 Copy Play Command button displayed for each disk in the list
- [ ] #6 Button click copies VLC command to clipboard
- [ ] #7 Visual feedback shown after successful copy (toast/button text change)
- [ ] #8 Correct protocol used based on disk format (bluray://, dvd://, file://)
- [ ] #9 JavaScript clipboard API implemented with fallback for older browsers
- [ ] #10 Error handling for clipboard failures
- [ ] #11 PLAY_URL_PREFIX passed to templates and used in command generation
- [ ] #12 All existing tests updated and passing

- [ ] #13 New tests added for PlayCommand method and protocol selection

- [ ] #14 Template test coverage maintained or improved (no decrease in templates_test.go coverage)
<!-- AC:END -->
