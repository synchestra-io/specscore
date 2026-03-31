import { writeFile, readFile, mkdir, rm } from 'node:fs/promises';
import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { randomUUID } from 'node:crypto';

const execFileAsync = promisify(execFile);

const MERMAID_BLOCK_RE = /```mermaid\n([\s\S]*?)```/g;

/**
 * Finds all mermaid code blocks in markdown, renders each to SVG via mmdc,
 * and replaces the code block with the inline SVG.
 *
 * @param {string} markdown - source markdown content
 * @returns {string} markdown with mermaid blocks replaced by inline SVGs
 */
export async function renderMermaidBlocks(markdown) {
  const blocks = [...markdown.matchAll(MERMAID_BLOCK_RE)];
  if (blocks.length === 0) return markdown;

  const workDir = join(tmpdir(), `mermaid-${randomUUID()}`);
  await mkdir(workDir, { recursive: true });

  // Find mmdc binary and puppeteer config from node_modules
  const mmdc = new URL(
    '../node_modules/.bin/mmdc',
    import.meta.url
  ).pathname;
  const puppeteerConfig = new URL(
    '../puppeteer-config.json',
    import.meta.url
  ).pathname;

  try {
    let result = markdown;

    for (let i = blocks.length - 1; i >= 0; i--) {
      const block = blocks[i];
      const mermaidSource = block[1];
      const inputFile = join(workDir, `diagram-${i}.mmd`);
      const outputFile = join(workDir, `diagram-${i}.svg`);

      await writeFile(inputFile, mermaidSource, 'utf-8');
      try {
        await execFileAsync(mmdc, [
          '-i', inputFile,
          '-o', outputFile,
          '-b', 'transparent',
          '--quiet',
          '--puppeteerConfigFile', puppeteerConfig,
        ]);

        const svg = await readFile(outputFile, 'utf-8');
        // Replace the full code block (``` to ```) with the SVG
        result =
          result.slice(0, block.index) +
          svg.trim() +
          result.slice(block.index + block[0].length);
      } catch (err) {
        console.warn(`  warning: mermaid rendering failed for diagram ${i}, keeping code block`);
      }
    }

    return result;
  } finally {
    await rm(workDir, { recursive: true, force: true });
  }
}
