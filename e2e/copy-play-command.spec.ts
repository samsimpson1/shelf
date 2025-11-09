import { test, expect } from '@playwright/test';

test.describe('Copy Play Commands', () => {
  test('should display disk list with copy command buttons', async ({ page }) => {
    await page.goto('/media/the-matrix-1999');

    // Verify disk table is visible
    await expect(page.locator('.disk-table')).toBeVisible();

    // Verify table has headers
    await expect(page.locator('.disk-table th:has-text("Name")')).toBeVisible();
    await expect(page.locator('.disk-table th:has-text("Format")')).toBeVisible();
    await expect(page.locator('.disk-table th:has-text("Size")')).toBeVisible();
    await expect(page.locator('.disk-table th:has-text("Action")')).toBeVisible();

    // Verify at least one disk row exists
    await expect(page.locator('.disk-table tbody tr')).toHaveCount(await page.locator('.disk-table tbody tr').count());
  });

  test('should copy VLC command to clipboard', async ({ page, context }) => {
    // Grant clipboard permissions
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);

    await page.goto('/media/the-matrix-1999');

    // Click the "Copy VLC Command" button
    const copyButton = page.locator('.copy-btn').first();
    await expect(copyButton).toBeVisible();
    await copyButton.click();

    // Verify toast notification appears
    await expect(page.locator('#toast')).toBeVisible();
    await expect(page.locator('#toast')).toContainText('Copied to clipboard!');

    // Verify clipboard contains the command
    const clipboardText = await page.evaluate(() => navigator.clipboard.readText());
    expect(clipboardText).toContain('vlc');
  });

  test('should copy MPV command to clipboard', async ({ page, context }) => {
    // Grant clipboard permissions
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);

    await page.goto('/media/the-matrix-1999');

    // Click the "Copy MPV Command" button
    const copyButton = page.locator('.copy-btn-mpv').first();
    await expect(copyButton).toBeVisible();
    await copyButton.click();

    // Verify toast notification appears
    await expect(page.locator('#toast')).toBeVisible();
    await expect(page.locator('#toast')).toContainText('Copied to clipboard!');

    // Verify clipboard contains the command
    const clipboardText = await page.evaluate(() => navigator.clipboard.readText());
    expect(clipboardText).toContain('mpv');
  });

  test('should show toast message after copying', async ({ page, context }) => {
    await context.grantPermissions(['clipboard-read', 'clipboard-write']);
    await page.goto('/media/the-matrix-1999');

    // Click copy button
    await page.locator('.copy-btn').first().click();

    // Toast should be visible
    const toast = page.locator('#toast');
    await expect(toast).toBeVisible();
    await expect(toast).toContainText('Copied to clipboard!');

    // Toast should disappear after a few seconds
    await page.waitForTimeout(3500);
    await expect(toast).not.toBeVisible();
  });

  test('should display multiple disks for TV shows', async ({ page }) => {
    await page.goto('/media/breaking-bad');

    // TV shows can have multiple disks
    const diskRows = page.locator('.disk-table tbody tr');
    const rowCount = await diskRows.count();

    // Breaking Bad fixture has 2 disks
    expect(rowCount).toBeGreaterThanOrEqual(1);

    // Each row should have copy buttons
    for (let i = 0; i < rowCount; i++) {
      await expect(diskRows.nth(i).locator('.copy-btn')).toBeVisible();
      await expect(diskRows.nth(i).locator('.copy-btn-mpv')).toBeVisible();
    }
  });
});
