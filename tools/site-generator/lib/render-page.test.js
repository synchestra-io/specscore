import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { renderMarkdownToHtml, injectIntoTemplate } from './render-page.js';

describe('renderMarkdownToHtml', () => {
  it('renders basic markdown to HTML', () => {
    const html = renderMarkdownToHtml('# Hello\n\nA paragraph.');
    assert.ok(html.includes('<h1>Hello</h1>'));
    assert.ok(html.includes('<p>A paragraph.</p>'));
  });

  it('renders GFM tables', () => {
    const md = '| A | B |\n|---|---|\n| 1 | 2 |';
    const html = renderMarkdownToHtml(md);
    assert.ok(html.includes('<table>'));
    assert.ok(html.includes('<td>1</td>'));
  });

  it('renders fenced code blocks with language class', () => {
    const md = '```js\nconsole.log("hi");\n```';
    const html = renderMarkdownToHtml(md);
    assert.ok(html.includes('<code class="language-js">'));
  });
});

describe('injectIntoTemplate', () => {
  const template = '<title>{{title}} - SpecScore</title>{{nav}}<main>{{content}}</main><a href="{{mdUrl}}">md</a>';

  const navItems = [
    { slug: 'index', navLabel: 'Home' },
    { slug: 'specifications', navLabel: 'Specifications' },
  ];

  it('replaces all placeholders', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '<h1>Feature</h1>',
      slug: 'feature-specification',
      navItems,
    });

    assert.ok(result.includes('<title>Feature Specification - SpecScore</title>'));
    assert.ok(result.includes('<h1>Feature</h1>'));
    assert.ok(result.includes('href="/feature-specification.md"'));
    assert.ok(result.includes('href="/"'));
    assert.ok(result.includes('href="/specifications"'));
  });

  it('marks the current page as active in nav', () => {
    const result = injectIntoTemplate(template, {
      title: 'Specifications',
      content: '<p>List</p>',
      slug: 'specifications',
      navItems,
    });

    assert.ok(result.includes('class="active"'));
  });
});
