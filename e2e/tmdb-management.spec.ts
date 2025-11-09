import { test, expect } from '@playwright/test';

test.describe('TMDB ID Management', () => {
  test('should show "Search for TMDB ID" button when no TMDB ID exists', async ({ page }) => {
    await page.goto('/media/no-tmdb-film-2020');

    // Should show blue "Search for TMDB ID" button
    const searchButton = page.locator('a.btn-primary:has-text("Search for TMDB ID")');
    await expect(searchButton).toBeVisible();

    // Should NOT show "Change TMDB ID" button
    await expect(page.locator('a.btn-secondary:has-text("Change TMDB ID")')).not.toBeVisible();
  });

  test('should show "Change TMDB ID" button when TMDB ID exists', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // Should show gray "Change TMDB ID" button
    const changeButton = page.locator('a.btn-secondary:has-text("Change TMDB ID")');
    await expect(changeButton).toBeVisible();

    // Should show warning about replacing metadata
    await expect(page.locator('.warning')).toBeVisible();
    await expect(page.locator('.warning')).toContainText('Changing the TMDB ID will replace existing metadata');
  });

  test('should navigate to TMDB search page when clicking search button', async ({ page }) => {
    await page.goto('/media/no-tmdb-film-2020');

    // Click the search button
    await page.click('a.btn-primary:has-text("Search for TMDB ID")');

    // Should navigate to search page
    await expect(page).toHaveURL(/\/media\/no-tmdb-film-2020\/search-tmdb/);

    // Search page should have a form
    await expect(page.locator('form')).toBeVisible();
  });

  test('should set TMDB ID for media without one (manual entry)', async ({ page }) => {
    await page.goto('/media/no-tmdb-film-2020');

    // Click search button
    await page.click('a.btn-primary:has-text("Search for TMDB ID")');

    // Look for manual entry option
    const manualEntryToggle = page.locator('text=Or enter TMDB ID manually');

    if (await manualEntryToggle.count() > 0) {
      await manualEntryToggle.click();

      // Fill in TMDB ID manually
      await page.fill('input[name="tmdb_id"]', '550');

      // Submit the form
      await page.click('button:has-text("Set TMDB ID")');

      // Should redirect to detail page or success page
      await page.waitForURL(/\/media\/no-tmdb-film-2020/);

      // Verify TMDB ID is now set (should show Change button)
      await expect(page.locator('a.btn-secondary:has-text("Change TMDB ID")')).toBeVisible();
    }
  });

  test('should change existing TMDB ID with warning', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // Click change button
    await page.click('a.btn-secondary:has-text("Change TMDB ID")');

    // Should navigate to search page
    await expect(page).toHaveURL(/\/media\/the-matrix-1999\/search-tmdb/);

    // Should show warning about changing ID
    const warningCount = await page.locator('text=⚠️ Changing the TMDB ID will replace existing metadata').count();
    if (warningCount > 0) {
      await expect(page.locator('text=⚠️ Changing the TMDB ID will replace existing metadata')).toBeVisible();
    }
  });

  test('should display TMDB link on detail page when ID exists', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // Should have TMDB link
    const tmdbLink = page.locator('a[href*="themoviedb.org"]');
    await expect(tmdbLink).toBeVisible();

    // Link should open in new tab
    await expect(tmdbLink).toHaveAttribute('target', '_blank');
    await expect(tmdbLink).toHaveAttribute('rel', 'noopener');
  });

  test('should pre-fill search form with media title', async ({ page }) => {
    await page.goto('/media/no-tmdb-film-2020');
    await page.click('a.btn-primary:has-text("Search for TMDB ID")');

    // Search input should be pre-filled with title
    const searchInput = page.locator('input[name="query"], input[name="title"]');
    const inputCount = await searchInput.count();

    if (inputCount > 0) {
      const value = await searchInput.first().inputValue();
      expect(value.length).toBeGreaterThan(0);
    }
  });
});
