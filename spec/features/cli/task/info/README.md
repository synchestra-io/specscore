# Feature: Task Info

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ftask%2Finfo) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore task info --task <slug>` returns detailed metadata for a single task: title, status, description, dependencies, and summary. Status is read from the board (`tasks/README.md`); the rest is read from `tasks/<slug>/README.md`.

## Synopsis

```
specscore task info --task <slug> [--format <yaml|json>] [--project <path>]
```

## Problem

Agents picking up work need the full context for a task in one call — title, status, dependencies, and the summary describing what done looks like. Parsing `tasks/<slug>/README.md` and cross-referencing the board is duplicated effort for every consumer. A single command with structured output removes that burden.

## Behavior

### Required input

The task slug is supplied via the `--task` flag.

#### REQ: task-flag-required

`--task` MUST be supplied. Absence MUST exit `2` (InvalidArgs).

### Data sources

The command merges two sources: the board row (for authoritative status) and the task's own README (for everything else).

#### REQ: status-from-board

The `status` field in the output MUST come from the task's row on the board (`tasks/README.md`). If the slug exists as a directory but not on the board, the command MUST exit `3` (NotFound).

#### REQ: body-from-readme

`title`, `description`, `dependencies`, and `summary` MUST come from `tasks/<slug>/README.md`. Missing sections MUST render as empty fields (not as errors) — consistent with lint, which is responsible for flagging absent required sections.

### Output format

Output is YAML by default; `--format json` switches to JSON. No text format is offered — task info is structured data consumed by agents more often than read directly.

#### REQ: format-values

`--format` MUST accept `yaml` (default) or `json`. Any other value MUST exit `2`.

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | Task found and info printed |
| `2` | Missing `--task`, invalid `--format` |
| `3` | No project found, task slug not on the board, or `tasks/<slug>/README.md` missing |
| `10` | Unexpected I/O failure |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [task](../../../task/README.md) | Defines the task file structure and board format. |
| [CLI Task group](../README.md) | Inherits scope rules and shared flags. |

## Acceptance Criteria

### AC: info-merges-board-and-readme

**Requirements:** cli/task/info#req:status-from-board, cli/task/info#req:body-from-readme

For a task whose board row shows `in_progress` and whose README has title `Build X` and two dependencies, `specscore task info --task <slug>` returns a YAML document with `status: in_progress`, `title: Build X`, and the two dependencies as a list.

### AC: missing-task-exits-3

**Requirements:** cli/task/info#req:status-from-board

`specscore task info --task does-not-exist` exits `3` with a message naming the missing slug. No partial output is written.

### AC: missing-flag-exits-2

**Requirements:** cli/task/info#req:task-flag-required

`specscore task info` (no `--task`) exits `2` with a message naming the missing flag.

## Outstanding Questions

- Should the `--task` flag become a positional argument (`specscore task info <slug>`) for consistency with `feature info <id>`? The inconsistency is tracked as an OQ on the parent group.

---
*This document follows the https://specscore.md/feature-specification*
