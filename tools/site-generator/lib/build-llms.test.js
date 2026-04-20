import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { buildLlmsIndex } from './build-llms.js';

describe('buildLlmsIndex', () => {
  const config = {
    pages: [
      { slug: 'index', title: 'Home' }, // no llmsDescription → excluded
      {
        slug: 'feature-specification',
        title: 'Feature Specification',
        llmsGroup: 'Specifications',
        llmsDescription: 'The atomic unit of product specification.',
      },
      {
        slug: 'plan-specification',
        title: 'Plan Specification',
        llmsGroup: 'Specifications',
        llmsDescription: 'Planning documents that bridge Feature specs to execution.',
      },
      {
        slug: 'install',
        title: 'Installation',
        llmsGroup: 'Getting Started',
        llmsDescription: 'How to install the specscore CLI.',
      },
    ],
  };

  it('starts with a # SpecScore heading and a blockquote summary', () => {
    const out = buildLlmsIndex(config);
    assert.match(out, /^# SpecScore\n/);
    assert.match(out, /\n> SpecScore is a Markdown-based specification language/);
  });

  it('includes only pages with an llmsDescription', () => {
    const out = buildLlmsIndex(config);
    assert.ok(out.includes('Feature Specification'));
    assert.ok(out.includes('Plan Specification'));
    assert.ok(out.includes('Installation'));
    assert.ok(!out.includes('Home'));
  });

  it('groups entries by llmsGroup in source order', () => {
    const out = buildLlmsIndex(config);
    const specsIdx = out.indexOf('## Specifications');
    const gettingIdx = out.indexOf('## Getting Started');
    assert.ok(specsIdx > 0);
    assert.ok(gettingIdx > specsIdx, 'Specifications group should appear before Getting Started');
  });

  it('emits absolute .md URLs to specscore.md', () => {
    const out = buildLlmsIndex(config);
    assert.ok(out.includes('https://specscore.md/feature-specification.md'));
    assert.ok(out.includes('https://specscore.md/install.md'));
  });

  it('uses an Unknown group bucket when llmsGroup is missing', () => {
    const localConfig = {
      pages: [
        { slug: 'orphan', title: 'Orphan', llmsDescription: 'No group.' },
      ],
    };
    const out = buildLlmsIndex(localConfig);
    assert.ok(out.includes('## Pages'));
    assert.ok(out.includes('Orphan'));
  });
});
