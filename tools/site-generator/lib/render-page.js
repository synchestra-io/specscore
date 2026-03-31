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
