---
id: task-003
title: Add -help parameter to show configuration options
status: Done
assignee: []
created_date: '2025-11-08 16:00'
updated_date: '2025-11-09 01:06'
labels:
  - enhancement
  - usability
  - cli
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a `-help` command-line flag that displays all available configuration options (environment variables) and their default values.

Currently, the application accepts the following environment variables:
- `MEDIA_DIR` - Path to media backup directory (default: `/home/sam/Scratch/media/backup`)
- `PORT` - HTTP server port (default: `8080`)
- `TMDB_API_KEY` - TMDB API key for metadata fetching (optional)

The `-help` flag should display this information in a user-friendly format when the user runs `./shelf -help` or `./shelf --help`.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Running `./shelf -help` displays all environment variables with descriptions and default values
- [x] #2 Running `./shelf --help` works the same as `-help`
- [x] #3 Help output is clear, well-formatted, and easy to read
- [x] #4 Application exits after displaying help (does not start server)
- [x] #5 Existing functionality remains unchanged when no flags are provided
<!-- AC:END -->
