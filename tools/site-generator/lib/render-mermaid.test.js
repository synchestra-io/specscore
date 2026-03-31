import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { renderMermaidBlocks } from './render-mermaid.js';

describe('renderMermaidBlocks', () => {
  it('replaces mermaid code blocks with SVG', async () => {
    const md = [
      'Some text.',
      '',
      '```mermaid',
      'graph LR',
      '    A --> B',
      '```',
      '',
      'More text.',
    ].join('\n');

    const result = await renderMermaidBlocks(md);

    // Mermaid block should be replaced with an SVG
    assert.ok(!result.includes('```mermaid'), 'mermaid code block should be removed');
    assert.ok(result.includes('<svg'), 'should contain inline SVG');
    assert.ok(result.includes('Some text.'), 'surrounding text preserved');
    assert.ok(result.includes('More text.'), 'surrounding text preserved');
  });

  it('returns markdown unchanged when no mermaid blocks', async () => {
    const md = 'Just regular markdown.\n\n```js\nconsole.log("hi");\n```';
    const result = await renderMermaidBlocks(md);
    assert.equal(result, md);
  });

  it('handles multiple mermaid blocks', async () => {
    const md = [
      '```mermaid',
      'graph LR',
      '    A --> B',
      '```',
      '',
      'Between.',
      '',
      '```mermaid',
      'graph TD',
      '    C --> D',
      '```',
    ].join('\n');

    const result = await renderMermaidBlocks(md);
    const svgCount = (result.match(/<svg/g) || []).length;
    assert.equal(svgCount, 2, 'should have 2 SVGs');
  });
});
