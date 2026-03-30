# Static Site Generator for specscore.org

## Goal

Publish SpecScore specification content from `spec/` as a static website at specscore.org, serving both HTML and raw markdown at each URL.

## URL Scheme

Each specification page is available as HTML (default) and markdown (`.md` extension):

| Source file | HTML URL | Markdown URL |
|---|---|---|
| `spec/README.md` | `/` | `/index.md` |
| `spec/features/README.md` | `/specifications` | `/specifications.md` |
| `spec/features/feature/README.md` | `/feature-specification` | `/feature-specification.md` |
| `spec/features/development-plan/README.md` | `/development-plan-specification` | `/development-plan-specification.md` |
| `spec/features/acceptance-criteria/README.md` | `/acceptance-criteria-specification` | `/acceptance-criteria-specification.md` |
| `spec/features/project-definition/README.md` | `/project-definition-specification` | `/project-definition-specification.md` |
| `spec/features/source-references/README.md` | `/source-references-specification` | `/source-references-specification.md` |

The mapping is explicit, defined in `site-config.json`. New pages are added by editing the config — there is no automatic discovery.

### Clean URLs

HTML files are generated as `{slug}.html` (e.g., `feature-specification.html`). The hosting layer is expected to serve clean URLs — requests to `/feature-specification` serve `feature-specification.html`. Most static hosts (GitHub Pages, Cloudflare Pages, Netlify) support this out of the box. If the chosen host does not, the build can be adjusted to output `{slug}/index.html` instead.

## Domain Separation

- **specscore.org** hosts the SpecScore framework specifications (this site).
- **synchestra.io/{host}/{org}/{repo}/{path}** hosts per-project specification views for Synchestra users (separate concern, not part of this design).

Users reference SpecScore specs with links like `https://specscore.org/feature-specification` in their own feature definitions.

## Build Pipeline

```
spec/**/*.md
     │
     ▼
  build.js (Node.js)
     │
     ├── reads site-config.json (URL mapping, page metadata, nav order)
     ├── converts markdown → HTML via markdown-it
     ├── pre-renders mermaid code blocks to SVG via @mermaid-js/mermaid-cli
     ├── wraps HTML in template.html (header, nav, content slot, footer)
     ├── rewrites internal links (spec-relative → site URLs)
     └── copies processed .md files alongside HTML
     │
     ▼
  public/ (.gitignored)
```

### Markdown Processing

- **Parser:** `markdown-it` with GFM tables, fenced code blocks.
- **Mermaid diagrams:** Extracted from fenced code blocks at build time, rendered to inline SVG via `@mermaid-js/mermaid-cli` (uses Puppeteer). No client-side mermaid.js needed.
- **Link rewriting:** Relative links in the source markdown (e.g., `[Development Plan](../development-plan/README.md)`) are rewritten to site URLs (e.g., `/development-plan-specification`) in both the HTML output and the published `.md` files. The original source files in `spec/` are never modified.

### Link Rewriting Details

The build script maintains a mapping from source paths to site URLs (derived from `site-config.json`). During processing:

1. Parse all markdown links and image references.
2. Resolve relative paths against the source file's location.
3. If the resolved path matches a known source file in the config, replace with the corresponding site URL.
4. Links to unknown paths (external URLs, anchors) are left unchanged.

This applies to both the HTML output and the published `.md` files, ensuring all internal navigation works regardless of format.

In HTML output, rewritten links point to the clean URL (e.g., `/feature-specification`). In published `.md` files, rewritten links point to the `.md` sibling (e.g., `/feature-specification.md`), so markdown-to-markdown navigation is self-contained.

## Project Layout

```
tools/site-generator/
├── package.json          # dependencies: markdown-it, @mermaid-js/mermaid-cli
├── build.js              # build script
├── site-config.json      # URL mapping, page titles, nav order
└── template.html         # HTML shell (header, nav, content slot, footer)
```

### site-config.json Structure

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
      "navOrder": 1
    },
    {
      "source": "spec/features/feature/README.md",
      "slug": "feature-specification",
      "title": "Feature Specification",
      "nav": true,
      "navOrder": 2
    }
  ]
}
```

Each entry defines: source file path, output slug, page title (for `<title>` and heading), whether it appears in navigation, and its nav ordering.

## Build Output

```
public/
├── index.html
├── index.md
├── specifications.html
├── specifications.md
├── feature-specification.html
├── feature-specification.md
├── development-plan-specification.html
├── development-plan-specification.md
├── acceptance-criteria-specification.html
├── acceptance-criteria-specification.md
├── project-definition-specification.html
├── project-definition-specification.md
├── source-references-specification.html
├── source-references-specification.md
└── assets/
    └── style.css
```

The `public/` directory is added to `.gitignore`. Mermaid SVGs are inlined in the HTML, not served as separate files.

## HTML Template

Single `template.html` with a simple structure:

```
┌─────────────────────────────────────┐
│  SpecScore wordmark    [theme toggle]│
│  Nav: Home | Specifications | GitHub│
├─────────────────────────────────────┤
│                                     │
│  Page content (rendered markdown)   │
│                                     │
│  "View as Markdown" link            │
│                                     │
├─────────────────────────────────────┤
│  Footer: Apache 2.0, GitHub link    │
└─────────────────────────────────────┘
```

- Header with nav links (generated from `site-config.json` entries where `nav: true`).
- Content area with rendered markdown.
- "View as Markdown" link pointing to the `.md` sibling.
- Footer with license and GitHub repo link.

## Styling

- **Single CSS file** (`assets/style.css`), no framework.
- **CSS custom properties** for all colors, enabling theme switching.
- Clean typography optimized for technical documentation: readable body text, clear heading hierarchy, styled code blocks and tables.
- **Responsive** layout — works on mobile, optimized for desktop reading.

## Theme Toggle

- **Default:** Respects `prefers-color-scheme` media query from the browser. Falls back to light if no preference detected.
- **Toggle button** in the header. Persists choice to `localStorage`.
- **Inline `<script>` in `<head>`** (~15 lines) reads `localStorage` and applies `data-theme` attribute before first paint, preventing flash of wrong theme.
- CSS uses `[data-theme="dark"]` selectors to override custom properties.

## JavaScript

The site ships minimal JavaScript:

1. **Theme toggle** — inlined in `<head>`, ~15 lines. Reads `localStorage`, applies theme, handles toggle clicks.
2. **Google Analytics** — standard GA snippet.

No other JavaScript. Mermaid diagrams are pre-rendered SVGs. No search, no dynamic features.

## CI / Deployment

- Build runs via `npm run build` from `tools/site-generator/`.
- Puppeteer browser binary cached in GitHub Actions via `actions/setup-node` npm cache support.
- Output is plain static files — deployable to GitHub Pages, Cloudflare Pages, Netlify, or any static host.
- Deployment target is not decided in this spec — the build produces host-agnostic output.

## What This Design Does Not Cover

- **Hosting provider selection** — the build output is static files, deployable anywhere.
- **Custom domain DNS setup** — depends on hosting choice.
- **Search functionality** — can be added later with client-side search (e.g., pagefind) if needed.
- **Source-references URL fix** — the `source-references` spec currently references `specscore.org` URLs where it should say `synchestra.io`. This is a separate fix.
- **Visual design details** — color palette, fonts, spacing. These will be decided during implementation.
