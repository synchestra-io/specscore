# Feature: Task (CLI)

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ftask) — graph, discussions, approvals

**Status:** In Progress

## Summary

`specscore task` commands read and create entries on the project task board. The MVP surface covers listing, inspecting, and creating tasks. Status transitions (claim, release, status updates) are intentionally out of scope for this group today.

## Problem

Task boards live as Markdown (`tasks/README.md` plus `tasks/<slug>/README.md` per entry) so they are diff-able and reviewable. Without command-line read and create surfaces, every consumer re-parses the Markdown, and every agent writes the files by hand. A minimal CLI — list, info, new — unlocks automation while keeping the human-readable format as the source of truth.

## Contents

| Directory | Description |
|---|---|
| [info/](info/README.md) | Show detailed task metadata for one slug |
| [list/](list/README.md) | List tasks from the board, optionally filtered by status |
| [new/](new/README.md) | Create a new task in `planning` status |

### info

Reads `tasks/<slug>/README.md` and the task's row on the board. Returns YAML or JSON with title, status, description, dependencies, and summary.

### list

Reads the `tasks/README.md` board. Returns all rows or those matching `--status`. YAML is the default format; `--format md` re-emits the board table unchanged for round-tripping.

### new

Writes a new `tasks/<slug>/README.md` and appends the row to the board. New tasks are always created with status `planning`.

## Behavior

### Scope of this group

Today the group covers the **read and seed** operations on the task board — no lifecycle transitions. That split is deliberate: transition semantics (who can claim, how conflicts are resolved, when status becomes terminal) warrant their own feature spec.

#### REQ: no-lifecycle-in-mvp

No subcommand in this group may mutate an existing task's status field. `new` only creates tasks in `planning`; it MUST NOT accept a `--status` argument for other values. Future lifecycle commands (e.g., `task claim`, `task status`) land under new subcommands with their own feature specs.

### Task slug argument

Every non-listing command takes a task slug via `--task <slug>` (a flag, not positional, matching the current implementation).

#### REQ: task-flag-required

`task info` and `task new` MUST require `--task`. A missing slug MUST exit `2` (InvalidArgs).

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [task](../../task/README.md) | Source of truth for task file structure, board format, and allowed status values. |
| [CLI](../README.md) | Inherits shared exit-code contract, `--project`, `--format`. |

## Outstanding Questions

- When should lifecycle commands (`task claim`, `task release`, `task status`) land? They depend on answering how concurrency and multi-agent claim semantics work in a git-backed board — which is a feature spec in its own right.
- Should `--task` move from a flag to a positional argument (`specscore task info <slug>`) for consistency with `feature info <id>`? Today the two groups diverge on this point.

---
*This document follows the https://specscore.md/feature-specification*
