# Sidebar Target Audience Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an expandable "SpecScore for" sidebar section with role-specific landing pages.

**Architecture:** Six lightweight markdown pages under `docs/for/`, registered in `site-config.json`. The sidebar renderer gets a special collapsible group that shows a comma-separated role summary by default and expands into individual links. A small JS toggle in `template.html` handles expand/collapse.

**Tech Stack:** Node.js site generator (markdown-it), vanilla HTML/CSS/JS

---

### Task 1: Create role markdown pages

**Files:**
- Create: `docs/for/developers.md`
- Create: `docs/for/product-owners.md`
- Create: `docs/for/business-analysts.md`
- Create: `docs/for/qas.md`
- Create: `docs/for/project-managers.md`
- Create: `docs/for/architects.md`

- [ ] **Step 1: Create `docs/for/` directory and all six markdown files**

```markdown
<!-- docs/for/developers.md -->
# SpecScore for Developers

SpecScore gives developers clear, structured specifications to build against — reducing ambiguity and rework.

- **Unambiguous acceptance criteria** — know exactly what "done" means before writing code
- **Structured feature specs** — consistent format across every feature in the project
- **Source references** — trace every requirement back to its origin
- **Development plans** — break features into concrete, estimable tasks
```

```markdown
<!-- docs/for/product-owners.md -->
# SpecScore for Product Owners

SpecScore helps product owners translate business needs into specifications that development teams can execute on confidently.

- **Project definitions** — capture vision, goals, and constraints in a standard format
- **Feature specifications** — describe what you need built, clearly and completely
- **Acceptance criteria** — define success conditions the whole team agrees on
- **Traceability** — link every feature back to the business need it serves
```

```markdown
<!-- docs/for/business-analysts.md -->
# SpecScore for Business Analysts

SpecScore provides a structured framework for capturing and organizing requirements that bridges business stakeholders and technical teams.

- **Source references** — document where each requirement originated
- **Feature specifications** — formalize analysis into actionable specs
- **Acceptance criteria** — translate business rules into verifiable conditions
- **Project definitions** — maintain a single source of truth for project scope
```

```markdown
<!-- docs/for/qas.md -->
# SpecScore for QAs

SpecScore gives QA teams the specification clarity they need to build comprehensive test strategies and catch issues early.

- **Acceptance criteria** — structured, testable conditions for every feature
- **Feature specifications** — complete context for what each feature should do
- **Development plans** — understand implementation approach to design better tests
- **Traceability** — verify that every requirement has corresponding test coverage
```

```markdown
<!-- docs/for/project-managers.md -->
# SpecScore for Project Managers

SpecScore helps project managers maintain visibility into what's being built, why, and how — with a consistent specification format across the project.

- **Project definitions** — standardized scope and goals documentation
- **Development plans** — structured task breakdowns for estimation and tracking
- **Feature specifications** — clear deliverables the whole team aligns on
- **Source references** — audit trail from stakeholder request to implementation
```

```markdown
<!-- docs/for/architects.md -->
# SpecScore for Architects

SpecScore gives architects a structured specification format that captures technical requirements alongside business context.

- **Feature specifications** — detailed enough to inform architectural decisions
- **Development plans** — validate that implementation approach aligns with system design
- **Source references** — understand the full context behind each requirement
- **Project definitions** — capture cross-cutting concerns and technical constraints
```

- [ ] **Step 2: Verify files exist**

Run: `ls -la docs/for/`
Expected: Six `.md` files listed.

- [ ] **Step 3: Commit**

```bash
git add docs/for/
git commit -m "feat(site): add draft role landing pages under docs/for/"
```

---

### Task 2: Register role pages in site config

**Files:**
- Modify: `tools/site-generator/site-config.json`

- [ ] **Step 1: Add six page entries to site-config.json**

Add these entries to the `pages` array, after the existing "Getting Started" entries (after `navOrder: 8`):

```json
{
  "source": "docs/for/developers.md",
  "slug": "for/developers",
  "title": "SpecScore for Developers",
  "nav": true,
  "navLabel": "Developers",
  "navGroup": "SpecScore for",
  "navOrder": 10
},
{
  "source": "docs/for/product-owners.md",
  "slug": "for/product-owners",
  "title": "SpecScore for Product Owners",
  "nav": true,
  "navLabel": "Product Owners",
  "navGroup": "SpecScore for",
  "navOrder": 11
},
{
  "source": "docs/for/business-analysts.md",
  "slug": "for/business-analysts",
  "title": "SpecScore for Business Analysts",
  "nav": true,
  "navLabel": "Business Analysts",
  "navGroup": "SpecScore for",
  "navOrder": 12
},
{
  "source": "docs/for/qas.md",
  "slug": "for/qas",
  "title": "SpecScore for QAs",
  "nav": true,
  "navLabel": "QAs",
  "navGroup": "SpecScore for",
  "navOrder": 13
},
{
  "source": "docs/for/project-managers.md",
  "slug": "for/project-managers",
  "title": "SpecScore for Project Managers",
  "nav": true,
  "navLabel": "Project Managers",
  "navGroup": "SpecScore for",
  "navOrder": 14
},
{
  "source": "docs/for/architects.md",
  "slug": "for/architects",
  "title": "SpecScore for Architects",
  "nav": true,
  "navLabel": "Architects",
  "navGroup": "SpecScore for",
  "navOrder": 15
}
```

- [ ] **Step 2: Run existing config tests to verify nothing broke**

Run: `node --test tools/site-generator/lib/load-config.test.js`
Expected: All existing tests pass.

- [ ] **Step 3: Commit**

```bash
git add tools/site-generator/site-config.json
git commit -m "feat(site): register role pages in site config"
```

---

### Task 3: Add collapsible sidebar group support

**Files:**
- Modify: `tools/site-generator/lib/load-config.js`
- Modify: `tools/site-generator/lib/load-config.test.js`

- [ ] **Step 1: Write failing test for `collapsible` flag on sidebar groups**

Add this test to `tools/site-generator/lib/load-config.test.js`:

```js
it('marks "SpecScore for" sidebar group as collapsible', async () => {
  const config = await loadConfig(new URL('../site-config.json', import.meta.url));
  const forGroup = config.sidebarGroups.find((g) => g.label === 'SpecScore for');
  assert.ok(forGroup, '"SpecScore for" group should exist');
  assert.equal(forGroup.collapsible, true);
});
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node --test tools/site-generator/lib/load-config.test.js`
Expected: FAIL — `forGroup.collapsible` is `undefined`.

- [ ] **Step 3: Add collapsible flag in loadConfig**

In `tools/site-generator/lib/load-config.js`, replace the `sidebarGroups` mapping at the end:

```js
const COLLAPSIBLE_GROUPS = new Set(['SpecScore for']);

const sidebarGroups = Array.from(groupMap.entries()).map(([label, items]) => ({
  label,
  items,
  collapsible: COLLAPSIBLE_GROUPS.has(label),
}));
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `node --test tools/site-generator/lib/load-config.test.js`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add tools/site-generator/lib/load-config.js tools/site-generator/lib/load-config.test.js
git commit -m "feat(site): add collapsible flag to sidebar groups"
```

---

### Task 4: Render collapsible sidebar HTML

**Files:**
- Modify: `tools/site-generator/lib/render-page.js`

- [ ] **Step 1: Update `buildSidebarHtml` to handle collapsible groups**

Replace the `buildSidebarHtml` function in `tools/site-generator/lib/render-page.js`:

```js
function buildSidebarHtml(sidebarGroups, currentSlug) {
  return sidebarGroups
    .map((group) => {
      const links = group.items
        .map((item) => {
          if (item.external) {
            return `<li><a href="${item.href}" target="_blank" rel="noopener noreferrer">${item.navLabel} &#8599;</a></li>`;
          }
          const href = item.slug === 'index' ? '/' : `/${item.slug}.html`;
          const active = item.slug === currentSlug ? ' class="active"' : '';
          return `<li><a href="${href}"${active}>${item.navLabel}</a></li>`;
        })
        .join('\n          ');

      if (group.collapsible) {
        const summary = group.items.map((item) => item.navLabel).join(', ');
        return `<div class="sidebar-section sidebar-collapsible" data-collapsible>
        <span class="sidebar-label">${group.label}</span>
        <p class="sidebar-summary">${summary}</p>
        <button class="sidebar-expand" type="button" aria-expanded="false" aria-label="Expand ${group.label}">Show all &#9662;</button>
        <ul class="sidebar-nav sidebar-collapsible-list">
          ${links}
        </ul>
      </div>`;
      }

      return `<div class="sidebar-section">
        <span class="sidebar-label">${group.label}</span>
        <ul class="sidebar-nav">
          ${links}
        </ul>
      </div>`;
    })
    .join('\n      ');
}
```

- [ ] **Step 2: Build the site and verify the HTML contains the collapsible markup**

Run: `node tools/site-generator/build.js && grep -c "sidebar-collapsible" public/specifications.html`
Expected: Build succeeds, grep returns `1` or more.

- [ ] **Step 3: Commit**

```bash
git add tools/site-generator/lib/render-page.js
git commit -m "feat(site): render collapsible sidebar group HTML"
```

---

### Task 5: Add expand/collapse CSS

**Files:**
- Modify: `tools/site-generator/assets/style.css`

- [ ] **Step 1: Add CSS for collapsible sidebar group**

Add after the `.sidebar-nav li a.active` rule block (after line 183):

```css
/* === COLLAPSIBLE SIDEBAR GROUP === */
.sidebar-summary {
  font-size: 0.8rem;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0;
  padding: 0 0.6rem;
}

.sidebar-expand {
  display: block;
  background: none;
  border: none;
  font-size: 0.75rem;
  color: var(--accent);
  cursor: pointer;
  padding: 0.3rem 0.6rem;
  margin-top: 0.25rem;
}

.sidebar-expand:hover {
  text-decoration: underline;
}

.sidebar-collapsible-list {
  display: none;
}

.sidebar-collapsible.is-expanded .sidebar-summary {
  display: none;
}

.sidebar-collapsible.is-expanded .sidebar-expand {
  display: none;
}

.sidebar-collapsible.is-expanded .sidebar-collapsible-list {
  display: block;
}
```

- [ ] **Step 2: Commit**

```bash
git add tools/site-generator/assets/style.css
git commit -m "feat(site): add CSS for collapsible sidebar group"
```

---

### Task 6: Add expand/collapse JavaScript

**Files:**
- Modify: `tools/site-generator/template.html`

- [ ] **Step 1: Add toggle script to template.html**

Add this script before the closing `</body>` tag, after the existing theme-toggle script:

```html
<script>
  document.querySelectorAll('[data-collapsible] .sidebar-expand').forEach(function(btn) {
    btn.addEventListener('click', function() {
      var section = btn.closest('[data-collapsible]');
      section.classList.toggle('is-expanded');
      var expanded = section.classList.contains('is-expanded');
      btn.setAttribute('aria-expanded', expanded);
    });
  });
</script>
```

- [ ] **Step 2: Build and manually verify**

Run: `node tools/site-generator/build.js`
Expected: Build succeeds. Open `public/specifications.html` in a browser — the "SpecScore for" section should show a comma-separated list with a "Show all" button that expands into individual links.

- [ ] **Step 3: Commit**

```bash
git add tools/site-generator/template.html
git commit -m "feat(site): add expand/collapse JS for collapsible sidebar"
```

---

### Task 7: Auto-expand when viewing a role page

**Files:**
- Modify: `tools/site-generator/template.html`

- [ ] **Step 1: Update the collapsible script to auto-expand when the current page is in the group**

Replace the collapsible script added in Task 6 with:

```html
<script>
  document.querySelectorAll('[data-collapsible]').forEach(function(section) {
    var btn = section.querySelector('.sidebar-expand');
    var hasActive = section.querySelector('.sidebar-nav li a.active');
    if (hasActive) {
      section.classList.add('is-expanded');
      btn.setAttribute('aria-expanded', 'true');
    }
    btn.addEventListener('click', function() {
      section.classList.toggle('is-expanded');
      var expanded = section.classList.contains('is-expanded');
      btn.setAttribute('aria-expanded', String(expanded));
    });
  });
</script>
```

- [ ] **Step 2: Build and verify**

Run: `node tools/site-generator/build.js && grep "is-expanded" public/for/developers.html | head -1`
Expected: Build succeeds. The `for/developers.html` page should NOT have `is-expanded` hard-coded in HTML (it's added by JS at runtime). Instead, verify by opening the file in a browser — the collapsible section should auto-expand because the active link is inside it.

- [ ] **Step 3: Commit**

```bash
git add tools/site-generator/template.html
git commit -m "feat(site): auto-expand collapsible sidebar on active role page"
```

---

### Task 8: Run full build and verify

- [ ] **Step 1: Run all tests**

Run: `node --test tools/site-generator/lib/load-config.test.js`
Expected: All tests pass.

- [ ] **Step 2: Run full build**

Run: `node tools/site-generator/build.js`
Expected: All pages build successfully, including the six new `for/*.html` pages.

- [ ] **Step 3: Verify all role pages were generated**

Run: `ls public/for/`
Expected:
```
architects.html
architects.md
business-analysts.html
business-analysts.md
developers.html
developers.md
product-owners.html
product-owners.md
project-managers.html
project-managers.md
qas.html
qas.md
```

- [ ] **Step 4: Spot-check sidebar HTML in a role page**

Run: `grep -A2 "sidebar-collapsible" public/for/developers.html | head -5`
Expected: Shows the `sidebar-collapsible` div with the `data-collapsible` attribute.
