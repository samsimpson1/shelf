# End-to-End Tests

This directory contains Playwright E2E tests for the Shelf media backup manager.

## Test Fixtures

Test fixtures are **NOT** stored in the repository. Instead, they are automatically generated before tests run.

### How It Works

1. **[setup-fixtures.ts](setup-fixtures.ts)** - Generates all test data including:
   - Disk directories with dummy content (100MB for media, 50MB for import)
   - Metadata files (tmdb.txt, description.txt, genre.txt, title.txt)
   - Poster images (1x1 pixel JPEG placeholders)

2. **[playwright.config.ts](../playwright.config.ts)** - The `webServer` command runs `setup-fixtures.ts` before starting the Go server, ensuring fixtures exist before tests run

### Generated Fixtures

#### Media Fixtures (for viewing/detail tests)
- **The Matrix (1999) [Film]** - With TMDB metadata
- **Breaking Bad [TV]** - With TMDB metadata, 2 disks
- **No TMDB Film (2020) [Film]** - Without TMDB metadata

#### Import Fixtures (for import workflow tests)
- **raw_bluray/** - Blu-ray disk structure (BDMV)
- **raw_dvd/** - DVD disk structure (VIDEO_TS)
- **raw_custom/** - Custom format (MKV file)

### Adding New Fixtures

To add new test fixtures, edit [setup-fixtures.ts](setup-fixtures.ts):

1. Add to `mediaFixtures` array for media library tests
2. Add to `importFixtures` array for import workflow tests
3. Run `npx tsx e2e/setup-fixtures.ts` to test locally

### Running Tests Locally

```bash
# Install dependencies (first time)
npm install
npx playwright install chromium

# Build the Go application
go build -o shelf .

# Run all tests (fixtures auto-generated)
npm test

# Run in headed mode (see browser)
npm run test:headed

# Run in UI mode (interactive debugging)
npm run test:ui
```

### CI/CD

In CI, fixtures are generated automatically as part of the `webServer` command in playwright.config.ts. No manual setup required.

### Why Not Store Fixtures in Git?

Previously, test fixtures were stored in the repository using `.gitkeep` files and complex gitignore rules. This had several issues:

- **Large binary files** in git history
- **Complex .gitignore patterns** to exclude some files but keep others
- **Fixtures could get out of sync** between local and CI
- **Hard to modify** fixture data without committing large files

The new approach:
- **Zero fixture files in git** - entire `e2e/fixtures/` directory is ignored
- **Programmatic generation** - fixtures defined as data structures in TypeScript
- **Always fresh** - fixtures regenerated for every test run
- **Easy to modify** - just edit the fixture definitions
- **Fast** - fixture generation takes ~1 second

## Test Structure

- **[copy-play-command.spec.ts](copy-play-command.spec.ts)** - Copy VLC/MPV commands
- **[import-workflow.spec.ts](import-workflow.spec.ts)** - Import media workflow
- **[media-details.spec.ts](media-details.spec.ts)** - Media detail page viewing
- **[tmdb-management.spec.ts](tmdb-management.spec.ts)** - TMDB ID management

See [../CLAUDE.md](../CLAUDE.md) for full E2E testing documentation.
