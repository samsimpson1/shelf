import { execSync } from 'child_process';

/**
 * Playwright global setup - runs before all tests
 * This creates the test fixtures needed for E2E tests
 */
async function globalSetup() {
  console.log('Running fixture setup...');
  try {
    execSync('npx tsx e2e/setup-fixtures.ts', { stdio: 'inherit' });
  } catch (error) {
    console.error('Failed to setup fixtures:', error);
    throw error;
  }
}

export default globalSetup;
