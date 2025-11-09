import { execSync } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';

/**
 * Setup script to generate E2E test fixtures.
 * This creates all necessary test data including:
 * - Disk directories with dummy content for size calculation
 * - Metadata files (tmdb.txt, description.txt, genre.txt, title.txt)
 * - Poster images (1x1 pixel JPG placeholders)
 */

const fixturesDir = path.join(__dirname, 'fixtures');

interface FixtureFile {
  path: string;
  sizeMB: number;
}

interface MediaFixture {
  name: string;
  tmdbId?: string;
  title?: string;
  description?: string;
  genres?: string;
  disks: FixtureFile[];
}

// Media fixtures (for viewing/detail tests)
const mediaFixtures: MediaFixture[] = [
  {
    name: 'The Matrix (1999) [Film]',
    tmdbId: '603',
    title: 'The Matrix',
    description: 'Set in the 22nd century, The Matrix tells the story of a computer hacker who joins a group of underground insurgents fighting the vast and powerful computers who now rule the earth.',
    genres: 'Action, Science Fiction',
    disks: [
      {
        path: 'media/The Matrix (1999) [Film]/Disk [Blu-Ray]/BDMV/index.bdmv',
        sizeMB: 100, // 100MB = 0.1 GB visible in UI
      },
    ],
  },
  {
    name: 'Breaking Bad [TV]',
    tmdbId: '1396',
    title: 'Breaking Bad',
    description: 'When Walter White, a New Mexico chemistry teacher, is diagnosed with Stage III cancer and given a prognosis of only two years left to live. He becomes filled with a sense of fearlessness and an unrelenting desire to secure his family\'s financial future at any cost as he enters the dangerous world of drugs and crime.',
    genres: 'Drama, Crime',
    disks: [
      {
        path: 'media/Breaking Bad [TV]/Series 1 Disk 1 [Blu-Ray]/BDMV/index.bdmv',
        sizeMB: 100,
      },
      {
        path: 'media/Breaking Bad [TV]/Series 1 Disk 2 [DVD]/VIDEO_TS/VIDEO_TS.IFO',
        sizeMB: 100,
      },
    ],
  },
  {
    name: 'No TMDB Film (2020) [Film]',
    // No TMDB metadata for this one
    disks: [
      {
        path: 'media/No TMDB Film (2020) [Film]/Disk [DVD]/VIDEO_TS/VIDEO_TS.IFO',
        sizeMB: 100,
      },
    ],
  },
];

// Import fixtures (for import workflow tests)
const importFixtures: FixtureFile[] = [
  {
    path: 'import/raw_bluray/BDMV/index.bdmv',
    sizeMB: 50,
  },
  {
    path: 'import/raw_dvd/VIDEO_TS/VIDEO_TS.IFO',
    sizeMB: 50,
  },
  {
    path: 'import/raw_custom/content.mkv',
    sizeMB: 50,
  },
];

/**
 * Create a dummy file with specified size using dd (Unix) or Node.js buffer (Windows/fallback)
 */
function createDummyFile(filePath: string, sizeMB: number): void {
  const fullPath = path.join(fixturesDir, filePath);
  const dir = path.dirname(fullPath);

  // Ensure directory exists
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }

  // Skip if file already exists with correct size
  if (fs.existsSync(fullPath)) {
    const stats = fs.statSync(fullPath);
    const expectedSize = sizeMB * 1024 * 1024;
    if (Math.abs(stats.size - expectedSize) < 1024) {
      console.log(`  ✓ ${filePath} (already exists with correct size)`);
      return;
    }
  }

  try {
    // Try Unix dd command first
    if (process.platform !== 'win32') {
      execSync(
        `dd if=/dev/zero of="${fullPath}" bs=1M count=${sizeMB} 2>/dev/null`,
        { stdio: 'pipe' }
      );
    } else {
      // Windows fallback: use Node.js to create file
      const buffer = Buffer.alloc(sizeMB * 1024 * 1024);
      fs.writeFileSync(fullPath, buffer);
    }
    console.log(`  ✓ ${filePath} (${sizeMB}MB)`);
  } catch (error) {
    // Fallback: use Node.js to create file
    console.log(`  ⚠ Using fallback method for ${filePath}`);
    const buffer = Buffer.alloc(sizeMB * 1024 * 1024);
    fs.writeFileSync(fullPath, buffer);
    console.log(`  ✓ ${filePath} (${sizeMB}MB)`);
  }
}

/**
 * Create a text file with given content
 */
function createTextFile(filePath: string, content: string): void {
  const fullPath = path.join(fixturesDir, filePath);
  const dir = path.dirname(fullPath);

  // Ensure directory exists
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }

  fs.writeFileSync(fullPath, content, 'utf-8');
  console.log(`  ✓ ${filePath}`);
}

/**
 * Create a 1x1 pixel JPEG image (smallest valid JPEG)
 * This is a base64-encoded minimal JPEG image
 */
function createPosterImage(filePath: string): void {
  const fullPath = path.join(fixturesDir, filePath);
  const dir = path.dirname(fullPath);

  // Ensure directory exists
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }

  // Minimal 1x1 pixel JPEG (base64 encoded)
  const jpegBase64 = '/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8A/9k=';
  const jpegBuffer = Buffer.from(jpegBase64, 'base64');

  fs.writeFileSync(fullPath, jpegBuffer);
  console.log(`  ✓ ${filePath}`);
}

/**
 * Create metadata files for a media fixture
 */
function createMetadataFiles(mediaFixture: MediaFixture): void {
  const mediaDir = `media/${mediaFixture.name}`;

  if (mediaFixture.tmdbId) {
    createTextFile(`${mediaDir}/tmdb.txt`, mediaFixture.tmdbId);
  }

  if (mediaFixture.title) {
    createTextFile(`${mediaDir}/title.txt`, mediaFixture.title);
  }

  if (mediaFixture.description) {
    createTextFile(`${mediaDir}/description.txt`, mediaFixture.description);
  }

  if (mediaFixture.genres) {
    createTextFile(`${mediaDir}/genre.txt`, mediaFixture.genres);
  }

  // Always create a poster image
  createPosterImage(`${mediaDir}/poster.jpg`);
}

/**
 * Main setup function
 */
function setupFixtures(): void {
  console.log('Setting up E2E test fixtures...\n');

  console.log('Creating media fixtures:');
  for (const fixture of mediaFixtures) {
    console.log(`\n  ${fixture.name}:`);

    // Create disk files
    for (const disk of fixture.disks) {
      createDummyFile(disk.path, disk.sizeMB);
    }

    // Create metadata files
    createMetadataFiles(fixture);
  }

  console.log('\n\nCreating import fixture files:');
  for (const fixture of importFixtures) {
    createDummyFile(fixture.path, fixture.sizeMB);
  }

  console.log('\nFixtures setup complete!');
}

// Run setup
setupFixtures();
