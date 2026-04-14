import { readFile, writeFile, mkdir, cp, rm, readdir } from 'node:fs/promises';
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
    if (!page.source) continue;  // blog index — built separately

    // Landing page: inject hand-crafted HTML instead of rendering markdown
    if (page.slug === 'index') {
      const landingContent = await readFile(join(__dirname, 'landing.html'), 'utf-8');
      const htmlPage = injectIntoTemplate(template, {
        title: page.title,
        content: landingContent,
        slug: page.slug,
        sidebarGroups: config.sidebarGroups,
        eyebrow: '',
        showViewMarkdown: false,
      });
      await writeFile(join(OUTPUT, 'index.html'), htmlPage, 'utf-8');
      console.log('  index.html (landing)');
      continue;
    }

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
      sidebarGroups: config.sidebarGroups,
      eyebrow: page.navGroup || '',
    });

    // 5. Write HTML file
    const htmlFile = join(OUTPUT, `${page.slug}.html`);
    await mkdir(dirname(htmlFile), { recursive: true });
    await writeFile(htmlFile, htmlPage, 'utf-8');

    // 6. Rewrite links for markdown output and write .md file
    const mdContent = rewriteLinks(rawMarkdown, page.source, config.sourceToSlug, 'md');
    const mdFile = join(OUTPUT, `${page.slug}.md`);
    await mkdir(dirname(mdFile), { recursive: true });
    await writeFile(mdFile, mdContent, 'utf-8');

    console.log(`  ${page.slug}.html + ${page.slug}.md`);
  }

  // --- Blog ---
  const { parseBlogPost, buildBlogIndex } = await import('./lib/build-blog.js');

  const blogDir = join(ROOT, 'blog');
  let blogFiles = [];
  try {
    blogFiles = (await readdir(blogDir)).filter(f => f.endsWith('.md') && f !== 'README.md').sort();
  } catch {
    // blog/ directory doesn't exist yet — skip
  }

  if (blogFiles.length > 0) {
    const blogPostTemplate = await readFile(join(__dirname, 'blog-post.html'), 'utf-8');
    const blogIndexTemplate = await readFile(join(__dirname, 'blog-index.html'), 'utf-8');

    const posts = [];
    for (const file of blogFiles) {
      const raw = await readFile(join(blogDir, file), 'utf-8');
      const post = parseBlogPost(raw, file);
      posts.push(post);

      // Render individual post
      const htmlContent = renderMarkdownToHtml(post.body);
      const postContent = blogPostTemplate
        .replace('{{postDate}}', post.date)
        .replace('{{content}}', htmlContent);

      const htmlPage = injectIntoTemplate(template, {
        title: post.title,
        content: postContent,
        slug: post.slug,
        sidebarGroups: config.sidebarGroups,
        eyebrow: 'Blog',
      });

      const postFile = join(OUTPUT, `${post.slug}.html`);
      await mkdir(dirname(postFile), { recursive: true });
      await writeFile(postFile, htmlPage, 'utf-8');
      console.log(`  ${post.slug}.html (blog)`);
    }

    // Render blog index
    const sorted = buildBlogIndex(posts);
    const blogEntries = sorted.map(p =>
      `<li>
        <a href="/${p.slug}">${p.title}</a>
        <span class="blog-entry-date">${p.date}</span>
        ${p.description ? `<div class="blog-entry-desc">${p.description}</div>` : ''}
      </li>`
    ).join('\n');

    const indexContent = blogIndexTemplate.replace('{{blogEntries}}', blogEntries);
    const blogIndexPage = injectIntoTemplate(template, {
      title: 'Blog',
      content: indexContent,
      slug: 'blog',
      sidebarGroups: config.sidebarGroups,
      eyebrow: '',
      showViewMarkdown: false,
    });
    await writeFile(join(OUTPUT, 'blog.html'), blogIndexPage, 'utf-8');
    console.log('  blog.html (index)');
  }

  console.log(`\nDone. Output: ${OUTPUT}`);
}

build().catch((err) => {
  console.error(err);
  process.exit(1);
});
