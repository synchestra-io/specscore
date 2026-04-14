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
