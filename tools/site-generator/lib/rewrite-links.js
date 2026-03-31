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
