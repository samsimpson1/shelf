---
id: task-008
title: List media disks in the details page
status: Done
assignee: []
created_date: '2025-11-09 11:07'
updated_date: '2025-11-09 11:50'
labels: []
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a detailed disk listing to the media details page that shows each individual disk with its name, format (Blu-Ray, Blu-Ray UHD, DVD), and file size in GB.

## Current State
- Media struct only tracks `DiskCount` (int) - total number of disks
- Scanner counts disks but doesn't collect individual disk information
- Detail page ([templates/detail.html:53](templates/detail.html#L53)) only displays total disk count
- Film disks match pattern: `^Disk \[.+\]$` (e.g., "Disk [Blu-Ray]")
- TV disks match pattern: `^Series (\d+) Disk (\d+) \[.+\]$` (e.g., "Series 1 Disk 1 [Blu-Ray UHD]")

## Required Changes
1. **Create Disk struct** in [models.go](models.go) to represent individual disks with fields:
   - Name (string) - disk name/identifier (e.g., "Disk 1", "Series 1 Disk 2")
   - Format (string) - disk format extracted from bracket notation (e.g., "Blu-Ray", "DVD", "Blu-Ray UHD")
   - SizeGB (float64) - disk size in gigabytes

2. **Update Media struct** in [models.go](models.go):
   - Add `Disks []Disk` field to store individual disk information
   - Keep `DiskCount` for backward compatibility (or replace with `len(Disks)`)

3. **Update Scanner** in [scanner.go](scanner.go):
   - Modify `countFilmDisks()` and `countTVDisks()` to collect detailed disk information
   - Parse disk format from directory names (extract content within brackets)
   - Calculate disk sizes by walking directory tree and summing file sizes
   - Store disk objects in Media.Disks slice

4. **Update detail template** in [templates/detail.html](templates/detail.html):
   - Add disk listing section below metadata grid
   - Display table/list showing: Disk Name, Format, Size (GB)
   - Format sizes nicely (e.g., "45.2 GB" with proper rounding)
   - Handle empty state (no disks found)

5. **Update tests**:
   - Add tests for Disk struct and methods
   - Update scanner tests to verify disk collection
   - Update handler/template tests for disk display
   - Add integration tests with various disk formats
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Disk struct created with Name, Format, and SizeGB fields
- [x] #2 Media struct includes Disks slice containing all disk information
- [x] #3 Scanner collects individual disk details (name, format, size) during scan
- [x] #4 Disk format correctly extracted from bracket notation in directory names
- [x] #5 Disk sizes calculated by walking directory tree and summing file sizes
- [x] #6 Detail page displays a list/table of all disks for the media item
- [x] #7 Disk list shows name, format, and size in GB for each disk
- [x] #8 File sizes formatted with proper precision (e.g., '45.2 GB')
- [x] #9 Both film and TV disk formats handled correctly
- [x] #10 All existing tests updated and passing
- [x] #11 New tests added for disk collection and display functionality
<!-- AC:END -->
