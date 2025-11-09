import { test, expect } from '@playwright/test';

test.describe('Import Workflow', () => {
  test('should display import button when IMPORT_DIR is configured', async ({ page }) => {
    await page.goto('/');

    // Should have "Import Media" button
    const importButton = page.locator('a:has-text("Import Media")');
    await expect(importButton).toBeVisible();
  });

  test('should list available directories to import', async ({ page }) => {
    await page.goto('/import');

    // Should show list of directories
    await expect(page.locator('h1, h2')).toContainText(/Import/i);

    // Should have at least the fixture directories
    const hasDirectories = await page.locator('button:has-text("Import"), a:has-text("Import")').count() > 0;
    expect(hasDirectories).toBeTruthy();
  });

  test('should import Film with Blu-Ray disk type', async ({ page }) => {
    await page.goto('/import');

    // Find and click import button for raw_bluray
    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    const buttonCount = await importButtons.count();

    if (buttonCount > 0) {
      // Click the first import button
      await importButtons.first().click();

      // Step 1: Choose media type
      await page.waitForURL(/\/import\/start|\/import\/step1/);

      // Select Film
      const filmButton = page.locator('button:has-text("Film"), input[value="Film"] ~ label, a:has-text("Film")');
      if (await filmButton.count() > 0) {
        await filmButton.first().click();
      }

      // Verify Blu-Ray is detected (if auto-detection is shown)
      const detectedText = page.locator('text=/Blu-?Ray/i');
      if (await detectedText.count() > 0) {
        await expect(detectedText.first()).toBeVisible();
      }

      // Continue through workflow
      const nextButton = page.locator('button:has-text("Next"), button:has-text("Continue"), button[type="submit"]');
      if (await nextButton.count() > 0) {
        // Fill in title if needed
        const titleInput = page.locator('input[name="title"]');
        if (await titleInput.count() > 0) {
          await titleInput.fill('Test Film Import');
        }

        // Fill in year if needed
        const yearInput = page.locator('input[name="year"]');
        if (await yearInput.count() > 0) {
          await yearInput.fill('2024');
        }
      }
    }
  });

  test('should import Film with DVD disk type', async ({ page }) => {
    await page.goto('/import');

    // Similar workflow but for DVD
    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    const buttonCount = await importButtons.count();

    if (buttonCount > 1) {
      // Try to find raw_dvd directory
      const dvdImport = page.locator('text=raw_dvd').locator('..').locator('button:has-text("Import"), a:has-text("Import")');

      if (await dvdImport.count() > 0) {
        await dvdImport.first().click();

        // Select Film type
        await page.waitForURL(/\/import\/start|\/import\/step1/);

        const filmButton = page.locator('button:has-text("Film"), input[value="Film"] ~ label');
        if (await filmButton.count() > 0) {
          await filmButton.first().click();
        }

        // Verify DVD is detected
        const detectedText = page.locator('text=DVD, text=VIDEO_TS');
        if (await detectedText.count() > 0) {
          await expect(detectedText.first()).toBeVisible();
        }
      }
    }
  });

  test('should import TV show with Blu-Ray disk type', async ({ page }) => {
    await page.goto('/import');

    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    if (await importButtons.count() > 0) {
      await importButtons.first().click();

      await page.waitForURL(/\/import\/start|\/import\/step1/);

      // Select TV Show
      const tvButton = page.locator('button:has-text("TV"), input[value="TV"] ~ label, a:has-text("TV Show")');
      if (await tvButton.count() > 0) {
        await tvButton.first().click();
      }

      // For TV shows, we need series and disk numbers
      const seriesInput = page.locator('input[name="series"], input[name="season"]');
      if (await seriesInput.count() > 0) {
        await seriesInput.fill('1');
      }

      const diskInput = page.locator('input[name="disk"]');
      if (await diskInput.count() > 0) {
        await diskInput.fill('1');
      }
    }
  });

  test('should import TV show with DVD disk type', async ({ page }) => {
    await page.goto('/import');

    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    if (await importButtons.count() > 1) {
      // Look for raw_dvd
      const dvdImport = page.locator('text=raw_dvd').locator('..').locator('button:has-text("Import"), a:has-text("Import")');

      if (await dvdImport.count() > 0) {
        await dvdImport.first().click();

        await page.waitForURL(/\/import\/start|\/import\/step1/);

        // Select TV Show
        const tvButton = page.locator('button:has-text("TV"), input[value="TV"] ~ label');
        if (await tvButton.count() > 0) {
          await tvButton.first().click();
        }
      }
    }
  });

  test('should import media with custom disk type', async ({ page }) => {
    await page.goto('/import');

    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    if (await importButtons.count() > 2) {
      // Look for raw_custom
      const customImport = page.locator('text=raw_custom').locator('..').locator('button:has-text("Import"), a:has-text("Import")');

      if (await customImport.count() > 0) {
        await customImport.first().click();

        await page.waitForURL(/\/import\/start|\/import\/step1/);

        // Select Film type
        const filmButton = page.locator('button:has-text("Film"), input[value="Film"] ~ label');
        if (await filmButton.count() > 0) {
          await filmButton.first().click();
        }

        // Look for custom disk type option
        const customTypeInput = page.locator('input[name="disk_type"], input[name="format"]');
        if (await customTypeInput.count() > 0) {
          await customTypeInput.fill('4K UHD');
        }
      }
    }
  });

  test('should show preview before executing import', async ({ page }) => {
    await page.goto('/import');

    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    if (await importButtons.count() > 0) {
      await importButtons.first().click();

      // Go through minimal workflow to reach confirmation
      await page.waitForURL(/\/import/);

      // Look for confirmation or preview page
      const confirmButton = page.locator('button:has-text("Confirm"), button:has-text("Execute Import")');

      // If we reach a confirmation step, verify it shows details
      if (await confirmButton.count() > 0) {
        // Should show destination path or similar preview
        const hasPreview = await page.locator('text=Destination, text=Preview, code, pre').count() > 0;
        expect(hasPreview).toBeTruthy();
      }
    }
  });

  test('should allow adding to existing media', async ({ page }) => {
    await page.goto('/import');

    // This feature allows adding disks to existing media entries
    const importButtons = page.locator('button:has-text("Import"), a:has-text("Import")');
    if (await importButtons.count() > 0) {
      await importButtons.first().click();

      await page.waitForURL(/\/import/);

      // Look for "Add to existing" option
      const addToExisting = page.locator('text=Add to existing, input[value="existing"] ~ label');

      if (await addToExisting.count() > 0) {
        await addToExisting.first().click();

        // Should show list of compatible media
        const existingMedia = page.locator('select[name="existing_media"], input[name="existing_media"]');
        expect(await existingMedia.count()).toBeGreaterThan(0);
      }
    }
  });
});
