import { readFile, writeFile, mkdir, cp, rm } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { loadConfig } from './lib/load-config.js';
import { rewriteLinks } from './lib/rewrite-links.js';
import { renderMermaidBlocks } from './lib/render-mermaid.js';
import { renderMarkdownToHtml, injectIntoTemplate } from './lib/render-page.js';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT = join(__dirname, '..', '..');
const OUTPUT = join(ROOT, 'public');

async function build() {
  const config = await loadConfig(new URL('./site-config.json', import.meta.url));
  const template = await readFile(join(__dirname, 'template.html'), 'utf-8');

  // Clean and create output directory
  await rm(OUTPUT, { recursive: true, force: true });
  await mkdir(OUTPUT, { recursive: true });

  // Copy static assets
  await mkdir(join(OUTPUT, 'assets'), { recursive: true });
  await cp(
    join(__dirname, 'assets', 'style.css'),
    join(OUTPUT, 'assets', 'style.css')
  );

  console.log(`Building ${config.pages.length} pages...`);

  for (const page of config.pages) {
    const sourcePath = join(ROOT, page.source);
    const rawMarkdown = await readFile(sourcePath, 'utf-8');

    // 1. Rewrite links for HTML output
    const htmlMarkdown = rewriteLinks(rawMarkdown, page.source, config.sourceToSlug, 'html');

    // 2. Pre-render mermaid diagrams (on the HTML-targeted markdown)
    const mermaidRendered = await renderMermaidBlocks(htmlMarkdown);

    // 3. Render markdown to HTML
    const htmlContent = renderMarkdownToHtml(mermaidRendered);

    // 4. Inject into template
    const htmlPage = injectIntoTemplate(template, {
      title: page.title,
      content: htmlContent,
      slug: page.slug,
      navItems: config.navItems,
    });

    // 5. Write HTML file
    const htmlFile = join(OUTPUT, `${page.slug}.html`);
    await writeFile(htmlFile, htmlPage, 'utf-8');

    // 6. Rewrite links for markdown output and write .md file
    const mdContent = rewriteLinks(rawMarkdown, page.source, config.sourceToSlug, 'md');
    const mdFile = join(OUTPUT, `${page.slug}.md`);
    await writeFile(mdFile, mdContent, 'utf-8');

    console.log(`  ${page.slug}.html + ${page.slug}.md`);
  }

  console.log(`\nDone. Output: ${OUTPUT}`);
}

build().catch((err) => {
  console.error(err);
  process.exit(1);
});
