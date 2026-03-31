import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { loadConfig } from './load-config.js';

describe('loadConfig', () => {
  it('loads pages from site-config.json', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(Array.isArray(config.pages));
    assert.ok(config.pages.length > 0);
    assert.equal(config.pages[0].slug, 'index');
    assert.equal(config.pages[0].source, 'spec/README.md');
  });

  it('builds a sourceToSlug map', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.equal(config.sourceToSlug.get('spec/README.md'), 'index');
    assert.equal(
      config.sourceToSlug.get('spec/features/feature/README.md'),
      'feature-specification'
    );
  });

  it('builds sorted nav items', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(config.navItems.length > 0);
    // nav items should be sorted by navOrder
    for (let i = 1; i < config.navItems.length; i++) {
      assert.ok(config.navItems[i].navOrder >= config.navItems[i - 1].navOrder);
    }
  });
});
