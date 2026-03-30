# Static Site Generator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Node.js script that converts `spec/` markdown into a static site with HTML + markdown output, pre-rendered mermaid diagrams, link rewriting, and dark/light theme support.

**Architecture:** A single `build.js` script reads a `site-config.json` mapping, processes each markdown source through `markdown-it`, rewrites internal links, pre-renders mermaid diagrams via `@mermaid-js/mermaid-cli`, and outputs HTML (wrapped in a template) plus processed `.md` files to `public/`. No framework — just a script with a few npm dependencies.

**Tech Stack:** Node.js, markdown-it, @mermaid-js/mermaid-cli (Puppeteer), vanilla CSS, vanilla JS (theme toggle only)

**Spec:** `docs/superpowers/specs/2026-03-30-static-site-generator-design.md`

---

## File Structure

```
tools/site-generator/
├── package.json              # npm project, scripts: build, clean
├── site-config.json          # URL mapping, page metadata, nav order
├── build.js                  # main build script (orchestrates all steps)
├── lib/
│   ├── load-config.js        # reads and validates site-config.json
│   ├── rewrite-links.js      # rewrites relative markdown links to site URLs
│   ├── render-mermaid.js     # extracts mermaid blocks, renders to SVG via mmdc
│   └── render-page.js        # markdown-it rendering + template injection
├── template.html             # HTML shell with nav, content slot, footer
└── assets/
    └── style.css             # all styling, light/dark themes via custom properties

public/                       # build output (.gitignored)
.gitignore                    # ignores public/ and node_modules
```

---

### Task 1: Project Scaffolding

**Files:**
- Create: `tools/site-generator/package.json`
- Create: `tools/site-generator/site-config.json`
- Create: `.gitignore`

- [ ] **Step 1: Create `package.json`**

```json
{
  "name": "specscore-site-generator",
  "version": "0.0.1",
  "private": true,
  "type": "module",
  "scripts": {
    "build": "node build.js",
    "clean": "rm -rf ../../public"
  },
  "dependencies": {
    "markdown-it": "^14.0.0",
    "@mermaid-js/mermaid-cli": "^11.4.0"
  }
}
```

Write this to `tools/site-generator/package.json`.

- [ ] **Step 2: Create `site-config.json`**

```json
{
  "pages": [
    {
      "source": "spec/README.md",
      "slug": "index",
      "title": "SpecScore",
      "nav": true,
      "navLabel": "Home",
      "navOrder": 0
    },
    {
      "source": "spec/features/README.md",
      "slug": "specifications",
      "title": "Specifications",
      "nav": true,
      "navLabel": "Specifications",
      "navOrder": 1
    },
    {
      "source": "spec/features/feature/README.md",
      "slug": "feature-specification",
      "title": "Feature Specification",
      "nav": true,
      "navLabel": "Feature",
      "navOrder": 2
    },
    {
      "source": "spec/features/development-plan/README.md",
      "slug": "development-plan-specification",
      "title": "Development Plan Specification",
      "nav": true,
      "navLabel": "Development Plan",
      "navOrder": 3
    },
    {
      "source": "spec/features/acceptance-criteria/README.md",
      "slug": "acceptance-criteria-specification",
      "title": "Acceptance Criteria Specification",
      "nav": true,
      "navLabel": "Acceptance Criteria",
      "navOrder": 4
    },
    {
      "source": "spec/features/project-definition/README.md",
      "slug": "project-definition-specification",
      "title": "Project Definition Specification",
      "nav": true,
      "navLabel": "Project Definition",
      "navOrder": 5
    },
    {
      "source": "spec/features/source-references/README.md",
      "slug": "source-references-specification",
      "title": "Source References Specification",
      "nav": true,
      "navLabel": "Source References",
      "navOrder": 6
    }
  ]
}
```

Write this to `tools/site-generator/site-config.json`.

- [ ] **Step 3: Create `.gitignore`**

```
public/
node_modules/
```

Write this to the repo root `.gitignore`.

- [ ] **Step 4: Install dependencies**

Run from `tools/site-generator/`:

```bash
npm install
```

Expected: `package-lock.json` created, `node_modules/` populated. Puppeteer browser downloaded (this may take a minute on first run).

- [ ] **Step 5: Commit**

```bash
git add .gitignore tools/site-generator/package.json tools/site-generator/package-lock.json tools/site-generator/site-config.json
git commit -m "feat(site): scaffold site generator project with config"
```

---

### Task 2: Config Loader

**Files:**
- Create: `tools/site-generator/lib/load-config.js`
- Create: `tools/site-generator/lib/load-config.test.js`

- [ ] **Step 1: Write the failing test**

Create `tools/site-generator/lib/load-config.test.js`:

```js
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { loadConfig } from './load-config.js';

describe('loadConfig', () => {
  it('loads pages from site-config.json', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(Array.isArray(config.pages));
    assert.ok(config.pages.length > 0);
    assert.equal(config.pages[0].slug, 'index');
    assert.equal(config.pages[0].source, 'spec/README.md');
  });

  it('builds a sourceToSlug map', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.equal(config.sourceToSlug.get('spec/README.md'), 'index');
    assert.equal(
      config.sourceToSlug.get('spec/features/feature/README.md'),
      'feature-specification'
    );
  });

  it('builds sorted nav items', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(config.navItems.length > 0);
    // nav items should be sorted by navOrder
    for (let i = 1; i < config.navItems.length; i++) {
      assert.ok(config.navItems[i].navOrder >= config.navItems[i - 1].navOrder);
    }
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run from `tools/site-generator/`:

```bash
node --test lib/load-config.test.js
```

Expected: FAIL — `Cannot find module './load-config.js'`

- [ ] **Step 3: Write implementation**

Create `tools/site-generator/lib/load-config.js`:

```js
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
```

- [ ] **Step 4: Run test to verify it passes**

```bash
node --test lib/load-config.test.js
```

Expected: 3 tests pass.

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/lib/load-config.js tools/site-generator/lib/load-config.test.js
git commit -m "feat(site): add config loader with sourceToSlug map and nav items"
```

---

### Task 3: Link Rewriting

**Files:**
- Create: `tools/site-generator/lib/rewrite-links.js`
- Create: `tools/site-generator/lib/rewrite-links.test.js`

- [ ] **Step 1: Write the failing test**

Create `tools/site-generator/lib/rewrite-links.test.js`:

```js
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { rewriteLinks } from './rewrite-links.js';

const sourceToSlug = new Map([
  ['spec/README.md', 'index'],
  ['spec/features/README.md', 'specifications'],
  ['spec/features/feature/README.md', 'feature-specification'],
  ['spec/features/development-plan/README.md', 'development-plan-specification'],
]);

describe('rewriteLinks', () => {
  it('rewrites relative links to site slugs for HTML', () => {
    const md = 'See [Features](features/README.md) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [Features](/specifications) for details.');
  });

  it('rewrites relative links to .md for markdown output', () => {
    const md = 'See [Features](features/README.md) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'md');
    assert.equal(result, 'See [Features](/specifications.md) for details.');
  });

  it('rewrites parent-relative links', () => {
    const md = 'See [Dev Plan](../development-plan/README.md) for planning.';
    const result = rewriteLinks(md, 'spec/features/feature/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [Dev Plan](/development-plan-specification) for planning.');
  });

  it('leaves external URLs unchanged', () => {
    const md = 'Visit [GitHub](https://github.com/synchestra-io/specscore).';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'Visit [GitHub](https://github.com/synchestra-io/specscore).');
  });

  it('leaves anchor-only links unchanged', () => {
    const md = 'See [below](#behavior) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [below](#behavior) for details.');
  });

  it('leaves links to unknown sources unchanged', () => {
    const md = 'See [unknown](../proposals/README.md) for proposals.';
    const result = rewriteLinks(md, 'spec/features/feature/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [unknown](../proposals/README.md) for proposals.');
  });

  it('handles links with anchors, preserving the fragment', () => {
    const md = 'See [statuses](../feature/README.md#feature-statuses).';
    const result = rewriteLinks(md, 'spec/features/development-plan/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [statuses](/feature-specification#feature-statuses).');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
node --test lib/rewrite-links.test.js
```

Expected: FAIL — `Cannot find module './rewrite-links.js'`

- [ ] **Step 3: Write implementation**

Create `tools/site-generator/lib/rewrite-links.js`:

```js
import { posix } from 'node:path';

/**
 * Rewrites markdown links from spec-relative paths to site URLs.
 *
 * @param {string} markdown - source markdown content
 * @param {string} sourcePath - repo-relative path of the source file (e.g. 'spec/features/feature/README.md')
 * @param {Map<string, string>} sourceToSlug - map from source paths to output slugs
 * @param {'html' | 'md'} format - output format determines link suffix
 * @returns {string} markdown with rewritten links
 */
export function rewriteLinks(markdown, sourcePath, sourceToSlug, format) {
  const sourceDir = posix.dirname(sourcePath);

  // Match markdown links: [text](url) — but not images ![text](url)
  return markdown.replace(/(?<!!)\[([^\]]*)\]\(([^)]+)\)/g, (match, text, href) => {
    // Skip external URLs and anchor-only links
    if (href.startsWith('http://') || href.startsWith('https://') || href.startsWith('#')) {
      return match;
    }

    // Split off anchor fragment
    const hashIndex = href.indexOf('#');
    const rawPath = hashIndex >= 0 ? href.slice(0, hashIndex) : href;
    const fragment = hashIndex >= 0 ? href.slice(hashIndex) : '';

    // Resolve relative path against source file's directory
    const resolved = posix.normalize(posix.join(sourceDir, rawPath));

    const slug = sourceToSlug.get(resolved);
    if (!slug) {
      return match; // unknown target — leave unchanged
    }

    const suffix = format === 'md' ? '.md' : '';
    const siteUrl = slug === 'index' ? `/${suffix ? 'index.md' : ''}` : `/${slug}${suffix}`;
    return `[${text}](${siteUrl}${fragment})`;
  });
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
node --test lib/rewrite-links.test.js
```

Expected: 7 tests pass.

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/lib/rewrite-links.js tools/site-generator/lib/rewrite-links.test.js
git commit -m "feat(site): add link rewriting for spec-relative to site URLs"
```

---

### Task 4: Mermaid Rendering

**Files:**
- Create: `tools/site-generator/lib/render-mermaid.js`
- Create: `tools/site-generator/lib/render-mermaid.test.js`

- [ ] **Step 1: Write the failing test**

Create `tools/site-generator/lib/render-mermaid.test.js`:

```js
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { renderMermaidBlocks } from './render-mermaid.js';

describe('renderMermaidBlocks', () => {
  it('replaces mermaid code blocks with SVG', async () => {
    const md = [
      'Some text.',
      '',
      '```mermaid',
      'graph LR',
      '    A --> B',
      '```',
      '',
      'More text.',
    ].join('\n');

    const result = await renderMermaidBlocks(md);

    // Mermaid block should be replaced with an SVG
    assert.ok(!result.includes('```mermaid'), 'mermaid code block should be removed');
    assert.ok(result.includes('<svg'), 'should contain inline SVG');
    assert.ok(result.includes('Some text.'), 'surrounding text preserved');
    assert.ok(result.includes('More text.'), 'surrounding text preserved');
  });

  it('returns markdown unchanged when no mermaid blocks', async () => {
    const md = 'Just regular markdown.\n\n```js\nconsole.log("hi");\n```';
    const result = await renderMermaidBlocks(md);
    assert.equal(result, md);
  });

  it('handles multiple mermaid blocks', async () => {
    const md = [
      '```mermaid',
      'graph LR',
      '    A --> B',
      '```',
      '',
      'Between.',
      '',
      '```mermaid',
      'graph TD',
      '    C --> D',
      '```',
    ].join('\n');

    const result = await renderMermaidBlocks(md);
    const svgCount = (result.match(/<svg/g) || []).length;
    assert.equal(svgCount, 2, 'should have 2 SVGs');
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
node --test lib/render-mermaid.test.js
```

Expected: FAIL — `Cannot find module './render-mermaid.js'`

- [ ] **Step 3: Write implementation**

Create `tools/site-generator/lib/render-mermaid.js`:

```js
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

  // Find mmdc binary from node_modules
  const mmdc = new URL(
    '../node_modules/.bin/mmdc',
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
      await execFileAsync(mmdc, [
        '-i', inputFile,
        '-o', outputFile,
        '-b', 'transparent',
        '--quiet',
      ]);

      const svg = await readFile(outputFile, 'utf-8');
      // Replace the full code block (``` to ```) with the SVG
      result =
        result.slice(0, block.index) +
        svg.trim() +
        result.slice(block.index + block[0].length);
    }

    return result;
  } finally {
    await rm(workDir, { recursive: true, force: true });
  }
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
node --test lib/render-mermaid.test.js
```

Expected: 3 tests pass. Note: first run may be slow (~5-10s) as Puppeteer launches a browser to render diagrams.

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/lib/render-mermaid.js tools/site-generator/lib/render-mermaid.test.js
git commit -m "feat(site): add mermaid diagram pre-rendering to SVG"
```

---

### Task 5: HTML Template

**Files:**
- Create: `tools/site-generator/template.html`
- Create: `tools/site-generator/assets/style.css`

- [ ] **Step 1: Create the HTML template**

Create `tools/site-generator/template.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{title}} - SpecScore</title>
  <link rel="stylesheet" href="/assets/style.css">
  <script>
    (function() {
      var saved = localStorage.getItem('specscore-theme');
      if (saved) {
        document.documentElement.setAttribute('data-theme', saved);
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.setAttribute('data-theme', 'dark');
      }
    })();
  </script>
</head>
<body>
  <header>
    <div class="header-inner">
      <a href="/" class="wordmark">SpecScore</a>
      <nav>
        {{nav}}
        <a href="https://github.com/synchestra-io/specscore">GitHub</a>
      </nav>
      <button id="theme-toggle" type="button" aria-label="Toggle dark/light theme">
        <span class="theme-icon-light">&#9789;</span>
        <span class="theme-icon-dark">&#9788;</span>
      </button>
    </div>
  </header>
  <main>
    {{content}}
    <p class="view-markdown"><a href="{{mdUrl}}">View as Markdown</a></p>
  </main>
  <footer>
    <div class="footer-inner">
      <span>Apache License 2.0</span>
      <a href="https://github.com/synchestra-io/specscore">GitHub</a>
    </div>
  </footer>
  <script>
    document.getElementById('theme-toggle').addEventListener('click', function() {
      var current = document.documentElement.getAttribute('data-theme');
      var next = current === 'dark' ? 'light' : 'dark';
      document.documentElement.setAttribute('data-theme', next);
      localStorage.setItem('specscore-theme', next);
    });
  </script>
</body>
</html>
```

- [ ] **Step 2: Create the CSS file**

Create `tools/site-generator/assets/style.css`:

```css
:root {
  --bg: #ffffff;
  --bg-secondary: #f6f8fa;
  --text: #1f2328;
  --text-secondary: #656d76;
  --border: #d0d7de;
  --link: #0969da;
  --code-bg: #f6f8fa;
  --header-bg: #f6f8fa;
  --header-border: #d0d7de;
  --table-stripe: #f6f8fa;
}

[data-theme="dark"] {
  --bg: #0d1117;
  --bg-secondary: #161b22;
  --text: #e6edf3;
  --text-secondary: #8b949e;
  --border: #30363d;
  --link: #58a6ff;
  --code-bg: #161b22;
  --header-bg: #161b22;
  --header-border: #30363d;
  --table-stripe: #161b22;
}

* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: var(--text);
  background: var(--bg);
}

/* Header */
header {
  background: var(--header-bg);
  border-bottom: 1px solid var(--header-border);
  padding: 0.75rem 1.5rem;
}

.header-inner {
  max-width: 52rem;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.wordmark {
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--text);
  text-decoration: none;
}

nav {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  flex: 1;
}

nav a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
}

nav a:hover,
nav a.active {
  color: var(--link);
}

#theme-toggle {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  color: var(--text);
  font-size: 1rem;
}

.theme-icon-dark { display: none; }
[data-theme="dark"] .theme-icon-light { display: none; }
[data-theme="dark"] .theme-icon-dark { display: inline; }

/* Main content */
main {
  max-width: 52rem;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

/* Typography */
h1 { font-size: 2rem; margin: 0 0 1rem; border-bottom: 1px solid var(--border); padding-bottom: 0.5rem; }
h2 { font-size: 1.5rem; margin: 2rem 0 0.75rem; border-bottom: 1px solid var(--border); padding-bottom: 0.25rem; }
h3 { font-size: 1.25rem; margin: 1.5rem 0 0.5rem; }
h4 { font-size: 1rem; margin: 1.25rem 0 0.5rem; }

p { margin: 0.75rem 0; }

a { color: var(--link); text-decoration: none; }
a:hover { text-decoration: underline; }

/* Code */
code {
  background: var(--code-bg);
  padding: 0.15rem 0.35rem;
  border-radius: 4px;
  font-size: 0.875em;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
}

pre {
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 1rem;
  overflow-x: auto;
  margin: 1rem 0;
}

pre code {
  background: none;
  padding: 0;
  border-radius: 0;
}

/* Tables */
table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
  font-size: 0.875rem;
}

th, td {
  border: 1px solid var(--border);
  padding: 0.5rem 0.75rem;
  text-align: left;
}

th { background: var(--bg-secondary); font-weight: 600; }
tr:nth-child(even) { background: var(--table-stripe); }

/* Lists */
ul, ol { padding-left: 1.5rem; margin: 0.5rem 0; }
li { margin: 0.25rem 0; }

/* Mermaid SVGs */
svg { max-width: 100%; height: auto; margin: 1rem 0; }

/* View as Markdown link */
.view-markdown {
  margin-top: 3rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
  font-size: 0.875rem;
  color: var(--text-secondary);
}

/* Footer */
footer {
  border-top: 1px solid var(--border);
  padding: 1.5rem;
  margin-top: 3rem;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.footer-inner {
  max-width: 52rem;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
}

/* Responsive */
@media (max-width: 640px) {
  .header-inner { flex-wrap: wrap; }
  nav { order: 3; width: 100%; }
  main { padding: 1rem; }
}
```

- [ ] **Step 3: Verify template renders locally**

Quick manual check — open `tools/site-generator/template.html` in a browser. It will show the placeholder tokens (`{{title}}`, `{{nav}}`, etc.) but you can verify the layout, theme toggle, and CSS load correctly. The theme toggle button should switch between light and dark when clicked.

- [ ] **Step 4: Commit**

```bash
git add tools/site-generator/template.html tools/site-generator/assets/style.css
git commit -m "feat(site): add HTML template and CSS with dark/light theme"
```

---

### Task 6: Page Renderer

**Files:**
- Create: `tools/site-generator/lib/render-page.js`
- Create: `tools/site-generator/lib/render-page.test.js`

- [ ] **Step 1: Write the failing test**

Create `tools/site-generator/lib/render-page.test.js`:

```js
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { renderMarkdownToHtml, injectIntoTemplate } from './render-page.js';

describe('renderMarkdownToHtml', () => {
  it('renders basic markdown to HTML', () => {
    const html = renderMarkdownToHtml('# Hello\n\nA paragraph.');
    assert.ok(html.includes('<h1>Hello</h1>'));
    assert.ok(html.includes('<p>A paragraph.</p>'));
  });

  it('renders GFM tables', () => {
    const md = '| A | B |\n|---|---|\n| 1 | 2 |';
    const html = renderMarkdownToHtml(md);
    assert.ok(html.includes('<table>'));
    assert.ok(html.includes('<td>1</td>'));
  });

  it('renders fenced code blocks with language class', () => {
    const md = '```js\nconsole.log("hi");\n```';
    const html = renderMarkdownToHtml(md);
    assert.ok(html.includes('<code class="language-js">'));
  });
});

describe('injectIntoTemplate', () => {
  const template = '<title>{{title}} - SpecScore</title>{{nav}}<main>{{content}}</main><a href="{{mdUrl}}">md</a>';

  const navItems = [
    { slug: 'index', navLabel: 'Home' },
    { slug: 'specifications', navLabel: 'Specifications' },
  ];

  it('replaces all placeholders', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '<h1>Feature</h1>',
      slug: 'feature-specification',
      navItems,
    });

    assert.ok(result.includes('<title>Feature Specification - SpecScore</title>'));
    assert.ok(result.includes('<h1>Feature</h1>'));
    assert.ok(result.includes('href="/feature-specification.md"'));
    assert.ok(result.includes('href="/"'));
    assert.ok(result.includes('href="/specifications"'));
  });

  it('marks the current page as active in nav', () => {
    const result = injectIntoTemplate(template, {
      title: 'Specifications',
      content: '<p>List</p>',
      slug: 'specifications',
      navItems,
    });

    assert.ok(result.includes('class="active"'));
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
node --test lib/render-page.test.js
```

Expected: FAIL — `Cannot find module './render-page.js'`

- [ ] **Step 3: Write implementation**

Create `tools/site-generator/lib/render-page.js`:

```js
import MarkdownIt from 'markdown-it';

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: false,
});

/**
 * Renders markdown string to HTML using markdown-it.
 * @param {string} markdown
 * @returns {string} HTML
 */
export function renderMarkdownToHtml(markdown) {
  return md.render(markdown);
}

/**
 * Injects rendered content into the HTML template.
 * @param {string} template - HTML template with {{placeholders}}
 * @param {{ title: string, content: string, slug: string, navItems: Array }} opts
 * @returns {string} complete HTML page
 */
export function injectIntoTemplate(template, { title, content, slug, navItems }) {
  const navHtml = navItems
    .map((item) => {
      const href = item.slug === 'index' ? '/' : `/${item.slug}`;
      const cls = item.slug === slug ? ' class="active"' : '';
      return `<a href="${href}"${cls}>${item.navLabel}</a>`;
    })
    .join('\n        ');

  const mdUrl = slug === 'index' ? '/index.md' : `/${slug}.md`;

  return template
    .replace(/\{\{title\}\}/g, title)
    .replace(/\{\{nav\}\}/g, navHtml)
    .replace(/\{\{content\}\}/g, content)
    .replace(/\{\{mdUrl\}\}/g, mdUrl);
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
node --test lib/render-page.test.js
```

Expected: 5 tests pass.

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/lib/render-page.js tools/site-generator/lib/render-page.test.js
git commit -m "feat(site): add markdown renderer and template injection"
```

---

### Task 7: Build Script

**Files:**
- Create: `tools/site-generator/build.js`

- [ ] **Step 1: Write the build script**

Create `tools/site-generator/build.js`:

```js
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
```

- [ ] **Step 2: Run the build**

Run from `tools/site-generator/`:

```bash
npm run build
```

Expected output:

```
Building 7 pages...
  index.html + index.md
  specifications.html + specifications.md
  feature-specification.html + feature-specification.md
  development-plan-specification.html + development-plan-specification.md
  acceptance-criteria-specification.html + acceptance-criteria-specification.md
  project-definition-specification.html + project-definition-specification.md
  source-references-specification.html + source-references-specification.md

Done. Output: /path/to/specscore/public
```

- [ ] **Step 3: Verify output**

```bash
ls ../../public/
```

Expected: `index.html`, `index.md`, `specifications.html`, `specifications.md`, `feature-specification.html`, `feature-specification.md`, etc., plus `assets/style.css`.

- [ ] **Step 4: Manual spot-check**

Open `public/index.html` in a browser. Verify:
- Page renders with header, nav, content, footer
- Nav links work (relative to `/`)
- Theme toggle switches colors
- "View as Markdown" link points to `/index.md`

Open `public/feature-specification.html`. Verify:
- Mermaid diagrams render as SVG images (not code blocks)
- Internal links point to other site pages (not `../development-plan/README.md`)

Open `public/feature-specification.md`. Verify:
- Links point to `.md` siblings (e.g., `/development-plan-specification.md`)
- Mermaid blocks are still in markdown format (not SVG)

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/build.js
git commit -m "feat(site): add build script orchestrating full site generation"
```

---

### Task 8: Run All Tests

**Files:** (none — verification only)

- [ ] **Step 1: Run all tests**

```bash
node --test lib/load-config.test.js lib/rewrite-links.test.js lib/render-mermaid.test.js lib/render-page.test.js
```

Expected: All 18 tests pass.

- [ ] **Step 2: Run a clean build**

```bash
npm run clean && npm run build
```

Expected: `public/` is deleted and rebuilt. All 7 pages generated.

- [ ] **Step 3: Verify no mermaid code blocks in HTML output**

```bash
grep -r '```mermaid' ../../public/*.html; echo "exit: $?"
```

Expected: No matches. Exit code 1 (grep found nothing).

- [ ] **Step 4: Verify all HTML files have SVG diagrams for pages that should have them**

```bash
grep -l '<svg' ../../public/feature-specification.html ../../public/development-plan-specification.html
```

Expected: Both files listed (these are the two specs with mermaid diagrams).

- [ ] **Step 5: Verify .md files do NOT contain SVG**

```bash
grep -r '<svg' ../../public/*.md; echo "exit: $?"
```

Expected: No matches. The `.md` files should contain original mermaid code blocks, not SVGs.

- [ ] **Step 6: Commit (if any fixes were needed)**

Only if previous steps required fixes. Otherwise skip.

---

### Task 9: Add npm test script

**Files:**
- Modify: `tools/site-generator/package.json`

- [ ] **Step 1: Add test script to `package.json`**

Add to the `"scripts"` section:

```json
"test": "node --test lib/**/*.test.js"
```

So the scripts section becomes:

```json
"scripts": {
  "build": "node build.js",
  "clean": "rm -rf ../../public",
  "test": "node --test lib/**/*.test.js"
}
```

- [ ] **Step 2: Verify `npm test` works**

```bash
npm test
```

Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add tools/site-generator/package.json
git commit -m "feat(site): add npm test script"
```
