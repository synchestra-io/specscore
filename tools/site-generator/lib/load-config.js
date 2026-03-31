import { readFile } from 'node:fs/promises';

/**
 * Loads site-config.json and builds derived lookup structures.
 * @param {URL} configPath - URL to site-config.json
 * @returns {{ pages: Array, sourceToSlug: Map<string, string>, navItems: Array }}
 */
export async function loadConfig(configPath) {
  const raw = await readFile(configPath, 'utf-8');
  const { pages } = JSON.parse(raw);

  const sourceToSlug = new Map();
  for (const page of pages) {
    sourceToSlug.set(page.source, page.slug);
  }

  const navItems = pages
    .filter((p) => p.nav)
    .sort((a, b) => a.navOrder - b.navOrder);

  return { pages, sourceToSlug, navItems };
}
