---
id: task-013
title: Add MPV play command option alongside VLC
status: Done
assignee: []
created_date: '2025-11-09 14:49'
updated_date: '2025-11-09 14:56'
labels: []
dependencies:
  - task-009
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a second "Copy MPV Command" button next to the existing VLC play command button. This gives users the option to copy an mpv command instead of VLC for playing media disks.

## Context
Task 009 implemented a "Copy Play Command" button that copies VLC commands. This task adds an alternative mpv command option, allowing users to choose their preferred media player.

## MPV Command Format
The MPV button should copy a command in this format:
```
mpv bd:// --bluray-device="${PLAY_URL_PREFIX}/${MEDIA_DIR}/${DISK_DIR}"
```

**Important differences from VLC:**
- MPV uses `bd://` protocol (not `bluray://`)
- MPV requires `--bluray-device` flag with the path
- For DVDs, use `dvd://` with `--dvd-device` flag
- MPV command structure is different from VLC's simple protocol://path format

Example outputs:
```bash
# Blu-Ray disk
mpv bd:// --bluray-device="/mnt/media/War of the Worlds (2025) [Film]/Disk [Blu-Ray UHD]"

# DVD disk
mpv dvd:// --dvd-device="/mnt/media/Some Movie (2020) [Film]/Disk [DVD]"

# Other formats (fallback to direct file path)
mpv "/mnt/media/Some Movie (2020) [Film]/Disk [Other]"
```

## Required Changes

1. **Add new method to Disk struct** in [models.go](models.go):
   - Add `MPVPlayCommand(prefix string)` method to generate MPV commands
   - Detect Blu-Ray format → use `mpv bd:// --bluray-device="..."`
   - Detect DVD format → use `mpv dvd:// --dvd-device="..."`
   - Other formats → use `mpv "..."`  (direct file path)
   - Reuse existing format detection logic from `PlayCommand()`

2. **Update detail template** in [templates/detail.html](templates/detail.html):
   - Add second button "Copy MPV Command" next to existing VLC button
   - Call new JavaScript function `copyMPVCommand()` or reuse `copyPlayCommand()` with MPV command
   - Consider button styling to differentiate VLC vs MPV (or keep same style)
   - Maintain responsive layout (both buttons should work on mobile)

3. **Update/reuse JavaScript functionality** in [templates/detail.html](templates/detail.html):
   - Reuse existing `copyPlayCommand()` function (it's player-agnostic)
   - Both buttons can call the same clipboard function
   - No JavaScript changes needed if reusing the function

4. **Add comprehensive tests**:
   - Add tests for `MPVPlayCommand()` method in [models_test.go](models_test.go)
   - Test Blu-Ray format → `mpv bd:// --bluray-device="..."`
   - Test DVD format → `mpv dvd:// --dvd-device="..."`
   - Test other formats → `mpv "..."`
   - Test with and without PLAY_URL_PREFIX
   - Test path quoting/escaping if needed
   - Update template tests in [templates_test.go](templates_test.go) to verify both buttons present
   - **CRITICAL**: Ensure template test coverage does not decrease

## Design Decisions to Make

**Button Layout Options:**
- **Option A**: Two separate buttons side-by-side (simpler, clearer)
- **Option B**: Dropdown menu to select player (more compact, scalable)
- **Option C**: Toggle/radio buttons to switch between VLC/MPV (cleaner but less discoverable)

**Recommendation**: Start with Option A (two buttons) for simplicity and clarity. Users can clearly see both options.

## Protocol Reference

**VLC Format** (existing):
```bash
vlc bluray:///path/to/disk
vlc dvd:///path/to/disk
```

**MPV Format** (new):
```bash
mpv bd:// --bluray-device="/path/to/disk"
mpv dvd:// --dvd-device="/path/to/disk"
```
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 MPVPlayCommand() method added to Disk struct
- [x] #2 MPVPlayCommand() generates correct mpv bd:// command for Blu-Ray disks
- [x] #3 MPVPlayCommand() generates correct mpv dvd:// command for DVD disks
- [x] #4 MPVPlayCommand() handles other formats with direct file path
- [x] #5 MPVPlayCommand() correctly uses PLAY_URL_PREFIX when provided
- [x] #6 MPVPlayCommand() correctly quotes/escapes paths with spaces
- [x] #7 Copy MPV Command button displayed next to VLC button for each disk
- [x] #8 MPV button click copies mpv command to clipboard
- [x] #9 Visual feedback shown after successful copy (reuses existing toast)
- [x] #10 Both buttons work correctly on mobile layout
- [x] #11 All existing tests continue passing
- [x] #12 New tests added for MPVPlayCommand method covering all formats
- [x] #13 Template test coverage maintained or improved
<!-- AC:END -->
