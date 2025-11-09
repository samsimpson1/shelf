# Running E2E Tests Locally

## Quick Start

```bash
# 1. Build the application
go build -o shelf .

# 2. Install dependencies
npm ci

# 3. Create test fixture content
dd if=/dev/zero of="e2e/fixtures/media/The Matrix (1999) [Film]/Disk [Blu-Ray]/BDMV/index.bdmv" bs=1M count=100 2>/dev/null
dd if=/dev/zero of="e2e/fixtures/media/Breaking Bad [TV]/Series 1 Disk 1 [Blu-Ray]/BDMV/index.bdmv" bs=1M count=100 2>/dev/null
dd if=/dev/zero of="e2e/fixtures/media/Breaking Bad [TV]/Series 1 Disk 2 [DVD]/VIDEO_TS/VIDEO_TS.IFO" bs=1M count=100 2>/dev/null
dd if=/dev/zero of="e2e/fixtures/media/No TMDB Film (2020) [Film]/Disk [DVD]/VIDEO_TS/VIDEO_TS.IFO" bs=1M count=100 2>/dev/null

# 4. Create import fixture directories (optional, for import tests)
mkdir -p e2e/fixtures/import/raw_bluray/BDMV
mkdir -p e2e/fixtures/import/raw_dvd/VIDEO_TS
mkdir -p e2e/fixtures/import/raw_custom

# 5. Install Playwright browsers
npx playwright install chromium

# 6. Run tests
npm test
```

## Troubleshooting

### Browser Crashes in Docker

If you see "Page crashed" errors when running tests in Docker or resource-constrained environments:

**Symptoms:**
- Tests fail with `page.goto: Page crashed`
- Tests timeout waiting for elements
- Chromium crashes immediately

**Cause:**
Chromium requires significant resources (memory, /dev/shm) that may not be available in Docker containers.

**Solutions:**

1. **Increase Docker resources** (recommended for local development):
   ```bash
   # If using Docker Desktop, increase memory to 4GB+
   # If using docker run, add:
   docker run --shm-size=2g ...
   ```

2. **Run tests in CI** (recommended approach):
   - Tests are configured to run in GitHub Actions with proper resources
   - CI environment has adequate memory and shared memory
   - See `.github/workflows/e2e-tests.yml`

3. **Run individual tests** instead of the full suite:
   ```bash
   npx playwright test e2e/media-details.spec.ts
   ```

4. **Use headed mode** to see what's happening:
   ```bash
   npm run test:headed
   ```

5. **Check system resources**:
   ```bash
   # Check available memory
   free -h

   # Check /dev/shm size (should be at least 1GB)
   df -h /dev/shm
   ```

### Verifying Server Works

If tests fail but you want to verify the server is working:

```bash
# Start server manually
MEDIA_DIR=./e2e/fixtures/media IMPORT_DIR=./e2e/fixtures/import PORT=8080 ./shelf &

# Test endpoints
curl http://localhost:8080/
curl http://localhost:8080/media/the-matrix-1999
curl http://localhost:8080/import

# Stop server
pkill shelf
```

## Test Structure

The E2E tests verify:
- ✅ Media detail page viewing
- ✅ Copy VLC/MPV commands
- ✅ TMDB ID management
- ✅ Import workflows

For detailed documentation, see [E2E_TESTING.md](E2E_TESTING.md).
