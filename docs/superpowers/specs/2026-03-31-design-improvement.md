# specscore.org Design Improvement

## Goal

Redesign the specscore.org static site with a "Clarity & Standards" aesthetic — editorial, authoritative, readable — replacing the current GitHub-docs-clone appearance.

## Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Personality | Clarity & Standards | Feels like a specification or RFC — W3C, Stripe Docs, MDN |
| Accent color | Indigo (`#4f46e5`) | Intellectual, precise, developer-native |
| Heading typeface | Source Serif 4 (Google Fonts) | Designed for long-form documentation; readable at all heading levels |
| Body typeface | System UI stack (unchanged) | No need to load an extra font for body text |
| Navigation | Left sidebar | Scales as spec pages grow; shows hierarchy |
| Homepage | Dedicated landing page | `spec/README.md` content moves to `/specifications`; `/` is hand-crafted HTML |

## Layout

A persistent sidebar sits left of the main content area. The top header contains only the wordmark and right-aligned GitHub link + theme toggle. The sidebar replaces the current top nav.

```
┌──────────────────────────────────────────────────────┐
│  SpecScore                          GitHub  [theme]  │  ← sticky header, 56px
├──────────┬───────────────────────────────────────────┤
│          │                                           │
│ Sidebar  │  Main content                             │
│ 240px    │  max-width 780px, padding 3rem 4rem       │
│          │                                           │
│ ─ Getting│  [eyebrow label]                          │
│   Started│  H1 (Source Serif 4)                      │
│   Home   │  ─── indigo rule ───                      │
│   Specs  │                                           │
│          │  Body text, code blocks, tables, diagrams │
│ ─ Feature│                                           │
│   Specs  │                                           │
│   Feature│                                           │
│   Acc Cri│                                           │
│   Dev Plan                                           │
│   Src Ref│                                           │
│   Proj Def                                           │
│          │                                           │
│ ─Resources                                           │
│   GitHub │                                           │
└──────────┴───────────────────────────────────────────┘
```

### Sidebar structure

Sidebar groups are defined by `navGroup` fields in `site-config.json`. Entries with `nav: false` do not appear. The active page link is highlighted with indigo background + indigo text.

`site-config.json` gains a `navGroup` field per page entry:

```json
{
  "source": "spec/features/feature/README.md",
  "slug": "feature-specification",
  "title": "Feature Specification",
  "nav": true,
  "navLabel": "Feature",
  "navOrder": 2,
  "navGroup": "Feature Specs"
}
```

Groups appear in the order they are first encountered by `navOrder`.

### Active page detection

`build.js` injects `data-slug="<slug>"` on `<body>`. A small inline script sets `.active` on the matching sidebar link by comparing `data-slug` to each link's `href`.

## Pages

### Spec pages (markdown-generated)

Each page has:
- An eyebrow label (the page's `navGroup` value from `site-config.json`, rendered above h1)
- H1 in Source Serif 4, followed by a 2.5rem × 3px indigo rule
- Body content rendered from markdown

### Landing page (`/`)

A hand-crafted `landing.html` in `tools/site-generator/`. `build.js` copies it to `public/index.html` instead of generating from `spec/README.md`.

The landing page has two sections:

**Hero:**
- H1: "The specification framework for AI-driven development"
- Tagline paragraph describing SpecScore in 1–2 sentences
- Two CTAs: "View Specifications" (primary, indigo button → `/specifications`) and "GitHub ↗" (secondary, outline button)

**Spec list:**
- Section label: "Specifications" (indigo uppercase)
- A simple `<ul>` list, one item per spec page (all entries where `navGroup === "Feature Specs"`)
- Each item: spec title as a link + one-line description
- Descriptions are hardcoded in `landing.html` (not pulled from `site-config.json`)

The landing page uses the same `template.html` shell (header, sidebar, footer) via a different injection path in `build.js`.

## Typography Scale

| Element | Font | Size | Weight |
|---|---|---|---|
| Wordmark | Source Serif 4 | 1.25rem | 700 |
| H1 | Source Serif 4 | 2.25rem | 700 |
| H2 | Source Serif 4 | 1.375rem | 600 |
| H3 | Source Serif 4 | 1.1rem | 600 |
| Body | System UI | 16px | 400 |
| Code (inline) | SFMono / Menlo | 0.85em | 400 |
| Code (block) | SFMono / Menlo | 0.875rem | 400 |
| Table headers | System UI | 0.8rem | 600 |
| Eyebrow / labels | System UI | 0.7rem | 700 |

Source Serif 4 is loaded from Google Fonts with `opsz` (optical size) axis and weights 400, 600, 700.

## Color System

### Light mode

```css
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
--accent-muted: #818cf8;
--code-bg: #f3f4f6;
--code-border: #e5e7eb;
--code-accent: #4f46e5;
--link: #4f46e5;
```

### Dark mode

Warmer dark palette — not GitHub's blue-grey. Feels like a well-lit reading environment at night.

```css
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
--accent-muted: #6366f1;
--code-bg: #27272a;
--code-border: #3f3f46;
--code-accent: #818cf8;
--link: #818cf8;
```

Dark mode is activated by `[data-theme="dark"]` on `<html>`, exactly as today. The theme toggle and localStorage persistence logic is unchanged.

## Component Details

### Code blocks

```
background: --code-bg
border: 1px solid --code-border
border-left: 3px solid --code-accent   ← indigo accent stripe
border-radius: 0 6px 6px 0
padding: 1rem 1.25rem
```

No syntax highlighting — color is not added to token classes. The indigo left border provides visual identity without requiring a highlight.js dependency.

### Tables

- Header row: `--bg-secondary` background, `text-transform: uppercase`, `font-size: 0.8rem`, 2px solid bottom border
- Data rows: subtle hover state (`--bg-secondary`)
- No alternating row stripes (cleaner at small sizes)

### Mermaid diagrams

Wrapped in a container:

```html
<div class="mermaid-wrap">
  <div class="mermaid-label">Diagram — [auto-generated caption]</div>
  <div class="mermaid-body"><!-- inline SVG --></div>
</div>
```

`build.js` already extracts mermaid fenced blocks and replaces them with SVG. The wrapping `<div class="mermaid-wrap">` is added at the same substitution point. The caption is either taken from the fenced block info string (e.g., ` ```mermaid title="..."`) or defaults to "Diagram".

`mermaid-wrap`: `border: 1px solid --border; border-radius: 8px; overflow: hidden;`
`mermaid-label`: small uppercase label, `--bg-secondary` background, bottom border
`mermaid-body`: `padding: 1.5rem; text-align: center;`

### "View as Markdown" link

Remains at the bottom of each spec page, separated by a top border. Unchanged in placement, slightly smaller (`0.8rem`, `--text-muted`).

## Changes to `site-config.json`

Add `navGroup` to every page entry. Proposed groupings:

| Page | navGroup |
|---|---|
| Home (landing) | *(not in sidebar — landing is a separate page)* |
| Specifications | Getting Started |
| Feature | Feature Specs |
| Acceptance Criteria | Feature Specs |
| Development Plan | Feature Specs |
| Source References | Feature Specs |
| Project Definition | Feature Specs |

The Home entry remains in `site-config.json` for slug/title metadata but `nav: false` so it does not appear in the sidebar.

## Changes to `build.js`

1. **Sidebar generation** — replace the existing `{{nav}}` flat link list with grouped sidebar HTML, derived from `site-config.json` entries sorted by `navOrder` and grouped by `navGroup`.
2. **Active page injection** — inject `data-slug="<slug>"` on `<body>` per page; add inline script to mark active link.
3. **Landing page handling** — detect the `index` slug; instead of rendering `spec/README.md`, copy `landing.html` (processed through the same template) to `public/index.html`.
4. **Mermaid wrapper** — wrap substituted SVG blocks in `.mermaid-wrap` / `.mermaid-label` / `.mermaid-body` divs.
5. **Eyebrow label** — inject `<div class="page-eyebrow">{{navGroup}}</div>` before the page's h1 (or before the first `<h1>` in the rendered HTML).

## Changes to `template.html`

- Replace `<nav>{{nav}}</nav>` in the header with just the wordmark + right-side controls (GitHub link, theme toggle).
- Add `<aside class="sidebar">{{sidebar}}</aside>` as a sibling of `<main>`, both inside a `.site-body` wrapper div.
- Add `data-slug="{{slug}}"` to `<body>`.
- Add inline active-link script after the sidebar HTML.

## What This Design Does Not Cover

- **Syntax highlighting** — can be added later with highlight.js or Shiki; deferred to keep the build simple.
- **Search** — can be added later with pagefind.
- **Landing page spec card descriptions** — the descriptions in `landing.html` are hardcoded; no mechanism to pull them from markdown or config.
- **Mobile sidebar** — responsive collapse (hamburger menu) is out of scope; on narrow screens the sidebar stacks above content as it does today, or can be hidden. A follow-on task.
