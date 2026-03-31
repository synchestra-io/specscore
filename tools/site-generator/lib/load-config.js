import { readFile } from 'node:fs/promises';

/**
 * Loads site-config.json and builds derived lookup structures.
 * @param {URL} configPath - URL to site-config.json
 * @returns {{ pages: Array, sourceToSlug: Map<string, string>, sidebarGroups: Array<{label: string, items: Array}> }}
 */
export async function loadConfig(configPath) {
  const raw = await readFile(configPath, 'utf-8');
  const { pages } = JSON.parse(raw);

  const sourceToSlug = new Map();
  for (const page of pages) {
    sourceToSlug.set(page.source, page.slug);
  }

  // Build grouped sidebar structure from nav:true pages
  const navPages = pages
    .filter((p) => p.nav)
    .sort((a, b) => a.navOrder - b.navOrder);

  const groupMap = new Map();
  for (const page of navPages) {
    const label = page.navGroup || 'Other';
    if (!groupMap.has(label)) groupMap.set(label, []);
    groupMap.get(label).push(page);
  }

  const sidebarGroups = Array.from(groupMap.entries()).map(([label, items]) => ({
    label,
    items,
  }));

  return { pages, sourceToSlug, sidebarGroups };
}
