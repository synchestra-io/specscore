# Feature: Feature List

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature%2Flist) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature list` returns every feature in the project as a flat, alphabetically sorted list of feature IDs. With `--fields`, it returns structured YAML (or JSON) with metadata columns per feature.

## Synopsis

```
specscore feature list [--fields <names>] [--format <yaml|json|text>] [--project <path>]
```

## Problem

Scripts, editor integrations, and shell completions need the simplest possible enumeration of what features exist. A flat ID list lets callers pipe into `grep`, `fzf`, or `xargs` without parsing. When more detail is needed, the same command returns structured data, so callers have one entry point for both shapes.

## Behavior

### Default output

Without flags, the command prints one feature ID per line, sorted alphabetically, to stdout.

#### REQ: flat-text-default

The default output format MUST be text: one feature ID per line, sorted alphabetically, no headers, no trailing blank line. This output is directly pipeable into standard Unix tools.

#### REQ: id-format

Feature IDs MUST match the path-identification rules in the [feature](../../../feature/README.md) spec: lowercase, hyphen-separated slugs joined by `/` for nested features (e.g., `cli/version`).

### Structured output

When `--fields` is supplied, the command switches to structured output (YAML by default, JSON with `--format json`) with a top-level list; each entry has the feature ID plus the requested fields.

#### REQ: fields-forces-structured

When `--fields` is non-empty, `--format text` MUST NOT be used. If the caller explicitly sets `--format text` alongside `--fields`, the command MUST auto-upgrade to YAML (for backward compatibility) rather than error.

### Filtering

The MVP does not support server-side filters. Callers filter with `grep` on the text output or `yq` on the structured output.

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | Listing printed (even if zero features exist) |
| `2` | Unknown `--fields` name, invalid `--format` value |
| `3` | No project found (no `specscore-spec-repo.yaml` reachable from `cwd` and no `--project`) |
| `10` | Unexpected I/O failure while reading the feature tree |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../../feature/README.md) | Source of truth for what counts as a feature and how IDs are formed. |
| [CLI Feature group](../README.md) | Inherits `--fields`, `--format`, `--project`. |

## Acceptance Criteria

### AC: default-listing-pipeable

**Requirements:** cli/feature/list#req:flat-text-default, cli/feature/list#req:id-format

`specscore feature list` in a repo with features `a`, `b/c`, and `d` prints exactly three lines: `a`, `b/c`, `d`, in that order, with no headers or trailing blank.

### AC: fields-returns-yaml

**Requirements:** cli/feature/list#req:fields-forces-structured

`specscore feature list --fields status` returns a YAML list where each entry has `id` and `status` keys. `--format text --fields status` upgrades to YAML rather than producing mixed text output.

## Outstanding Questions

- Should `--status <value>` and `--parent <id>` become first-class filters, or is post-filtering with `yq` / `grep` adequate for the foreseeable future?

---
*This document follows the https://specscore.md/feature-specification*
