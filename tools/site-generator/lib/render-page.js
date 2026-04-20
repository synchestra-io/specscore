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
          if (item.external) {
            return `<li><a href="${item.href}" target="_blank" rel="noopener noreferrer">${item.navLabel} &#8599;</a></li>`;
          }
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
  const alternateLinkHtml = showViewMarkdown
    ? `<link rel="alternate" type="text/markdown" href="${mdUrl}">`
    : '';

  return template
    .replace(/\{\{title\}\}/g, title)
    .replace(/\{\{sidebar\}\}/g, sidebarHtml)
    .replace(/\{\{eyebrow\}\}/g, eyebrowHtml)
    .replace(/\{\{content\}\}/g, content)
    .replace(/\{\{viewMarkdown\}\}/g, viewMarkdownHtml)
    .replace(/\{\{alternateLink\}\}/g, alternateLinkHtml);
}
