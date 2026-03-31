import { describe, it, before } from 'node:test';
import assert from 'node:assert/strict';
import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import { renderMermaidBlocks } from './render-mermaid.js';

const execFileAsync = promisify(execFile);

// Check if mmdc/Chrome can actually run
let mmdcAvailable = false;
before(async () => {
  const mmdc = new URL('../node_modules/.bin/mmdc', import.meta.url).pathname;
  try {
    await execFileAsync(mmdc, ['--version']);
    mmdcAvailable = true;
  } catch {
    mmdcAvailable = false;
  }
});

describe('renderMermaidBlocks', () => {
  it('replaces mermaid code blocks with SVG', { skip: !mmdcAvailable && 'mmdc/Chrome not available' }, async () => {
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

  it('handles multiple mermaid blocks', { skip: !mmdcAvailable && 'mmdc/Chrome not available' }, async () => {
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

  it('wraps SVG in mermaid-wrap container', { skip: !mmdcAvailable && 'mmdc/Chrome not available' }, async () => {
    const md = ['```mermaid', 'graph LR', '    A --> B', '```'].join('\n');
    const result = await renderMermaidBlocks(md);
    assert.ok(result.includes('class="mermaid-wrap"'), 'should have mermaid-wrap wrapper');
    assert.ok(result.includes('class="mermaid-label"'), 'should have mermaid-label');
    assert.ok(result.includes('class="mermaid-body"'), 'should have mermaid-body');
  });
});
