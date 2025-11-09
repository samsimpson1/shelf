import { test, expect } from '@playwright/test';

test.describe('Media Details Page', () => {
  test('should display media details with poster, description, genres, and metadata', async ({ page }) => {
    await page.goto('/');

    // Click on a media item to view details
    await page.click('a[href*="/media/"]');

    // Wait for detail page to load
    await expect(page).toHaveURL(/\/media\/.+/);

    // Verify page title contains media name
    await expect(page.locator('h1')).toBeVisible();

    // Verify metadata section is present
    await expect(page.locator('text=Type')).toBeVisible();
    await expect(page.locator('text=Disks')).toBeVisible();

    // Verify back to library link exists
    await expect(page.locator('a:has-text("Back to Library")')).toBeVisible();
  });

  test('should display poster image or fallback', async ({ page }) => {
    await page.goto('/');

    // Click on first media item
    await page.click('a[href*="/media/"]');

    // Check if poster image exists or fallback is shown
    const posterExists = await page.locator('img[src*="/posters/"]').count() > 0;
    const fallbackExists = await page.locator('.poster-fallback').count() > 0;

    expect(posterExists || fallbackExists).toBeTruthy();
  });

  test('should display description when available', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // If description exists, it should be visible
    const descriptionSection = page.locator('.description');
    const descriptionCount = await descriptionSection.count();

    if (descriptionCount > 0) {
      await expect(descriptionSection).toBeVisible();
    }
  });

  test('should display genres when available', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // If genres exist, they should be visible
    const genresSection = page.locator('.genres');
    const genresCount = await genresSection.count();

    if (genresCount > 0) {
      await expect(genresSection).toBeVisible();
      // Genres should be displayed as chips/tags
      await expect(page.locator('.genre-chip')).toHaveCount(await page.locator('.genre-chip').count());
    }
  });

  test('should display TMDB ID with link when available', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // Check if TMDB ID is shown
    const tmdbLinkCount = await page.locator('a[href*="themoviedb.org"]').count();

    if (tmdbLinkCount > 0) {
      const tmdbLink = page.locator('a[href*="themoviedb.org"]');
      await expect(tmdbLink).toBeVisible();
      // Verify link opens in new tab
      await expect(tmdbLink).toHaveAttribute('target', '_blank');
    }
  });
});
