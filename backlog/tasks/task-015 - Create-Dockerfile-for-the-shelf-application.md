---
id: task-015
title: Create Dockerfile for the shelf application
status: In Progress
assignee: []
created_date: '2025-11-09 17:49'
updated_date: '2025-11-16 14:22'
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
- [x] #1 Dockerfile exists in project root
- [x] #2 Multi-stage build is used for optimal image size
- [x] #3 All required files (templates/, static/) are copied to final image
- [x] #4 Port 8080 is exposed
- [x] #5 Environment variables are configurable
- [ ] #6 Image builds successfully with 'docker build'
- [ ] #7 Container runs and serves the web interface
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Dockerfile created with multi-stage build using golang:1.21-alpine for building and alpine:latest for runtime. All required files (templates/, static/) are copied. Port 8080 is exposed. Environment variables (PORT, MEDIA_DIR, IMPORT_DIR, TMDB_API_KEY) are configurable with defaults. CA certificates are included for TMDB API HTTPS calls.

Note: Acceptance criteria #6 and #7 (actual Docker build and run testing) cannot be verified in this development environment as Docker is not available. The Dockerfile follows best practices and should build successfully when tested in an environment with Docker installed.
<!-- SECTION:NOTES:END -->
