import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { parseBlogPost, buildBlogIndex } from './build-blog.js';

describe('parseBlogPost', () => {
  it('extracts frontmatter and body from a blog post', () => {
    const raw = `---
title: Why Specs Matter
description: A look at spec-driven development
date: 2026-04-14
---

# Why Specs Matter

Content here.`;

    const post = parseBlogPost(raw, '2026-04-14-why-specs-matter.md');
    assert.equal(post.title, 'Why Specs Matter');
    assert.equal(post.description, 'A look at spec-driven development');
    assert.equal(post.date, '2026-04-14');
    assert.equal(post.slug, 'blog/why-specs-matter');
    assert.ok(post.body.includes('# Why Specs Matter'));
    assert.ok(!post.body.includes('title:'));
  });

  it('throws on missing frontmatter', () => {
    assert.throws(
      () => parseBlogPost('# No frontmatter', 'bad.md'),
      /frontmatter/i
    );
  });
});

describe('buildBlogIndex', () => {
  it('sorts posts newest first', () => {
    const posts = [
      { title: 'Old', date: '2026-01-01', slug: 'blog/old', description: 'Old post' },
      { title: 'New', date: '2026-04-14', slug: 'blog/new', description: 'New post' },
    ];
    const sorted = buildBlogIndex(posts);
    assert.equal(sorted[0].title, 'New');
    assert.equal(sorted[1].title, 'Old');
  });
});
