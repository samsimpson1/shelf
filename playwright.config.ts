import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['list']
  ],
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    launchOptions: {
      args: [
        '--disable-dev-shm-usage', // Avoid /dev/shm issues in Docker
        '--no-sandbox',            // Required for Docker environments
      ],
    },
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  /* Setup fixtures before starting tests */
  globalSetup: './e2e/global-setup.ts',

  /* Run the Go server before starting tests */
  webServer: {
    command: 'MEDIA_DIR=./e2e/fixtures/media IMPORT_DIR=./e2e/fixtures/import PORT=8080 TMDB_API_KEY=test_key_for_e2e ./shelf',
    url: 'http://localhost:8080',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});
