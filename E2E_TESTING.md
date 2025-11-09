# End-to-End Testing with Playwright

This document explains how to run and write E2E tests for the Shelf media backup manager.

## Overview

The project uses [Playwright](https://playwright.dev/) for end-to-end testing of the web interface. These tests verify complete user workflows by running a real browser and interacting with the application as a user would.

## Test Coverage

The E2E tests cover all critical user journeys:

1. **Media Details Viewing** - Viewing media with posters, descriptions, genres, and metadata
2. **Copy Play Commands** - Copying VLC and MPV commands from the disk list
3. **TMDB ID Management** - Setting new TMDB IDs and changing existing ones
4. **Import Workflows** - Importing films and TV shows with various disk types

## Setup

### Prerequisites

- Node.js 20 or later
- Go 1.25 or later
- npm or yarn

### Installation

1. Install Node.js dependencies:
   ```bash
   npm install
   ```

2. Install Playwright browsers:
   ```bash
   npx playwright install
   ```

3. Build the Go application:
   ```bash
   go build -o shelf .
   ```

## Running Tests

### Run All Tests

```bash
npm test
```

This command:
- Builds the Go application
- Starts the server on port 8080
- Runs all Playwright tests
- Shuts down the server when done

### Run Tests in Headed Mode

To see the browser while tests run:

```bash
npm run test:headed
```

### Run Tests in UI Mode

For interactive debugging with Playwright's UI:

```bash
npm run test:ui
```

### Run Tests in Debug Mode

To debug tests step-by-step:

```bash
npm run test:debug
```

### Run Specific Tests

Run a single test file:

```bash
npx playwright test e2e/media-details.spec.ts
```

Run tests matching a pattern:

```bash
npx playwright test --grep "TMDB"
```

### View Test Report

After running tests, view the HTML report:

```bash
npm run test:report
```

## Test Structure

### Test Files

All E2E tests are in the `e2e/` directory:

- `e2e/media-details.spec.ts` - Media detail page viewing tests
- `e2e/copy-play-command.spec.ts` - Copy command functionality tests
- `e2e/tmdb-management.spec.ts` - TMDB ID management tests
- `e2e/import-workflow.spec.ts` - Import workflow tests

### Test Fixtures

Test fixtures are **automatically generated** before tests run - they are NOT stored in the repository.

**How it works:**
- `e2e/setup-fixtures.ts` - Generates all test data (disk files, metadata, posters)
- `e2e/global-setup.ts` - Playwright global setup that runs fixture generation
- Fixtures are created in `e2e/fixtures/` at test time
- See [e2e/README.md](e2e/README.md) for details on how fixtures work

**Generated fixtures:**
- `e2e/fixtures/media/` - 3 sample media items (The Matrix, Breaking Bad, No TMDB Film)
- `e2e/fixtures/import/` - 3 import directories (raw_bluray, raw_dvd, raw_custom)

The test server is configured to use these auto-generated fixtures instead of real media directories.

## Writing New Tests

### Basic Test Structure

```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test('should do something specific', async ({ page }) => {
    // Navigate to a page
    await page.goto('/');

    // Interact with elements
    await page.click('button:has-text("Click Me")');

    // Verify expectations
    await expect(page.locator('h1')).toContainText('Expected Text');
  });
});
```

### Best Practices

1. **Use Descriptive Test Names** - Test names should clearly describe what they verify
2. **Test User Workflows** - Focus on complete user journeys, not individual functions
3. **Use Reliable Selectors** - Prefer text-based selectors over CSS classes when possible
4. **Handle Async Properly** - Always await async operations
5. **Verify Outcomes** - Each test should have clear expectations
6. **Isolate Tests** - Tests should not depend on each other
7. **Clean Up** - Test fixtures are automatically regenerated for each test run

### Common Patterns

#### Navigating to Pages

```typescript
await page.goto('/');
await page.goto('/media/the-matrix-1999');
```

#### Clicking Elements

```typescript
await page.click('button:has-text("Submit")');
await page.click('a.btn-primary');
```

#### Filling Forms

```typescript
await page.fill('input[name="title"]', 'The Matrix');
await page.fill('input[type="text"]', 'Value');
```

#### Waiting for Navigation

```typescript
await page.waitForURL(/\/media\/.+/);
```

#### Checking Visibility

```typescript
await expect(page.locator('.disk-table')).toBeVisible();
await expect(page.locator('h1')).toContainText('The Matrix');
```

#### Testing Clipboard

```typescript
// Grant clipboard permissions
await context.grantPermissions(['clipboard-read', 'clipboard-write']);

// Click copy button
await page.click('.copy-btn');

// Verify clipboard
const clipboardText = await page.evaluate(() => navigator.clipboard.readText());
expect(clipboardText).toContain('vlc');
```

## Configuration

The Playwright configuration is in `playwright.config.ts`:

- **Base URL**: `http://localhost:8080`
- **Test Directory**: `./e2e`
- **Browser**: Chromium (Desktop Chrome)
- **Retries**: 2 in CI, 0 locally
- **Web Server**: Automatically starts Go app with test fixtures

### Environment Variables

The test server uses these environment variables:

- `MEDIA_DIR=./e2e/fixtures/media` - Test media directory
- `IMPORT_DIR=./e2e/fixtures/import` - Test import directory
- `PORT=8080` - Server port

## Continuous Integration

E2E tests run automatically in GitHub Actions on every push and pull request.

### CI Configuration

See `.github/workflows/e2e-tests.yml` for the full CI configuration.

The CI workflow:
1. Checks out code
2. Sets up Go and Node.js
3. Builds the application
4. Installs Node.js dependencies (including `tsx`)
5. Installs Playwright browsers
6. Runs all tests (fixtures auto-generated via `globalSetup`)
7. Uploads test reports as artifacts

### Viewing CI Results

- Test results are available in the GitHub Actions tab
- Failed tests upload screenshots and videos
- HTML reports are available as downloadable artifacts

## Debugging Failed Tests

### Local Debugging

1. Run tests in headed mode to see the browser:
   ```bash
   npm run test:headed
   ```

2. Use debug mode for step-by-step execution:
   ```bash
   npm run test:debug
   ```

3. Add `await page.pause()` to pause execution at a specific point

### CI Debugging

1. Check the GitHub Actions logs for error messages
2. Download the test report artifact from the failed workflow
3. Extract and open `playwright-report/index.html` in a browser
4. Review screenshots and traces of failed tests

## Adding Tests for New Features

When adding a new feature, you must add corresponding E2E tests:

1. **Identify User Workflows** - What actions will users take?
2. **Create Test File** - Add a new `.spec.ts` file in `e2e/`
3. **Write Test Cases** - Cover happy paths and edge cases
4. **Add Fixtures** - Edit `e2e/setup-fixtures.ts` to add any necessary test data
5. **Verify Locally** - Run tests locally to ensure they pass
6. **Document** - Add test descriptions explaining what's being verified

### Example: Adding Tests for a New Feature

```typescript
// e2e/new-feature.spec.ts
import { test, expect } from '@playwright/test';

test.describe('New Feature', () => {
  test('should perform main workflow', async ({ page }) => {
    // 1. Navigate to feature
    await page.goto('/new-feature');

    // 2. Interact with UI
    await page.fill('input[name="data"]', 'test value');
    await page.click('button:has-text("Submit")');

    // 3. Verify outcome
    await expect(page.locator('.success')).toBeVisible();
    await expect(page.locator('.result')).toContainText('test value');
  });

  test('should handle errors gracefully', async ({ page }) => {
    // Test error scenarios
    await page.goto('/new-feature');
    await page.click('button:has-text("Submit")');

    // Verify error message
    await expect(page.locator('.error')).toBeVisible();
  });
});
```

## Resources

- [Playwright Documentation](https://playwright.dev/docs/intro)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Writing Tests Guide](https://playwright.dev/docs/writing-tests)

## Troubleshooting

### Tests Fail Locally But Pass in CI

- Check Node.js version matches CI (v20)
- Ensure Playwright browsers are installed: `npx playwright install`
- Verify test fixtures are generated: `npx tsx e2e/setup-fixtures.ts`
- Check that `tsx` is installed: `npm install`

### Tests Timeout

- Increase timeout in `playwright.config.ts`
- Check if server is starting correctly
- Verify baseURL is accessible: `curl http://localhost:8080`

### Flaky Tests

- Add explicit waits for elements: `await page.waitForSelector('.element')`
- Use `waitForURL` instead of `expect(page).toHaveURL` for navigation
- Increase retries in CI: `retries: 2`

### Browser Not Found

Run: `npx playwright install chromium`

## Summary

E2E tests provide confidence that the application works correctly from a user's perspective. When adding features or making changes:

1. Write E2E tests for new user workflows
2. Run tests locally before committing
3. Verify tests pass in CI before merging
4. Keep tests maintainable and focused on user value

For questions or issues, refer to the [Playwright documentation](https://playwright.dev/) or check existing test files for examples.
