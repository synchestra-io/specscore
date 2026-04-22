# SpecScore Marketing Website Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform the SpecScore website from a documentation-only static site into a marketing-ready website at specscore.md — with updated branding, a redesigned landing page, blog infrastructure, an ecosystem page, and initial launch content.

**Architecture:** The existing Node.js site generator (`tools/site-generator/`) converts Markdown to HTML via `markdown-it` and injects into `template.html`. We extend it with blog support (date-prefixed Markdown files rendered as posts with an index page) and new content pages, while updating all branding references from `specscore.org` to `specscore.md`. No new frameworks or dependencies — extend what exists.

**Tech Stack:** Node.js site generator, markdown-it, HTML/CSS, Firebase Hosting

**Spec:** `docs/superpowers/specs/2026-04-14-specscore-marketing-plan-design.md`

---

## Task 1: Update Domain References (specscore.org → specscore.md)

**Files:**
- Modify: `tools/site-generator/template.html`
- Modify: `tools/site-generator/assets/style.css` (no domain refs, but verify)
- Modify: `firebase.json`
- Modify: `.firebaserc`

- [x] **Step 1: Update template.html wordmark**

In `tools/site-generator/template.html`, change line 24:

```html
<!-- Before -->
<a href="/" class="wordmark"><span class="wordmark-clef">𝄞</span> SpecScore.org</a>

<!-- After -->
<a href="/" class="wordmark"><span class="wordmark-clef">𝄞</span> SpecScore.md</a>
```

- [x] **Step 2: Update template.html footer**

In `tools/site-generator/template.html`, change the footer link (line 47):

```html
<!-- Before -->
<a href="https://specscore.org">specscore.org</a>

<!-- After -->
<a href="https://specscore.md">specscore.md</a>
```

- [x] **Step 3: Update Firebase hosting site name**

In `firebase.json`, update the site identifier:

```json
{
  "hosting": {
    "site": "specscore-md",
    "public": "public",
    "ignore": ["firebase.json", "**/.*"],
    "cleanUrls": true,
    "trailingSlash": false
  }
}
```

**Note:** The actual Firebase site must be created in the Firebase console before deploying. This change prepares the config.

- [x] **Step 4: Build and verify no broken references**

```bash
cd tools/site-generator && pnpm build
```

Verify: `grep -r "specscore.org" public/` returns no results. `grep -r "specscore.md" public/` shows updated references in header and footer.

- [x] **Step 5: Commit**

```bash
git add tools/site-generator/template.html firebase.json
git commit -m "brand: update domain from specscore.org to specscore.md"
```

---

## Task 2: Update Landing Page Hero with Marketing Positioning

**Files:**
- Modify: `tools/site-generator/landing.html`

The current landing page has a good structure but needs updated copy to match the approved positioning (standard-first, "spec-driven development" as owned term).

- [x] **Step 1: Rewrite landing.html with marketing positioning**

Replace the entire contents of `tools/site-generator/landing.html`:

```html
<div class="landing-hero">
  <h1>The open specification standard for AI-driven development</h1>
  <p class="landing-tagline">SpecScore is an open specification format that makes requirements machine-readable without making them human-unreadable. Markdown and YAML — version-controlled, portable, no vendor lock-in. A linter catches ambiguity before your agents do.</p>
  <div class="landing-cta">
    <a href="/specifications" class="btn-primary">Read the Spec</a>
    <a href="/install" class="btn-secondary">Install the CLI</a>
  </div>
</div>

<div class="landing-specs">
  <div class="landing-section-label">The Problem</div>
  <p class="landing-prose">AI agents work from specifications. Most specifications are ambiguous, scattered across Jira tickets and Notion docs, and impossible to validate automatically. When agents guess what you meant, you get rework. When specs live in chat logs, you get context loss.</p>
</div>

<div class="landing-specs">
  <div class="landing-section-label">The Format</div>
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
      <a href="/plan-specification">Plan</a>
      <span>Structured task breakdowns that bridge specs to execution</span>
    </li>
    <li>
      <a href="/task-specification">Task</a>
      <span>Discrete units of work with status lifecycle and dependencies</span>
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

<div class="landing-specs">
  <div class="landing-section-label">The Tooling</div>
  <p class="landing-prose">Write specs. Validate them. Ship them.</p>
  <ul class="spec-list">
    <li>
      <a href="/install">CLI</a>
      <span>Lint and validate specifications from the command line</span>
    </li>
    <li>
      <a>LSP</a>
      <span>Editor integration for real-time spec validation (coming soon)</span>
    </li>
  </ul>
</div>

<div class="landing-specs">
  <div class="landing-section-label">Find Your Path In</div>
  <ul class="spec-list">
    <li>
      <a href="/for/developers">Developers</a>
      <span>Clear specs, source references, structured features</span>
    </li>
    <li>
      <a href="/for/product-owners">Product Owners</a>
      <span>Standardized format for vision, goals, acceptance criteria</span>
    </li>
    <li>
      <a href="/for/qas">QAs</a>
      <span>Testable acceptance criteria and full feature context</span>
    </li>
    <li>
      <a href="/for/business-analysts">Business Analysts</a>
      <span>Formalize analysis into actionable, traceable specs</span>
    </li>
    <li>
      <a href="/for/project-managers">Project Managers</a>
      <span>Structured scope, task breakdowns, estimation visibility</span>
    </li>
    <li>
      <a href="/for/architects">Architects</a>
      <span>Technical decisions and system design documented consistently</span>
    </li>
  </ul>
</div>

<div class="landing-specs">
  <div class="landing-section-label">Ecosystem</div>
  <p class="landing-prose">SpecScore defines <em>what</em> gets built. When you're ready to test specs automatically or orchestrate agents across them, the ecosystem is there. But SpecScore works beautifully on its own.</p>
  <ul class="spec-list">
    <li>
      <a href="https://rehearse.ink" target="_blank" rel="noopener noreferrer">Rehearse</a>
      <span>Test your specifications automatically</span>
    </li>
    <li>
      <a href="https://synchestra.io" target="_blank" rel="noopener noreferrer">Synchestra</a>
      <span>Spec-driven orchestration for AI-assisted development</span>
    </li>
  </ul>
</div>
```

- [x] **Step 2: Add CSS for the new landing-prose class**

Add to `tools/site-generator/assets/style.css`, after the `.spec-list li span` rule (around line 512):

```css
/* === LANDING PROSE === */
.landing-prose {
  font-size: 1rem;
  color: var(--text-secondary);
  line-height: 1.65;
  max-width: 580px;
  margin-bottom: 1.25rem;
}
```

- [x] **Step 3: Build and visually verify**

```bash
cd tools/site-generator && pnpm build
```

Open `public/index.html` in a browser. Verify:
- Hero has updated tagline and "Read the Spec" as primary CTA
- "The Problem" section appears below hero
- "The Format" shows spec list
- "The Tooling" shows CLI and LSP
- "Find Your Path In" shows six role cards
- "Ecosystem" shows Rehearse and Synchestra links
- Responsive layout works on mobile (resize browser)

- [x] **Step 4: Commit**

```bash
git add tools/site-generator/landing.html tools/site-generator/assets/style.css
git commit -m "marketing: redesign landing page with standard-first positioning"
```

---

## Task 3: Add Ecosystem Page

**Files:**
- Create: `docs/ecosystem.md`
- Modify: `tools/site-generator/site-config.json`

- [x] **Step 1: Create ecosystem page content**

Create `docs/ecosystem.md`:

```markdown
# The Ecosystem

SpecScore is the specification layer in a broader ecosystem of tools for AI-driven development. Each tool works independently but is designed to complement the others.

## SpecScore — The Score

Defines *what* gets built. An open specification format with validation tooling.

- **Format:** Markdown and YAML — version-controlled, portable, no vendor lock-in
- **Tooling:** CLI linter, LSP for editor integration
- **Standalone value:** Write better specs, validate them automatically, link them to code

## Rehearse — The Rehearsal

Tests specifications automatically *before* implementation begins.

- Validates that specs are complete, consistent, and testable
- Catches ambiguity and gaps before agents start building
- Works with any project that uses the SpecScore format

[rehearse.ink](https://rehearse.ink)

## Synchestra — The Performance

Orchestrates multi-agent work *across* specifications.

- Coordinates AI agents claiming and executing tasks from SpecScore specs
- Schema-validated state at every commit — git is the protocol
- Scales from a single developer to distributed teams

[synchestra.io](https://synchestra.io)

## How They Fit Together

```
Write specs with SpecScore
         ↓
Validate with Rehearse
         ↓
Orchestrate with Synchestra
         ↓
Agents coordinate using SpecScore as the shared protocol
```

## Each Tool Standalone, Better Together

You don't need the full ecosystem. SpecScore works with any orchestration tool — Jira, Linear, your own scripts. Rehearse works in any project that uses SpecScore-formatted specs. Synchestra is optimized for SpecScore but is not required.

The recommended path: start with SpecScore. Add Rehearse when you want automated validation. Add Synchestra when you need multi-agent coordination.
```

- [x] **Step 2: Add ecosystem page to site-config.json**

Add a new page entry to the `pages` array in `tools/site-generator/site-config.json`, after the last "SpecScore for" entry:

```json
{
  "source": "docs/ecosystem.md",
  "slug": "ecosystem",
  "title": "Ecosystem",
  "nav": true,
  "navLabel": "Ecosystem",
  "navGroup": "Resources",
  "navOrder": 20
}
```

- [x] **Step 3: Build and verify**

```bash
cd tools/site-generator && pnpm build
```

Verify: `public/ecosystem.html` exists. Open it — content renders correctly. Sidebar shows "Ecosystem" under Resources. Navigation links work.

- [x] **Step 4: Commit**

```bash
git add docs/ecosystem.md tools/site-generator/site-config.json
git commit -m "content: add ecosystem page describing SpecScore, Rehearse, Synchestra"
```

---

## Task 4: Add Blog Infrastructure to Site Generator

**Files:**
- Create: `blog/` (directory at repo root)
- Create: `tools/site-generator/lib/build-blog.js`
- Create: `tools/site-generator/lib/build-blog.test.js`
- Create: `tools/site-generator/blog-post.html` (post template fragment)
- Create: `tools/site-generator/blog-index.html` (index template fragment)
- Modify: `tools/site-generator/build.js`
- Modify: `tools/site-generator/site-config.json`
- Modify: `tools/site-generator/assets/style.css`

The blog system reads Markdown files from `blog/` with a naming convention of `YYYY-MM-DD-slug.md`. Each file has YAML frontmatter with `title`, `description`, and `date`. The build generates individual post pages and an index page.

- [x] **Step 1: Write the failing test for blog post parsing**

Create `tools/site-generator/lib/build-blog.test.js`:

```js
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { parseBlogPost, buildBlogIndex } from './build-blog.js';

describe('parseBlogPost', () => {
  it('extracts frontmatter and body from a blog post', () => {
    const raw = `---
title: Why Specs Matter
description: A look at spec-driven development
date: 2026-04-14
---

# Why Specs Matter

Content here.`;

    const post = parseBlogPost(raw, '2026-04-14-why-specs-matter.md');
    assert.equal(post.title, 'Why Specs Matter');
    assert.equal(post.description, 'A look at spec-driven development');
    assert.equal(post.date, '2026-04-14');
    assert.equal(post.slug, 'blog/why-specs-matter');
    assert.ok(post.body.includes('# Why Specs Matter'));
    assert.ok(!post.body.includes('title:'));
  });

  it('throws on missing frontmatter', () => {
    assert.throws(
      () => parseBlogPost('# No frontmatter', 'bad.md'),
      /frontmatter/i
    );
  });
});

describe('buildBlogIndex', () => {
  it('sorts posts newest first', () => {
    const posts = [
      { title: 'Old', date: '2026-01-01', slug: 'blog/old', description: 'Old post' },
      { title: 'New', date: '2026-04-14', slug: 'blog/new', description: 'New post' },
    ];
    const sorted = buildBlogIndex(posts);
    assert.equal(sorted[0].title, 'New');
    assert.equal(sorted[1].title, 'Old');
  });
});
```

- [x] **Step 2: Run test to verify it fails**

```bash
cd tools/site-generator && node --test lib/build-blog.test.js
```

Expected: FAIL — `build-blog.js` does not exist.

- [x] **Step 3: Implement build-blog.js**

Create `tools/site-generator/lib/build-blog.js`:

```js
/**
 * Blog build utilities — parses frontmatter from blog posts
 * and builds a sorted index for the blog listing page.
 */

/**
 * Parses a blog post Markdown file with YAML frontmatter.
 * @param {string} raw - raw file content
 * @param {string} filename - e.g. '2026-04-14-why-specs-matter.md'
 * @returns {{ title: string, description: string, date: string, slug: string, body: string }}
 */
export function parseBlogPost(raw, filename) {
  const fmMatch = raw.match(/^---\n([\s\S]*?)\n---\n([\s\S]*)$/);
  if (!fmMatch) {
    throw new Error(`Missing frontmatter in ${filename}`);
  }

  const frontmatter = fmMatch[1];
  const body = fmMatch[2].trim();

  const title = frontmatter.match(/^title:\s*(.+)$/m)?.[1]?.trim();
  const description = frontmatter.match(/^description:\s*(.+)$/m)?.[1]?.trim();
  const date = frontmatter.match(/^date:\s*(.+)$/m)?.[1]?.trim();

  if (!title || !date) {
    throw new Error(`Missing required frontmatter fields (title, date) in ${filename}`);
  }

  // Derive slug from filename: 2026-04-14-why-specs-matter.md → blog/why-specs-matter
  const slugPart = filename.replace(/^\d{4}-\d{2}-\d{2}-/, '').replace(/\.md$/, '');
  const slug = `blog/${slugPart}`;

  return { title, description: description || '', date, slug, body };
}

/**
 * Returns posts sorted by date descending (newest first).
 * @param {Array<{title: string, date: string, slug: string, description: string}>} posts
 * @returns {Array} sorted copy
 */
export function buildBlogIndex(posts) {
  return [...posts].sort((a, b) => b.date.localeCompare(a.date));
}
```

- [x] **Step 4: Run test to verify it passes**

```bash
cd tools/site-generator && node --test lib/build-blog.test.js
```

Expected: PASS — both tests green.

- [x] **Step 5: Create blog post template fragment**

Create `tools/site-generator/blog-post.html`:

```html
<div class="blog-post-meta">
  <time>{{postDate}}</time>
</div>
{{content}}
```

- [x] **Step 6: Create blog index template fragment**

Create `tools/site-generator/blog-index.html`:

```html
<h1>Blog</h1>
<p class="landing-prose">Thoughts on spec-driven development, AI agent coordination, and building specifications that work.</p>
<ul class="blog-post-list">
{{blogEntries}}
</ul>
```

- [x] **Step 7: Add blog CSS**

Append to `tools/site-generator/assets/style.css`:

```css
/* === BLOG === */
.blog-post-meta {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 1.5rem;
}

.blog-post-meta time {
  font-variant-numeric: tabular-nums;
}

.blog-post-list {
  list-style: none;
  padding: 0;
}

.blog-post-list li {
  border-bottom: 1px solid var(--border);
  padding: 1rem 0;
}

.blog-post-list li:first-child {
  border-top: 1px solid var(--border);
}

.blog-post-list li a {
  font-family: 'Source Serif 4', Georgia, serif;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text);
  text-decoration: none;
  display: block;
  margin-bottom: 0.25rem;
}

.blog-post-list li a:hover {
  color: var(--accent);
}

.blog-post-list li .blog-entry-date {
  font-size: 0.8rem;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
}

.blog-post-list li .blog-entry-desc {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}
```

- [x] **Step 8: Integrate blog build into build.js**

Add blog building to `tools/site-generator/build.js`. After the existing `for (const page of config.pages)` loop (before the closing `console.log('\nDone...')`), add:

```js
  // --- Blog ---
  const { readdir } = await import('node:fs/promises');
  const { parseBlogPost, buildBlogIndex } = await import('./lib/build-blog.js');

  const blogDir = join(ROOT, 'blog');
  let blogFiles = [];
  try {
    blogFiles = (await readdir(blogDir)).filter(f => f.endsWith('.md')).sort();
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
```

- [x] **Step 9: Add blog to site navigation**

Add to the `pages` array in `tools/site-generator/site-config.json`:

```json
{
  "source": null,
  "slug": "blog",
  "title": "Blog",
  "nav": true,
  "navLabel": "Blog",
  "navGroup": "Resources",
  "navOrder": 19
}
```

**Note:** The blog index page is generated by the blog build step, not from a source Markdown file. The `source: null` entry is only for sidebar navigation. `load-config.js` will need a guard to skip null-source pages during the main build loop.

- [x] **Step 10: Update load-config.js to handle null-source pages**

In `tools/site-generator/lib/load-config.js`, the `sourceToSlug` map should skip pages with `null` source. Find where `sourceToSlug` is built and add a guard:

```js
// Only map pages that have a source file
if (page.source) {
  sourceToSlug.set(page.source, page.slug);
}
```

And in `build.js`, skip null-source pages in the main loop:

```js
for (const page of config.pages) {
  if (!page.source) continue;  // blog index — built separately
  // ... rest of loop
}
```

- [x] **Step 11: Create the blog directory with a README**

Create `blog/README.md`:

```markdown
# Blog

Blog posts for specscore.md. Each post is a Markdown file with YAML frontmatter.

## Naming Convention

Files are named `YYYY-MM-DD-slug.md` where the date prefix determines sort order and the slug becomes the URL path (`/blog/slug`).

## Frontmatter

```yaml
---
title: Post Title
description: One-line description for the index page
date: YYYY-MM-DD
---
```

## Outstanding Questions

None at this time.
```

- [x] **Step 12: Build and verify blog infrastructure works with no posts**

```bash
cd tools/site-generator && pnpm build
```

Expected: Build succeeds. `public/blog.html` does not exist (no posts yet). No errors. Sidebar shows "Blog" under Resources.

- [x] **Step 13: Commit**

```bash
git add tools/site-generator/lib/build-blog.js tools/site-generator/lib/build-blog.test.js tools/site-generator/blog-post.html tools/site-generator/blog-index.html tools/site-generator/build.js tools/site-generator/assets/style.css tools/site-generator/site-config.json tools/site-generator/lib/load-config.js blog/README.md
git commit -m "feat: add blog infrastructure to site generator"
```

---

## Task 5: Write Launch Blog Post — "What Is Spec-Driven Development?"

**Files:**
- Create: `blog/2026-04-14-what-is-spec-driven-development.md`

This is the flagship methodology article — establishes the term "spec-driven development" and positions SpecScore as its reference implementation.

- [x] **Step 1: Write the blog post**

Create `blog/2026-04-14-what-is-spec-driven-development.md`:

```markdown
---
title: What Is Spec-Driven Development?
description: Why AI agents need formal specifications — and how SpecScore makes it practical.
date: 2026-04-14
---

# What Is Spec-Driven Development?

Running one AI agent on one task works fine. Tell it what to build, it builds something. Maybe it's right. Maybe you spend an hour fixing what it misunderstood.

Running five agents across three platforms on a tree of interdependent tasks? That's where things collapse. Not because the agents are bad — because the specifications are.

## The Problem With Informal Specs

Most teams specify work in one of three ways:

1. **Jira tickets** — a title, a description, maybe some acceptance criteria written by someone in a hurry
2. **Notion docs** — rich, detailed, and never read by the agent that actually does the work
3. **Conversation** — "just make it work like the mockup"

None of these are machine-readable. None of them are validatable. None of them link back to the code that implements them.

When an AI agent works from an ambiguous spec, it guesses. Sometimes it guesses right. When it doesn't, you get rework — and rework with AI agents is expensive because it compounds across every downstream task that depended on the wrong assumption.

## Spec-Driven Development

Spec-driven development is a practice where specifications are:

- **Structured** — features, requirements, acceptance criteria, plans, and tasks follow a consistent format
- **Machine-readable** — AI agents can parse the spec and understand what to build without human translation
- **Validatable** — a linter catches ambiguity, missing fields, and structural problems before anyone starts building
- **Traceable** — source code links back to the spec it implements via inline annotations

The specification is the source of truth. Not the ticket. Not the conversation. Not the PR description.

## What This Looks Like in Practice

A SpecScore feature specification is Markdown with structure:

```yaml
title: User Authentication
status: draft
priority: high
```

Inside the feature directory, you define requirements with acceptance criteria:

```markdown
## Requirements

### R1: Email/password login

**Acceptance Criteria:**
- AC1: Given valid credentials, the system returns a session token
- AC2: Given invalid credentials, the system returns a 401 with a descriptive error
- AC3: Passwords are never logged or returned in API responses
```

Then you run `specscore lint` and it tells you what's missing, what's ambiguous, and what doesn't conform to the schema.

The developer (or agent) implementing this feature adds a source reference:

```go
// specscore:user-authentication/R1
func Login(email, password string) (*Session, error) {
```

Now you can trace from spec to code and from code to spec.

## Why Now?

Three things changed:

1. **AI agents are doing real work** — not just autocomplete, but multi-step implementation across codebases
2. **Specifications became the bottleneck** — agent capability outpaced spec quality
3. **Markdown won** — it's the lingua franca of developer tools, and now BAs and PMs write it too (thanks to AI assistants)

Spec-driven development is what test-driven development was for code quality — but for the layer above code. The spec is the contract. Everything flows from it.

## Getting Started

SpecScore is the open specification framework for spec-driven development:

- [Read the spec](/specifications) — understand the format
- [Install the CLI](/install) — start validating your specifications
- [Find your role](/for/developers) — guides for developers, PMs, QAs, BAs, architects, and project managers

The format is Markdown and YAML. The schema is published and versioned. The CLI is open source. No vendor lock-in.
```

- [x] **Step 2: Build and verify**

```bash
cd tools/site-generator && pnpm build
```

Verify: `public/blog/what-is-spec-driven-development.html` exists. `public/blog.html` exists and lists the post. Open both in a browser — content renders, links work, date shows.

- [x] **Step 3: Commit**

```bash
git add blog/2026-04-14-what-is-spec-driven-development.md
git commit -m "content: add flagship blog post on spec-driven development"
```

---

## Task 6: Write Launch Blog Post — "SpecScore vs. No Specs"

**Files:**
- Create: `blog/2026-04-14-specscore-vs-no-specs.md`

Comparison article targeting teams that don't use formal specifications today.

- [x] **Step 1: Write the blog post**

Create `blog/2026-04-14-specscore-vs-no-specs.md`:

```markdown
---
title: "SpecScore vs. No Specs: What You're Actually Losing"
description: What happens when AI agents work without formal specifications — and what changes when they have them.
date: 2026-04-14
---

# SpecScore vs. No Specs: What You're Actually Losing

You don't need SpecScore to build software with AI agents. Plenty of teams ship without formal specifications. The agents are smart enough to figure it out, right?

Sometimes. Here's what happens when they don't.

## The Cost of Ambiguity

Without formal specs, every agent interaction starts with interpretation. The agent reads your Jira ticket, your Slack message, or your PR description and decides what you meant. Every interpretation is a coin flip between "exactly right" and "plausible but wrong."

With one agent, this is manageable. You review the output, catch the misinterpretation, and iterate. The cost is your time.

With multiple agents working on interconnected tasks, misinterpretation compounds. Agent A misreads a requirement. Agent B builds on Agent A's output. Agent C depends on both. By the time you notice, three agents have built the wrong thing.

## What Formal Specs Change

| Without specs | With SpecScore |
|---|---|
| Agent guesses what "login should work" means | Agent reads structured acceptance criteria with pass/fail conditions |
| Requirements live in Jira, Notion, Slack, and someone's head | Requirements live in the repo, version-controlled, next to the code |
| "Is this done?" requires a human to check | `specscore lint` validates completeness and structure automatically |
| Code has no link to what it implements | `specscore:feature/R1` annotations trace code to requirements |
| New team member reads 40 tickets to understand a feature | New team member reads one feature spec with all context |

## The Adoption Path

SpecScore is designed for incremental adoption:

1. **Start with one feature.** Pick your next feature. Write a SpecScore spec for it. Use the CLI to lint it. See if it catches anything your Jira ticket missed.
2. **Add source references.** When you implement the feature, add `specscore:` annotations. See if traceability helps during code review.
3. **Expand to the team.** If it helped, write specs for the next sprint. Show your PM the role-based guide.

You don't need to convert your entire backlog. You don't need to change your tools. SpecScore specs are Markdown files in your repo — they work with whatever you already use.

## What SpecScore Is Not

- **Not a project management tool.** It doesn't replace Jira or Linear. It makes what's in them more precise.
- **Not a testing framework.** It defines what should be tested, not how to test it. (That's what [Rehearse](https://rehearse.ink) is for.)
- **Not a proprietary format.** It's open source, Markdown-based, and designed for portability.

## Try It

```bash
go install github.com/synchestra-io/specscore-cli/cmd/specscore@latest
specscore lint ./spec
```

[Read the spec](/specifications) to understand the format. [Find your role](/for/developers) for a guided introduction.
```

- [x] **Step 2: Build and verify**

```bash
cd tools/site-generator && pnpm build
```

Verify: `public/blog/specscore-vs-no-specs.html` exists. Blog index now lists two posts, newest first.

- [x] **Step 3: Commit**

```bash
git add blog/2026-04-14-specscore-vs-no-specs.md
git commit -m "content: add comparison blog post — SpecScore vs. no specs"
```

---

## Task 7: Write Launch Blog Post — "Linting Specs Like You Lint Code"

**Files:**
- Create: `blog/2026-04-14-linting-specs-like-you-lint-code.md`

Tooling-focused article showcasing the CLI.

- [x] **Step 1: Write the blog post**

Create `blog/2026-04-14-linting-specs-like-you-lint-code.md`:

```markdown
---
title: Linting Specs Like You Lint Code
description: How SpecScore's CLI catches ambiguity in your specifications the same way ESLint catches bugs in your JavaScript.
date: 2026-04-14
---

# Linting Specs Like You Lint Code

You lint your code. You lint your CSS. You lint your commit messages. Why aren't you linting your specifications?

Specifications are the upstream input to everything an AI agent builds. If the spec is ambiguous, the output is unpredictable. If the spec has missing fields, the agent fills in blanks with assumptions. If the spec structure is inconsistent, every agent interaction starts with "what format is this in?"

SpecScore brings the same discipline to specifications that ESLint brought to JavaScript.

## What the Linter Catches

The `specscore lint` command validates your specification files against the SpecScore schema:

```bash
$ specscore lint ./spec

spec/features/auth/README.md
  ✗ missing required field: status
  ✗ acceptance criteria AC3 has no testable condition

spec/features/search/README.md
  ✗ requirement R2 references undefined feature: indexing
  ⚠ no source references found for this feature

2 features, 3 errors, 1 warning
```

The linter checks for:

- **Structural completeness** — required fields, proper hierarchy, valid metadata
- **Referential integrity** — features, requirements, and acceptance criteria reference things that exist
- **Testability** — acceptance criteria contain verifiable conditions, not vague statements
- **Consistency** — naming conventions, status lifecycle, and format conformance

## Pluggable Rules

The linting engine is designed for extension. The core rules ship with SpecScore. Teams can add custom rules for their own conventions — naming patterns, required metadata fields, domain-specific validation.

## In Your Editor

The SpecScore LSP (coming soon) brings linting into your editor. Write a spec, see violations inline, fix them before you commit. The same feedback loop you have for code, applied to specifications.

## In CI/CD

Add `specscore lint` to your CI pipeline:

```yaml
# .github/workflows/specs.yml
- name: Lint specifications
  run: specscore lint ./spec
```

Specs that don't pass the linter don't get merged. The same gate you have for tests and code quality, applied to specifications.

## Try It

```bash
go install github.com/synchestra-io/specscore-cli/cmd/specscore@latest
specscore lint ./spec
```

[Read the specification format](/specifications) to understand what the linter validates.
```

- [x] **Step 2: Build and verify**

```bash
cd tools/site-generator && pnpm build
```

Verify: `public/blog/linting-specs-like-you-lint-code.html` exists. Blog index lists three posts.

- [x] **Step 3: Commit**

```bash
git add blog/2026-04-14-linting-specs-like-you-lint-code.md
git commit -m "content: add tooling blog post — linting specs like you lint code"
```

---

## Task 8: Prepare GitHub Repository for Public Launch

**Files:**
- Modify: `README.md` (repo root)
- Create: `CONTRIBUTING.md`

- [x] **Step 1: Rewrite root README.md for public-facing positioning**

Read the current `README.md` first, then rewrite to match the marketing positioning. The README should serve as both a GitHub landing page and a marketing surface.

Structure:
1. One-line description matching the tagline
2. What SpecScore is (elevator pitch)
3. Quick start (install + first lint)
4. Spec format overview (link to docs)
5. Role-based guides (links)
6. Ecosystem (SpecScore → Rehearse → Synchestra)
7. Contributing link
8. License

Keep it under 120 lines. Lead with the value proposition, not the technical architecture.

- [x] **Step 2: Create CONTRIBUTING.md**

Create `CONTRIBUTING.md`:

```markdown
# Contributing to SpecScore

SpecScore is an open-source specification framework. Contributions are welcome — whether that's improving the spec format, adding linter rules, fixing bugs, improving documentation, or writing guides.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<you>/specscore.git`
3. Install Go (1.22+): see [go.dev/dl](https://go.dev/dl)
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run ./...`

## What to Contribute

### Specification Format
The spec format lives in `spec/`. Changes to the format should include:
- Updated `README.md` for the affected feature
- Updated acceptance criteria
- Updated linter rules if applicable

### Linter Rules
Linter rules live in `pkg/lint/`. Each rule should:
- Have a clear, descriptive name
- Include tests in `pkg/lint/*_test.go`
- Be documented in the rule's error message

### Documentation
Guides live in `docs/`. Role-based guides are in `docs/for/`. Blog posts are in `blog/`.

## Code Standards

- Format with `gofmt`
- Pass `golangci-lint run ./...`
- Pass `go test ./...`
- Pass `go vet ./...`
- Handle all errors explicitly (see AGENTS.md)

## Pull Requests

- One feature per PR
- Include tests for new functionality
- Update documentation if behavior changes
- Reference the relevant spec feature if applicable

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
```

- [x] **Step 3: Verify links in README**

```bash
# Check that all linked paths exist
ls spec/README.md docs/installation.md docs/for/developers.md docs/for/product-owners.md docs/for/qas.md docs/for/business-analysts.md docs/for/project-managers.md docs/for/architects.md docs/ecosystem.md CONTRIBUTING.md
```

All files should exist.

- [x] **Step 4: Commit**

```bash
git add README.md CONTRIBUTING.md
git commit -m "docs: rewrite README and add CONTRIBUTING for public launch"
```

---

## Summary

| Task | What it delivers |
|---|---|
| 1 | Domain updated from specscore.org to specscore.md |
| 2 | Landing page with marketing positioning |
| 3 | Ecosystem page (SpecScore / Rehearse / Synchestra) |
| 4 | Blog infrastructure in site generator |
| 5 | Flagship blog post: "What Is Spec-Driven Development?" |
| 6 | Comparison blog post: "SpecScore vs. No Specs" |
| 7 | Tooling blog post: "Linting Specs Like You Lint Code" |
| 8 | Public-facing README and CONTRIBUTING.md |

After completing all 8 tasks, the Phase 0 (Foundation) from the launch strategy is complete — the site, content, and repo are ready for the Phase 1 soft launch.
