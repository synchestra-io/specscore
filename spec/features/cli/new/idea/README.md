# Feature: New Idea

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fnew%2Fidea) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore new idea <slug>` scaffolds a lint-clean Idea artifact at `spec/ideas/<slug>.md`. Each required section is emitted with an HTML-comment prompt describing what belongs there. Content can be supplied via flags or gathered interactively with `-i`.

## Synopsis

```
specscore new idea <slug> [--title <text>] [--owner <id>] [--hmw <text>] [--context <text>] [--recommended-direction <text>] [--mvp <text>] [--not-doing "<thing> — <reason>" ...] [--interactive] [--force] [--project <path>]
```

## Problem

Ideas have a strict required-sections contract defined by the [idea](../../../idea/README.md) feature. Hand-writing them means remembering HMW framing, `Not Doing` entries, owner attribution, and status rules. The scaffolder produces a skeleton that is lint-clean on creation, so authors edit content rather than fight structure.

## Behavior

### Slug argument

The `<slug>` positional argument becomes the file name (`spec/ideas/<slug>.md`) and the Idea's declared `id`.

#### REQ: slug-required

`<slug>` MUST be supplied. Absence MUST exit `2` (InvalidArgs).

#### REQ: slug-format

The slug MUST satisfy the Idea slug-format rule: lowercase, hyphen-separated, URL-safe characters only. Invalid slugs MUST exit `2` with a message naming the offending slug.

### Content sources

Content for required fields can come from flags, interactive stdin prompts, or fallback defaults.

| Flag | Target section / field |
|---|---|
| `--title` | Idea title (defaults to title-cased slug) |
| `--owner` | Owner / author (defaults to `$USER`) |
| `--hmw` | Problem Statement sentence |
| `--context` | Context section body |
| `--recommended-direction` | Recommended Direction body |
| `--mvp` | MVP Scope body |
| `--not-doing` | Repeatable `Not Doing` entry (format: `"<thing> — <reason>"`) |

#### REQ: interactive-mode

`-i` / `--interactive` MUST prompt for each field on stdin, using the flag value as the default when supplied. Non-interactive runs use flag values plus defaults without prompting.

#### REQ: sensible-defaults

When a field has no flag value and the command is not interactive, the scaffolder MUST emit an HTML-comment prompt (`<!-- TODO: ... -->`) in place of the missing content. The resulting file MUST still be lint-clean — prompts are valid Markdown and do not break required-sections checks.

### Overwrite behavior

By default, the command MUST refuse to overwrite an existing file. `--force` opts in to overwrite.

#### REQ: no-clobber-default

If `spec/ideas/<slug>.md` already exists, the command MUST exit `1` (Conflict) with a message naming the path, unless `--force` is supplied. No partial write may occur before the collision check.

## Parameters

| Name | Required | Description |
|---|---|---|
| `slug` | Yes | Idea slug — becomes the file name. |

## Exit codes

| Code | Condition |
|---|---|
| `0` | Idea file created |
| `1` | File already exists and `--force` not supplied |
| `2` | Missing or invalid `slug`, invalid flag value |
| `10` | Unexpected I/O failure while writing |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [idea](../../../idea/README.md) | Source of truth for required Idea sections, slug rules, and lint checks. |
| [new](../README.md) | Parent group. Inherits the lint-clean-on-create guarantee. |

## Acceptance Criteria

### AC: scaffolded-idea-is-lint-clean

**Requirements:** cli/new/idea#req:sensible-defaults

`specscore new idea my-idea` with no other flags creates `spec/ideas/my-idea.md`. `specscore spec lint` immediately afterwards reports no new violations, even though several fields contain `<!-- TODO: ... -->` prompts.

### AC: existing-file-conflict

**Requirements:** cli/new/idea#req:no-clobber-default

Running the command twice for the same slug, without `--force`, exits `1` on the second run and leaves the existing file untouched. With `--force`, the second run overwrites and exits `0`.

### AC: invalid-slug-rejected

**Requirements:** cli/new/idea#req:slug-format

`specscore new idea My_Idea` exits `2` with a message that the slug contains invalid characters. No file is created.

## Outstanding Questions

- Should `new idea` accept an `--edit` flag that opens the scaffolded file in `$EDITOR` after writing, for a smoother author flow?
- Should the `--owner` default be GitHub username (from git config / `gh auth status`) instead of `$USER`, matching how `gh` attributes authorship?

---
*This document follows the https://specscore.md/feature-specification*
