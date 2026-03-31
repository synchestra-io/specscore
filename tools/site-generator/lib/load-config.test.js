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

  it('builds sidebar groups from navGroup fields', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    assert.ok(Array.isArray(config.sidebarGroups));
    assert.ok(config.sidebarGroups.length > 0);
    for (const group of config.sidebarGroups) {
      assert.ok(typeof group.label === 'string');
      assert.ok(Array.isArray(group.items));
      assert.ok(group.items.length > 0);
    }
  });

  it('sidebar groups contain Feature Specs entries', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    const featureGroup = config.sidebarGroups.find((g) => g.label === 'Feature Specs');
    assert.ok(featureGroup, 'Feature Specs group should exist');
    assert.ok(featureGroup.items.some((p) => p.slug === 'feature-specification'));
  });

  it('excludes nav:false pages from sidebar groups', async () => {
    const config = await loadConfig(new URL('../site-config.json', import.meta.url));
    const allItems = config.sidebarGroups.flatMap((g) => g.items);
    assert.ok(!allItems.some((p) => p.slug === 'index'), 'index should not appear in sidebar');
  });
});
