# Sidebar Target Audience Section

## Goal

Add a "SpecScore for" section to the sidebar that helps visitors self-select into role-specific landing pages, while positioning SpecScore as useful across multiple professions.

## Design

### Sidebar group: "SpecScore for"

Placed between "Getting Started" and "Resources" in the sidebar.

**Collapsed (default):**

- Group label: **SpecScore for:**
- Below the label: comma-separated role names as plain text (e.g., "Developers, Product Owners, Business Analysts, QAs, Project Managers, Architects")
- A small expand/toggle button (chevron or similar)

**Expanded:**

- Same label
- Standard `<ul class="sidebar-nav">` with individual links, matching the style of other sidebar groups

### Roles (initial set)

| Role               | Slug                  | Nav label          |
| ------------------ | --------------------- | ------------------ |
| Developers         | `for/developers`      | Developers         |
| Product Owners     | `for/product-owners`  | Product Owners     |
| Business Analysts  | `for/business-analysts` | Business Analysts |
| QAs                | `for/qas`             | QAs                |
| Project Managers   | `for/project-managers` | Project Managers  |
| Architects         | `for/architects`      | Architects         |

### Pages

Each role gets a dedicated markdown file under `docs/for/<role-slug>.md` with lightweight draft content:

- A short intro paragraph on how SpecScore helps that role
- A few bullet points highlighting key benefits
- Content will be expanded in dedicated sessions later

Pages are rendered at `/for/<role-slug>.html`.

### Site config changes

Each role page added to `site-config.json` with:
- `navGroup`: `"SpecScore for"`
- `nav`: `true`
- Appropriate `navOrder` values

### Sidebar rendering changes

The `buildSidebarHtml` function in `render-page.js` needs to handle the "SpecScore for" group specially:

- Render a collapsed summary by default (comma-separated role names)
- Add an expand button that reveals the standard link list
- Expand/collapse state resets on page load (no localStorage persistence)

### CSS changes

- Styling for the collapsed summary text
- Styling for the expand/toggle button
- Transition for expand/collapse animation (optional, keep simple)

## Out of scope

- Deep role-specific content (separate sessions)
- localStorage persistence for expand/collapse state
- Landing page cards linking to role pages
