import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { rewriteLinks } from './rewrite-links.js';

const sourceToSlug = new Map([
  ['spec/README.md', 'index'],
  ['spec/features/README.md', 'specifications'],
  ['spec/features/feature/README.md', 'feature-specification'],
  ['spec/features/plan/README.md', 'plan-specification'],
]);

describe('rewriteLinks', () => {
  it('rewrites relative links to site slugs for HTML', () => {
    const md = 'See [Features](features/README.md) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [Features](/specifications) for details.');
  });

  it('rewrites relative links to .md for markdown output', () => {
    const md = 'See [Features](features/README.md) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'md');
    assert.equal(result, 'See [Features](/specifications.md) for details.');
  });

  it('rewrites parent-relative links', () => {
    const md = 'See [Plan](../plan/README.md) for planning.';
    const result = rewriteLinks(md, 'spec/features/feature/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [Plan](/plan-specification) for planning.');
  });

  it('leaves external URLs unchanged', () => {
    const md = 'Visit [GitHub](https://github.com/synchestra-io/specscore).';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'Visit [GitHub](https://github.com/synchestra-io/specscore).');
  });

  it('leaves anchor-only links unchanged', () => {
    const md = 'See [below](#behavior) for details.';
    const result = rewriteLinks(md, 'spec/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [below](#behavior) for details.');
  });

  it('leaves links to unknown sources unchanged', () => {
    const md = 'See [unknown](../proposals/README.md) for proposals.';
    const result = rewriteLinks(md, 'spec/features/feature/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [unknown](../proposals/README.md) for proposals.');
  });

  it('handles links with anchors, preserving the fragment', () => {
    const md = 'See [statuses](../feature/README.md#feature-statuses).';
    const result = rewriteLinks(md, 'spec/features/plan/README.md', sourceToSlug, 'html');
    assert.equal(result, 'See [statuses](/feature-specification#feature-statuses).');
  });
});
