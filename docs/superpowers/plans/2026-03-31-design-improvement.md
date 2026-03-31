# specscore.org Design Improvement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign specscore.org with a "Clarity & Standards" aesthetic — Source Serif 4 headings, indigo accent, sidebar navigation, and a dedicated landing page.

**Architecture:** Update the existing static site generator's config, template, CSS, and build script. No new dependencies. The sidebar structure is driven by a new `navGroup` field in `site-config.json`. The landing page is a hand-crafted HTML file injected through the same template.

**Tech Stack:** Node.js, markdown-it, Google Fonts (Source Serif 4), CSS custom properties.

---

## File Map

**Modified files:**
- `tools/site-generator/site-config.json` — add `navGroup` field, set index to `nav: false`
- `tools/site-generator/lib/load-config.js` — replace `navItems` with `sidebarGroups`
- `tools/site-generator/lib/load-config.test.js` — update tests for new return shape
- `tools/site-generator/lib/render-page.js` — update `injectIntoTemplate` for sidebar, eyebrow, view-markdown
- `tools/site-generator/lib/render-page.test.js` — update tests for new signature
- `tools/site-generator/lib/render-mermaid.js` — wrap SVG output in `.mermaid-wrap` container
- `tools/site-generator/lib/render-mermaid.test.js` — add wrapper test
- `tools/site-generator/template.html` — new layout with sidebar + header-right
- `tools/site-generator/assets/style.css` — full replacement with new design system
- `tools/site-generator/build.js` — use `sidebarGroups`, handle landing page, pass eyebrow

**New files:**
- `tools/site-generator/landing.html` — landing page content (hero + spec list)

---

### Task 1: Update site-config.json with navGroup fields

**Files:**
- Modify: `tools/site-generator/site-config.json`

- [ ] **Step 1: Add navGroup to all entries and set index to nav:false**

Replace the full contents of `tools/site-generator/site-config.json` with:

```json
{
  "pages": [
    {
      "source": "spec/README.md",
      "slug": "index",
      "title": "SpecScore",
      "nav": false,
      "navLabel": "Home",
      "navOrder": 0
    },
    {
      "source": "spec/features/README.md",
      "slug": "specifications",
      "title": "Specifications",
      "nav": true,
      "navLabel": "Specifications",
      "navGroup": "Getting Started",
      "navOrder": 1
    },
    {
      "source": "spec/features/feature/README.md",
      "slug": "feature-specification",
      "title": "Feature Specification",
      "nav": true,
      "navLabel": "Feature",
      "navGroup": "Feature Specs",
      "navOrder": 2
    },
    {
      "source": "spec/features/development-plan/README.md",
      "slug": "development-plan-specification",
      "title": "Development Plan Specification",
      "nav": true,
      "navLabel": "Development Plan",
      "navGroup": "Feature Specs",
      "navOrder": 3
    },
    {
      "source": "spec/features/acceptance-criteria/README.md",
      "slug": "acceptance-criteria-specification",
      "title": "Acceptance Criteria Specification",
      "nav": true,
      "navLabel": "Acceptance Criteria",
      "navGroup": "Feature Specs",
      "navOrder": 4
    },
    {
      "source": "spec/features/project-definition/README.md",
      "slug": "project-definition-specification",
      "title": "Project Definition Specification",
      "nav": true,
      "navLabel": "Project Definition",
      "navGroup": "Feature Specs",
      "navOrder": 5
    },
    {
      "source": "spec/features/source-references/README.md",
      "slug": "source-references-specification",
      "title": "Source References Specification",
      "nav": true,
      "navLabel": "Source References",
      "navGroup": "Feature Specs",
      "navOrder": 6
    }
  ]
}
```

Changes: `index` entry now has `nav: false` (removed from sidebar). All other entries gain a `navGroup` field — `"Getting Started"` for Specifications, `"Feature Specs"` for the five feature spec pages.

- [ ] **Step 2: Commit**

```bash
cd tools/site-generator && git add site-config.json && git commit -m "feat(site): add navGroup fields to site-config.json for sidebar grouping"
```

---

### Task 2: Update loadConfig to return sidebarGroups (TDD)

**Files:**
- Test: `tools/site-generator/lib/load-config.test.js`
- Modify: `tools/site-generator/lib/load-config.js`

- [ ] **Step 1: Write failing tests for sidebarGroups**

Replace the full contents of `tools/site-generator/lib/load-config.test.js` with:

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

  it('builds sidebar groups from navGroup fields', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(Array.isArray(config.sidebarGroups));
    assert.ok(config.sidebarGroups.length > 0);
    for (const group of config.sidebarGroups) {
      assert.ok(typeof group.label === 'string');
      assert.ok(Array.isArray(group.items));
      assert.ok(group.items.length > 0);
    }
  });

  it('sidebar groups contain Feature Specs entries', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    const featureGroup = config.sidebarGroups.find((g) => g.label === 'Feature Specs');
    assert.ok(featureGroup, 'Feature Specs group should exist');
    assert.ok(featureGroup.items.some((p) => p.slug === 'feature-specification'));
  });

  it('excludes nav:false pages from sidebar groups', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    const allItems = config.sidebarGroups.flatMap((g) => g.items);
    assert.ok(!allItems.some((p) => p.slug === 'index'), 'index should not appear in sidebar');
  });
});
```

- [ ] **Step 2: Run tests to verify failure**

Run: `cd tools/site-generator && node --test lib/load-config.test.js`

Expected: 3 new tests FAIL — `config.sidebarGroups` is `undefined`.

- [ ] **Step 3: Update load-config.js to return sidebarGroups**

Replace the full contents of `tools/site-generator/lib/load-config.js` with:

```js
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd tools/site-generator && node --test lib/load-config.test.js`

Expected: All 5 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add lib/load-config.js lib/load-config.test.js && git commit -m "feat(site): update loadConfig to return sidebarGroups for sidebar navigation"
```

---

### Task 3: Update injectIntoTemplate for sidebar, eyebrow, and view-markdown (TDD)

**Files:**
- Test: `tools/site-generator/lib/render-page.test.js`
- Modify: `tools/site-generator/lib/render-page.js`

- [ ] **Step 1: Write failing tests for new injectIntoTemplate signature**

Replace the full contents of `tools/site-generator/lib/render-page.test.js` with:

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
  const template =
    '<title>{{title}} - SpecScore</title>' +
    '<aside>{{sidebar}}</aside>' +
    '{{eyebrow}}' +
    '<main>{{content}}</main>' +
    '{{viewMarkdown}}';

  const sidebarGroups = [
    {
      label: 'Getting Started',
      items: [
        { slug: 'specifications', navLabel: 'Specifications', navOrder: 1 },
      ],
    },
    {
      label: 'Feature Specs',
      items: [
        { slug: 'feature-specification', navLabel: 'Feature', navOrder: 2 },
        { slug: 'source-references-specification', navLabel: 'Source References', navOrder: 6 },
      ],
    },
  ];

  it('replaces title and content placeholders', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '<h1>Feature</h1>',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('<title>Feature Specification - SpecScore</title>'));
    assert.ok(result.includes('<h1>Feature</h1>'));
  });

  it('generates sidebar HTML with group labels and links', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('Getting Started'));
    assert.ok(result.includes('Feature Specs'));
    assert.ok(result.includes('href="/specifications"'));
    assert.ok(result.includes('href="/feature-specification"'));
  });

  it('marks the current page as active in sidebar', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('class="active"'));
    // The active link should be the feature-specification one
    assert.match(result, /href="\/feature-specification"[^>]*class="active"/);
  });

  it('renders eyebrow label when provided', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('class="page-eyebrow"'));
    assert.ok(result.includes('Feature Specs'));
  });

  it('omits eyebrow div when eyebrow is empty', () => {
    const result = injectIntoTemplate(template, {
      title: 'SpecScore',
      content: '',
      slug: 'index',
      sidebarGroups,
      eyebrow: '',
    });

    assert.ok(!result.includes('class="page-eyebrow"'));
  });

  it('includes view-markdown link by default', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('href="/feature-specification.md"'));
    assert.ok(result.includes('View as Markdown'));
  });

  it('omits view-markdown link when showViewMarkdown is false', () => {
    const result = injectIntoTemplate(template, {
      title: 'SpecScore',
      content: '',
      slug: 'index',
      sidebarGroups,
      eyebrow: '',
      showViewMarkdown: false,
    });

    assert.ok(!result.includes('View as Markdown'));
  });
});
```

- [ ] **Step 2: Run tests to verify failure**

Run: `cd tools/site-generator && node --test lib/render-page.test.js`

Expected: New `injectIntoTemplate` tests FAIL because the function still expects the old `navItems` parameter and the template doesn't have `{{sidebar}}`, `{{eyebrow}}`, `{{viewMarkdown}}` placeholders.

- [ ] **Step 3: Update render-page.js with new injectIntoTemplate**

Replace the full contents of `tools/site-generator/lib/render-page.js` with:

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
 * Builds sidebar HTML from grouped navigation items.
 * @param {Array<{label: string, items: Array}>} sidebarGroups
 * @param {string} currentSlug
 * @returns {string} HTML
 */
function buildSidebarHtml(sidebarGroups, currentSlug) {
  return sidebarGroups
    .map((group) => {
      const links = group.items
        .map((item) => {
          const href = item.slug === 'index' ? '/' : `/${item.slug}`;
          const active = item.slug === currentSlug ? ' class="active"' : '';
          return `<li><a href="${href}"${active}>${item.navLabel}</a></li>`;
        })
        .join('\n          ');
      return `<div class="sidebar-section">
        <span class="sidebar-label">${group.label}</span>
        <ul class="sidebar-nav">
          ${links}
        </ul>
      </div>`;
    })
    .join('\n      ');
}

/**
 * Injects rendered content into the HTML template.
 * @param {string} template - HTML template with {{placeholders}}
 * @param {object} opts
 * @param {string} opts.title - page title
 * @param {string} opts.content - rendered HTML content
 * @param {string} opts.slug - page slug for active state
 * @param {Array} opts.sidebarGroups - grouped nav items from loadConfig
 * @param {string} opts.eyebrow - eyebrow label (navGroup of this page)
 * @param {boolean} [opts.showViewMarkdown=true] - whether to show the "View as Markdown" link
 * @returns {string} complete HTML page
 */
export function injectIntoTemplate(template, { title, content, slug, sidebarGroups, eyebrow, showViewMarkdown = true }) {
  const sidebarHtml = buildSidebarHtml(sidebarGroups, slug);
  const mdUrl = slug === 'index' ? '/index.md' : `/${slug}.md`;
  const eyebrowHtml = eyebrow
    ? `<div class="page-eyebrow">${eyebrow}</div>`
    : '';
  const viewMarkdownHtml = showViewMarkdown
    ? `<p class="view-markdown"><a href="${mdUrl}">View as Markdown</a></p>`
    : '';

  return template
    .replace(/\{\{title\}\}/g, title)
    .replace(/\{\{sidebar\}\}/g, sidebarHtml)
    .replace(/\{\{eyebrow\}\}/g, eyebrowHtml)
    .replace(/\{\{content\}\}/g, content)
    .replace(/\{\{viewMarkdown\}\}/g, viewMarkdownHtml);
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd tools/site-generator && node --test lib/render-page.test.js`

Expected: All 10 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add lib/render-page.js lib/render-page.test.js && git commit -m "feat(site): update injectIntoTemplate for sidebar, eyebrow, and view-markdown"
```

---

### Task 4: Update renderMermaidBlocks to wrap SVG in container (TDD)

**Files:**
- Test: `tools/site-generator/lib/render-mermaid.test.js`
- Modify: `tools/site-generator/lib/render-mermaid.js`

- [ ] **Step 1: Add failing test for mermaid wrapper**

In `tools/site-generator/lib/render-mermaid.test.js`, add the following test inside the existing `describe('renderMermaidBlocks', ...)` block, after the last `it(...)`:

```js
  it('wraps SVG in mermaid-wrap container', { skip: !mmdcAvailable && 'mmdc/Chrome not available' }, async () => {
    const md = ['```mermaid', 'graph LR', '    A --> B', '```'].join('\n');
    const result = await renderMermaidBlocks(md);
    assert.ok(result.includes('class="mermaid-wrap"'), 'should have mermaid-wrap wrapper');
    assert.ok(result.includes('class="mermaid-label"'), 'should have mermaid-label');
    assert.ok(result.includes('class="mermaid-body"'), 'should have mermaid-body');
  });
```

- [ ] **Step 2: Run tests to verify failure**

Run: `cd tools/site-generator && node --test lib/render-mermaid.test.js`

Expected: New test FAILS — SVG output does not contain `mermaid-wrap`. (If mmdc is not available, the test is skipped — this is fine, CI will catch it.)

- [ ] **Step 3: Update render-mermaid.js to wrap SVG**

In `tools/site-generator/lib/render-mermaid.js`, find the block inside the `for` loop that reads the SVG and replaces the mermaid code block:

```js
        const svg = await readFile(outputFile, 'utf-8');
        // Replace the full code block (``` to ```) with the SVG
        result =
          result.slice(0, block.index) +
          svg.trim() +
          result.slice(block.index + block[0].length);
```

Replace it with:

```js
        const svg = await readFile(outputFile, 'utf-8');
        const wrapped =
          '<div class="mermaid-wrap">\n' +
          '<div class="mermaid-label">Diagram</div>\n' +
          '<div class="mermaid-body">' + svg.trim() + '</div>\n' +
          '</div>';
        result =
          result.slice(0, block.index) +
          wrapped +
          result.slice(block.index + block[0].length);
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd tools/site-generator && node --test lib/render-mermaid.test.js`

Expected: All tests PASS (or skipped if mmdc unavailable).

- [ ] **Step 5: Commit**

```bash
git add lib/render-mermaid.js lib/render-mermaid.test.js && git commit -m "feat(site): wrap mermaid SVGs in labeled container"
```

---

### Task 5: Replace template.html with sidebar layout

**Files:**
- Modify: `tools/site-generator/template.html`

- [ ] **Step 1: Write new template.html**

Replace the full contents of `tools/site-generator/template.html` with:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{title}} - SpecScore</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Source+Serif+4:ital,opsz,wght@0,8..60,400;0,8..60,600;0,8..60,700;1,8..60,400&display=swap" rel="stylesheet">
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
  <header class="site-header">
    <a href="/" class="wordmark">SpecScore</a>
    <div class="header-right">
      <a href="https://github.com/synchestra-io/specscore">GitHub</a>
      <button id="theme-toggle" type="button" aria-label="Toggle dark/light theme">
        <span class="theme-icon-light">&#9789;</span>
        <span class="theme-icon-dark">&#9788;</span>
      </button>
    </div>
  </header>
  <div class="site-body">
    <aside class="sidebar">
      {{sidebar}}
      <div class="sidebar-section">
        <span class="sidebar-label">Resources</span>
        <ul class="sidebar-nav">
          <li><a href="https://github.com/synchestra-io/specscore">GitHub ↗</a></li>
        </ul>
      </div>
    </aside>
    <main>
      {{eyebrow}}
      {{content}}
      {{viewMarkdown}}
    </main>
  </div>
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

Key changes from old template:
- Header: flat structure with wordmark on left, `.header-right` (GitHub link + theme toggle) on right. No nav in header.
- Body: `.site-body` flex container with `<aside class="sidebar">` and `<main>`.
- Sidebar: `{{sidebar}}` placeholder for dynamically-generated groups + hardcoded Resources section at bottom.
- Main: `{{eyebrow}}` before `{{content}}`, `{{viewMarkdown}}` after.
- Google Fonts: Source Serif 4 loaded via `<link>` with `preconnect`.
- Old `{{nav}}` and `{{mdUrl}}` placeholders are gone.

- [ ] **Step 2: Commit**

```bash
git add template.html && git commit -m "feat(site): replace template with sidebar layout and new placeholders"
```

---

### Task 6: Create landing page content

**Files:**
- Create: `tools/site-generator/landing.html`

- [ ] **Step 1: Write landing.html**

Create `tools/site-generator/landing.html` with:

```html
<div class="landing-hero">
  <h1>The specification framework for AI-driven development</h1>
  <p class="landing-tagline">SpecScore defines how software specifications are structured, validated, and linked to the code they govern — so AI agents and humans work from the same source of truth.</p>
  <div class="landing-cta">
    <a href="/specifications" class="btn-primary">View Specifications</a>
    <a href="https://github.com/synchestra-io/specscore" class="btn-secondary">GitHub ↗</a>
  </div>
</div>
<div class="landing-specs">
  <div class="landing-section-label">Specifications</div>
  <ul class="spec-list">
    <li>
      <a href="/feature-specification">Feature</a>
      <span>Structure and metadata for a single feature specification</span>
    </li>
    <li>
      <a href="/acceptance-criteria-specification">Acceptance Criteria</a>
      <span>Machine-verifiable conditions that define feature completeness</span>
    </li>
    <li>
      <a href="/development-plan-specification">Development Plan</a>
      <span>Ordered implementation steps linked to features and criteria</span>
    </li>
    <li>
      <a href="/source-references-specification">Source References</a>
      <span>Annotations that link source code back to the governing spec</span>
    </li>
    <li>
      <a href="/project-definition-specification">Project Definition</a>
      <span>Project-level configuration, structure, and metadata</span>
    </li>
  </ul>
</div>
```

This is a content fragment — it gets injected into the `{{content}}` slot of `template.html` by `build.js`. Not a full HTML document.

- [ ] **Step 2: Commit**

```bash
git add landing.html && git commit -m "feat(site): add landing page content with hero and spec list"
```

---

### Task 7: Update build.js for sidebarGroups and landing page

**Files:**
- Modify: `tools/site-generator/build.js`

- [ ] **Step 1: Update build.js**

Replace the full contents of `tools/site-generator/build.js` with:

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

Changes from old build.js:
- `config.sidebarGroups` used instead of `config.navItems`.
- `index` slug triggers landing page path: reads `landing.html`, passes `showViewMarkdown: false`, no `.md` file generated.
- All other pages pass `eyebrow: page.navGroup || ''` and use `sidebarGroups`.

- [ ] **Step 2: Run unit tests**

Run: `cd tools/site-generator && npm test`

Expected: All tests PASS. (The build script itself is not unit-tested; it is integration-tested via `npm run build` in the next step.)

- [ ] **Step 3: Run build**

Run: `cd tools/site-generator && npm run build`

Expected: Builds all pages. Output shows `index.html (landing)` for the home page and `{slug}.html + {slug}.md` for all other pages. No errors.

- [ ] **Step 4: Verify output structure**

Run: `ls -la ../../public/`

Expected: `index.html` (no `index.md`), plus `specifications.html`, `specifications.md`, `feature-specification.html`, `feature-specification.md`, etc., and `assets/style.css`.

- [ ] **Step 5: Commit**

```bash
git add build.js && git commit -m "feat(site): update build script for sidebar groups and landing page"
```

---

### Task 8: Replace style.css with new design system

**Files:**
- Modify: `tools/site-generator/assets/style.css`

- [ ] **Step 1: Write new style.css**

Replace the full contents of `tools/site-generator/assets/style.css` with:

```css
/* === COLOR SYSTEM === */
:root {
  --bg: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-sidebar: #f3f4f6;
  --text: #111827;
  --text-secondary: #6b7280;
  --text-muted: #9ca3af;
  --border: #e5e7eb;
  --border-subtle: #f3f4f6;
  --accent: #4f46e5;
  --accent-light: #eef2ff;
  --code-bg: #f3f4f6;
  --code-border: #e5e7eb;
  --link: #4f46e5;
}

[data-theme="dark"] {
  --bg: #18181b;
  --bg-secondary: #27272a;
  --bg-sidebar: #1f1f23;
  --text: #f4f4f5;
  --text-secondary: #a1a1aa;
  --text-muted: #71717a;
  --border: #3f3f46;
  --border-subtle: #27272a;
  --accent: #818cf8;
  --accent-light: #1e1b4b;
  --code-bg: #27272a;
  --code-border: #3f3f46;
  --link: #818cf8;
}

/* === RESET === */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.65;
  color: var(--text);
  background: var(--bg);
}

/* === HEADER === */
.site-header {
  background: var(--bg);
  border-bottom: 1px solid var(--border);
  height: 56px;
  padding: 0 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  z-index: 10;
}

.wordmark {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text);
  text-decoration: none;
  letter-spacing: -0.01em;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}

.header-right a {
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-decoration: none;
}

.header-right a:hover { color: var(--accent); }

#theme-toggle {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 1rem;
  line-height: 1;
}

.theme-icon-dark { display: none; }
[data-theme="dark"] .theme-icon-light { display: none; }
[data-theme="dark"] .theme-icon-dark { display: inline; }

/* === LAYOUT === */
.site-body {
  display: flex;
  min-height: calc(100vh - 56px);
}

/* === SIDEBAR === */
.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border);
  padding: 2rem 0;
  position: sticky;
  top: 56px;
  height: calc(100vh - 56px);
  overflow-y: auto;
}

.sidebar-section {
  margin-bottom: 1.5rem;
  padding: 0 1.25rem;
}

.sidebar-label {
  display: block;
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.sidebar-nav { list-style: none; }

.sidebar-nav li a {
  display: block;
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-decoration: none;
  padding: 0.3rem 0.6rem;
  border-radius: 6px;
  margin-bottom: 1px;
}

.sidebar-nav li a:hover {
  background: var(--border);
  color: var(--text);
}

.sidebar-nav li a.active {
  background: var(--accent-light);
  color: var(--accent);
  font-weight: 600;
}

/* === MAIN CONTENT === */
main {
  flex: 1;
  min-width: 0;
  padding: 3rem 4rem;
  max-width: 820px;
}

/* === EYEBROW === */
.page-eyebrow {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

/* === TYPOGRAPHY === */
h1 {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 2.25rem;
  font-weight: 700;
  color: var(--text);
  line-height: 1.2;
  letter-spacing: -0.02em;
  margin: 0 0 0.5rem;
  border: none;
  padding: 0;
}

/* Indigo accent rule after h1 */
main > h1::after,
.page-eyebrow + h1::after {
  content: '';
  display: block;
  width: 2.5rem;
  height: 3px;
  background: var(--accent);
  border-radius: 2px;
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
}

/* No accent rule on landing page h1 */
.landing-hero h1::after { display: none; }

h2 {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 1.375rem;
  font-weight: 600;
  color: var(--text);
  line-height: 1.3;
  letter-spacing: -0.01em;
  margin: 2.5rem 0 0.75rem;
  padding-bottom: 0.4rem;
  border-bottom: 1px solid var(--border);
}

h3 {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text);
  margin: 1.75rem 0 0.5rem;
}

h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text);
  margin: 1.25rem 0 0.4rem;
}

p { margin: 0.75rem 0; }

a { color: var(--link); text-decoration: none; }
a:hover { text-decoration: underline; }

/* === CODE === */
code {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 0.85em;
  background: var(--code-bg);
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  color: var(--text);
}

pre {
  background: var(--code-bg);
  border: 1px solid var(--code-border);
  border-left: 3px solid var(--accent);
  border-radius: 0 6px 6px 0;
  padding: 1rem 1.25rem;
  overflow-x: auto;
  margin: 1.25rem 0;
}

pre code {
  background: none;
  padding: 0;
  font-size: 0.875rem;
  line-height: 1.6;
  color: var(--text);
}

/* === TABLES === */
table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.25rem 0;
  font-size: 0.875rem;
}

th {
  background: var(--bg-secondary);
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  text-align: left;
  padding: 0.6rem 0.875rem;
  border-bottom: 2px solid var(--border);
  color: var(--text-secondary);
}

td {
  padding: 0.6rem 0.875rem;
  border-bottom: 1px solid var(--border-subtle);
  vertical-align: top;
}

tr:hover td { background: var(--bg-secondary); }

/* === LISTS === */
ul, ol { padding-left: 1.5rem; margin: 0.75rem 0; }
li { margin: 0.3rem 0; }

/* === MERMAID === */
.mermaid-wrap {
  margin: 1.5rem 0;
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

.mermaid-label {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-muted);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  padding: 0.4rem 0.875rem;
}

.mermaid-body {
  padding: 1.5rem;
  background: var(--bg);
  text-align: center;
}

.mermaid-body svg { max-width: 100%; height: auto; }

/* Standalone SVGs (backward compat for any un-wrapped diagrams) */
svg { max-width: 100%; height: auto; margin: 1rem 0; }

/* === VIEW AS MARKDOWN === */
.view-markdown {
  margin-top: 3rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
  font-size: 0.8rem;
  color: var(--text-muted);
}

/* === FOOTER === */
footer {
  background: var(--bg-secondary);
  border-top: 1px solid var(--border);
  padding: 1.5rem 2rem;
}

.footer-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.footer-inner a { color: var(--text-secondary); text-decoration: none; }
.footer-inner a:hover { color: var(--accent); }

/* === LANDING PAGE === */
.landing-hero {
  padding: 4rem 0 3rem;
  border-bottom: 1px solid var(--border);
  margin-bottom: 3rem;
}

.landing-hero h1 {
  font-size: 2.75rem;
  line-height: 1.15;
  letter-spacing: -0.025em;
  max-width: 520px;
  margin-bottom: 1rem;
}

.landing-tagline {
  font-size: 1.1rem;
  color: var(--text-secondary);
  max-width: 480px;
  line-height: 1.6;
  margin-bottom: 2rem;
}

.landing-cta {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  align-items: center;
}

.btn-primary {
  background: var(--accent);
  color: #fff;
  padding: 0.6rem 1.25rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  text-decoration: none;
  display: inline-block;
}

.btn-primary:hover { opacity: 0.9; text-decoration: none; }

.btn-secondary {
  color: var(--text-secondary);
  padding: 0.6rem 1.25rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  text-decoration: none;
  border: 1px solid var(--border);
  display: inline-block;
}

.btn-secondary:hover { border-color: var(--accent); color: var(--accent); text-decoration: none; }

.landing-section-label {
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--accent);
  margin-bottom: 1.25rem;
}

.spec-list { list-style: none; padding: 0; }

.spec-list li {
  border-bottom: 1px solid var(--border);
  padding: 0.875rem 0;
  display: flex;
  align-items: baseline;
  gap: 1rem;
}

.spec-list li:first-child { border-top: 1px solid var(--border); }

.spec-list li a {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text);
  text-decoration: none;
  flex-shrink: 0;
  min-width: 160px;
}

.spec-list li a:hover { color: var(--accent); }

.spec-list li span {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

/* === RESPONSIVE === */
@media (max-width: 768px) {
  .site-body { flex-direction: column; }
  .sidebar {
    width: 100%;
    height: auto;
    position: static;
    border-right: none;
    border-bottom: 1px solid var(--border);
    padding: 1rem 0;
  }
  main { padding: 1.5rem 1rem; }
  h1 { font-size: 1.75rem; }
  .landing-hero h1 { font-size: 2rem; }
  .site-header { padding: 0 1rem; }
  .spec-list li { flex-direction: column; gap: 0.25rem; }
  .spec-list li a { min-width: auto; }
}
```

- [ ] **Step 2: Rebuild site**

Run: `cd tools/site-generator && npm run build`

Expected: Build succeeds with no errors.

- [ ] **Step 3: Commit**

```bash
git add assets/style.css && git commit -m "feat(site): replace CSS with new design system — indigo accent, Source Serif 4, sidebar layout"
```

---

### Task 9: Final verification

- [ ] **Step 1: Run all unit tests**

Run: `cd tools/site-generator && npm test`

Expected: All tests PASS.

- [ ] **Step 2: Run full build**

Run: `cd tools/site-generator && npm run build`

Expected: All pages build. Console output shows `index.html (landing)` and `{slug}.html + {slug}.md` for each spec page.

- [ ] **Step 3: Verify output files**

Run: `ls ../../public/`

Expected files:
```
index.html
specifications.html
specifications.md
feature-specification.html
feature-specification.md
development-plan-specification.html
development-plan-specification.md
acceptance-criteria-specification.html
acceptance-criteria-specification.md
project-definition-specification.html
project-definition-specification.md
source-references-specification.html
source-references-specification.md
assets/
```

No `index.md` (landing page has no markdown sibling).

- [ ] **Step 4: Spot-check HTML output**

Run: `head -30 ../../public/feature-specification.html`

Verify:
- `<link>` to Google Fonts for Source Serif 4 is present
- `<aside class="sidebar">` is present with grouped nav
- `<div class="page-eyebrow">Feature Specs</div>` appears before the h1
- `<p class="view-markdown">` appears at the end of main

Run: `head -20 ../../public/index.html`

Verify:
- Contains `class="landing-hero"`
- Does NOT contain `class="page-eyebrow"`
- Does NOT contain `View as Markdown`
