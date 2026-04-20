import { readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const SITE_URL = 'https://specscore.md';

/**
 * Builds a line entry for llms.txt.
 * @param {object} page - a page from site-config.json with llmsDescription set
 * @returns {string} a markdown list item
 */
function llmsEntry(page) {
  const url = `${SITE_URL}/${page.slug}.md`;
  return `- [${page.title}](${url}): ${page.llmsDescription}`;
}

/**
 * Generates the contents of public/llms.txt.
 * Pages with an `llmsDescription` field are included, grouped by `llmsGroup`.
 * @param {object} config - the loaded site config
 * @returns {string}
 */
export function buildLlmsIndex(config) {
  const entries = config.pages.filter((p) => p.llmsDescription && p.slug);

  const groupOrder = [];
  const grouped = new Map();
  for (const page of entries) {
    const group = page.llmsGroup || 'Pages';
    if (!grouped.has(group)) {
      grouped.set(group, []);
      groupOrder.push(group);
    }
    grouped.get(group).push(page);
  }

  const sections = groupOrder.map((group) => {
    const lines = grouped.get(group).map(llmsEntry).join('\n');
    return `## ${group}\n\n${lines}`;
  });

  const header = [
    '# SpecScore',
    '',
    '> SpecScore is a Markdown-based specification language for AI-assisted development. Specs are versioned files that humans author, AI agents read and write, and lint tooling validates. Every URL on this site has a raw-Markdown sibling at `/{slug}.md` for agent consumption.',
    '',
    '',
  ].join('\n');

  return `${header}${sections.join('\n\n')}\n`;
}

/**
 * Generates the contents of public/llms-full.txt.
 * Pages with `llmsFull: true` contribute their raw source markdown, concatenated
 * in the order they appear in site-config.json.
 * @param {object} config - the loaded site config
 * @param {string} rootDir - repo root, used to resolve page.source
 * @returns {Promise<string>}
 */
export async function buildLlmsFull(config, rootDir) {
  const entries = config.pages.filter((p) => p.llmsFull && p.source);

  const header = [
    '# SpecScore — Full Documentation',
    '',
    '> Concatenated source Markdown of every SpecScore specification plus install and demo guides. Intended as a one-shot context warm-up for AI agents. Canonical, rendered HTML and per-page raw Markdown live at https://specscore.md/.',
    '',
    '',
  ].join('\n');

  const blocks = [];
  for (const page of entries) {
    const raw = await readFile(join(rootDir, page.source), 'utf-8');
    blocks.push(`---\n\n## File: ${page.slug}\n\nSource: \`${page.source}\`\nURL: ${SITE_URL}/${page.slug}.md\n\n${raw.trimEnd()}\n`);
  }

  return `${header}${blocks.join('\n')}\n`;
}

/**
 * Writes public/llms.txt and public/llms-full.txt.
 * @param {object} config - loaded site config
 * @param {string} outputDir - public/ output directory
 * @param {string} rootDir - repo root
 */
export async function writeLlmsFiles(config, outputDir, rootDir) {
  const index = buildLlmsIndex(config);
  await writeFile(join(outputDir, 'llms.txt'), index, 'utf-8');

  const full = await buildLlmsFull(config, rootDir);
  await writeFile(join(outputDir, 'llms-full.txt'), full, 'utf-8');
}
