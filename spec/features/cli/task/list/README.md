# Feature: Task List

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ftask%2Flist) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore task list` reads the project's task board (`tasks/README.md`) and prints every task row. `--status <value>` filters to one status. Default output is YAML; `--format md` re-emits the board table unchanged for round-tripping.

## Synopsis

```
specscore task list [--status <status>] [--format <yaml|json|md>] [--project <path>]
```

## Problem

Agents and scripts need to know what tasks exist, in which status. Parsing `tasks/README.md` themselves means re-implementing the board format. A command that returns structured data (or re-emits the original Markdown) keeps the board as source of truth while unlocking automation.

## Behavior

### Board source

The command MUST read the project's canonical board file at `tasks/README.md` under the project root.

#### REQ: reads-board

Output rows MUST reflect the board file exactly — not be recomputed from individual `tasks/<slug>/README.md` files. If a `<slug>/README.md` exists but its slug is not on the board, the task is not listed.

### Status filter

`--status` accepts any of the values declared by the [task](../../../task/README.md) feature.

#### REQ: status-filter

`--status <value>` MUST return only rows whose status column equals `<value>`. Unknown status values MUST exit `2` (InvalidArgs) naming the offending value. An empty `--status` (e.g., `--status ""`) MUST be treated as unset.

### Output formats

| Format | Shape |
|---|---|
| `yaml` (default) | List of objects, each with `slug`, `title`, `status`, `description` (if present). |
| `json` | Same shape, JSON-encoded. |
| `md` | Re-emit the original board table verbatim (possibly filtered by status). |

#### REQ: md-format-roundtrip

`--format md` MUST emit a Markdown table whose header and column layout match `tasks/README.md`. When no filter is applied, the output is byte-for-byte equivalent to the source file (allowing for a trailing newline normalization).

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | Listing printed (even if the board is empty) |
| `2` | Invalid `--status` or `--format` value |
| `3` | No project found, or project has no `tasks/README.md` |
| `10` | Unexpected I/O failure |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [task](../../../task/README.md) | Defines the board format and allowed status values. |
| [CLI Task group](../README.md) | Inherits scope rules and shared flags. |

## Acceptance Criteria

### AC: default-returns-yaml

**Requirements:** cli/task/list#req:reads-board

`specscore task list` on a board with three tasks returns a YAML list of three objects, each with `slug`, `title`, and `status` at minimum.

### AC: md-format-roundtrips

**Requirements:** cli/task/list#req:md-format-roundtrip

`specscore task list --format md` re-emits the `tasks/README.md` table. Diffing the command output against the source file (modulo trailing newline) yields no changes.

### AC: unknown-status-rejected

**Requirements:** cli/task/list#req:status-filter

`specscore task list --status not-a-status` exits `2` with a message naming the invalid value.

## Outstanding Questions

- Should `--status` accept a comma-separated list for multi-status filters (e.g., `--status queued,in_progress`), or is single-status filtering plus `yq` post-filtering adequate?

---
*This document follows the https://specscore.md/feature-specification*
