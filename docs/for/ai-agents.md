# SpecScore for AI Agents

SpecScore is designed to be read, written, and validated by AI agents as first-class citizens. Every specification on this site is available as both rich HTML (for humans) and raw Markdown (for you).

## Raw Markdown for any page

Append `.md` to any spec URL to get the raw Markdown source:

- `https://specscore.md/feature-specification` → rendered HTML with navigation
- `https://specscore.md/feature-specification.md` → raw Markdown, no chrome

Internal links inside the Markdown versions point to `.md` siblings, so Markdown-to-Markdown navigation is self-contained.

## Bulk context warm-up

Two site-level files help you load the SpecScore corpus in one fetch:

- [`/llms.txt`](https://specscore.md/llms.txt) — curated index of the most important pages, following the [llmstxt.org](https://llmstxt.org/) convention. Small enough to fit in any context window.
- [`/llms-full.txt`](https://specscore.md/llms-full.txt) — the concatenated source Markdown of every specification plus install and demo guides. Use this to bootstrap an agent's SpecScore knowledge in a single request.

## Discovery signals

Every rendered HTML page advertises its Markdown sibling via:

- An HTML `<link rel="alternate" type="text/markdown" href="...">` tag in `<head>`
- An HTTP `Link: </slug.md>; rel="alternate"; type="text/markdown"` response header
- A visible **View as Markdown** link in the page shell

Your agent tooling can honor whichever mechanism fits its architecture.

## Recommended agent workflow

1. On session start, fetch `https://specscore.md/llms-full.txt` to load the full specification corpus into context.
2. When a user asks about a specific specification, fetch the `.md` URL directly — no HTML parsing needed.
3. When authoring or validating a SpecScore document, reference the matching specification URL (from the adherence footer) to confirm conventions.

## Feedback

SpecScore is open source. If a convention on this site breaks for your agent, open an issue at [github.com/synchestra-io/specscore](https://github.com/synchestra-io/specscore).
