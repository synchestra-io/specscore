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
  const template =
    '<title>{{title}} - SpecScore</title>' +
    '<aside>{{sidebar}}</aside>' +
    '{{eyebrow}}' +
    '<main>{{content}}</main>' +
    '{{viewMarkdown}}';

  const sidebarGroups = [
    {
      label: 'Getting Started',
      items: [
        { slug: 'specifications', navLabel: 'Specifications', navOrder: 1 },
      ],
    },
    {
      label: 'Feature Specs',
      items: [
        { slug: 'feature-specification', navLabel: 'Feature', navOrder: 2 },
        { slug: 'source-references-specification', navLabel: 'Source References', navOrder: 6 },
      ],
    },
  ];

  it('replaces title and content placeholders', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '<h1>Feature</h1>',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('<title>Feature Specification - SpecScore</title>'));
    assert.ok(result.includes('<h1>Feature</h1>'));
  });

  it('generates sidebar HTML with group labels and links', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('Getting Started'));
    assert.ok(result.includes('Feature Specs'));
    assert.ok(result.includes('href="/specifications"'));
    assert.ok(result.includes('href="/feature-specification"'));
  });

  it('marks the current page as active in sidebar', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('class="active"'));
    // The active link should be the feature-specification one
    assert.match(result, /href="\/feature-specification"[^>]*class="active"/);
  });

  it('renders eyebrow label when provided', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('class="page-eyebrow"'));
    assert.ok(result.includes('Feature Specs'));
  });

  it('omits eyebrow div when eyebrow is empty', () => {
    const result = injectIntoTemplate(template, {
      title: 'SpecScore',
      content: '',
      slug: 'index',
      sidebarGroups,
      eyebrow: '',
    });

    assert.ok(!result.includes('class="page-eyebrow"'));
  });

  it('includes view-markdown link by default', () => {
    const result = injectIntoTemplate(template, {
      title: 'Feature Specification',
      content: '',
      slug: 'feature-specification',
      sidebarGroups,
      eyebrow: 'Feature Specs',
    });

    assert.ok(result.includes('href="/feature-specification.md"'));
    assert.ok(result.includes('View as Markdown'));
  });

  it('omits view-markdown link when showViewMarkdown is false', () => {
    const result = injectIntoTemplate(template, {
      title: 'SpecScore',
      content: '',
      slug: 'index',
      sidebarGroups,
      eyebrow: '',
      showViewMarkdown: false,
    });

    assert.ok(!result.includes('View as Markdown'));
  });
});
