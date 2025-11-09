---
id: task-015
title: Create Dockerfile for the shelf application
status: To Do
assignee: []
created_date: '2025-11-09 17:49'
labels: []
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a Dockerfile to containerize the shelf media backup manager application. The Dockerfile should:
- Use multi-stage build for optimal image size
- Build the Go binary in one stage
- Copy the binary and required assets (templates, static files) to a minimal runtime image
- Expose port 8080
- Support environment variables (MEDIA_DIR, IMPORT_DIR, PORT, TMDB_API_KEY)
- Follow Docker best practices
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Dockerfile exists in project root
- [ ] #2 Multi-stage build is used for optimal image size
- [ ] #3 All required files (templates/, static/) are copied to final image
- [ ] #4 Port 8080 is exposed
- [ ] #5 Environment variables are configurable
- [ ] #6 Image builds successfully with 'docker build'
- [ ] #7 Container runs and serves the web interface
<!-- AC:END -->
