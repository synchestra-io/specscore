# Feature: Feature Info

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature%2Finfo) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature info <feature_id>` returns structured metadata for a single feature: status, parent, children, dependency counts, and a section table-of-contents with line ranges.

## Synopsis

```
specscore feature info <feature_id> [--format <yaml|json|text>] [--project <path>]
```

## Problem

Tools that operate on specific features — LSP servers, editor plugins, agents stepping through a tree — need a stable, machine-readable snapshot of a feature's metadata and section layout. Parsing `README.md` on every query is fragile and slow. A single command that returns everything a consumer needs, in a known schema, removes that burden.

## Behavior

### Output shape

Output is a YAML (default) / JSON / text document describing the feature.

#### REQ: required-fields

The output MUST include, at minimum:

- `path` — the canonical feature ID
- `status` — the feature's declared status
- `deps` — list of feature IDs the feature depends on (empty list if none)
- `refs` — list of feature IDs that reference this feature (empty list if none)
- `sections` — an ordered list of top-level sections from the README, each with `title`, `lines` (range), and nested `children` or `items` counts as appropriate

Additional fields MAY be added in later releases; consumers MUST tolerate unknown fields (see [parent feature#req:stable-yaml-keys](../../README.md#req-stable-yaml-keys)).

### Feature resolution

The `<feature_id>` argument MUST match the path-identification rules from the [feature](../../../feature/README.md) spec.

#### REQ: not-found

An unresolved `<feature_id>` MUST exit `3` (NotFound) with a message that names the requested ID.

### Format selection

`--format` MUST accept `yaml` (default), `json`, or `text`. Text output is a condensed, human-readable rendering — suitable for quick CLI reads, not for parsing.

## Parameters

| Name | Required | Description |
|---|---|---|
| `feature_id` | Yes | Feature to inspect, using the path-identification rules (e.g., `cli/version`). |

## Exit codes

| Code | Condition |
|---|---|
| `0` | Feature found and info printed |
| `2` | Missing `feature_id` argument, invalid `--format` value |
| `3` | Feature not found |
| `10` | Unexpected I/O failure while reading the feature |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../../feature/README.md) | Defines the structural rules (required sections, REQ subsections) that the section TOC reflects. |
| [CLI Feature group](../README.md) | Inherits `--format`, `--project`. |

## Acceptance Criteria

### AC: info-returns-sections

**Requirements:** cli/feature/info#req:required-fields

`specscore feature info cli/version` returns a YAML document containing `path: cli/version`, a valid `status`, and a `sections` list whose titles match the feature's README top-level headings in order.

### AC: not-found-exits-3

**Requirements:** cli/feature/info#req:not-found

`specscore feature info does-not-exist` exits `3` with a stderr message naming the missing ID. No partial output is written to stdout.

## Outstanding Questions

- Should `feature info` optionally include the section *bodies* (with `--fields body`) for consumers that want the full spec in one call, or should body reads stay separate via `cat`?

---
*This document follows the https://specscore.md/feature-specification*
