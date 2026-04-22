# Feature: Task New

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ftask%2Fnew) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore task new` creates a new task: writes `tasks/<slug>/README.md` with the required sections and appends a row to the task board at `tasks/README.md`. New tasks are always created in `planning` status.

## Synopsis

```
specscore task new --task <slug> --title <text> [--description <text>] [--depends-on <slugs>] [--format <yaml|json>] [--project <path>]
```

## Problem

Board and task-file creation need to stay in lock-step: a row without a file is an orphan, a file without a row is invisible. Doing both by hand is a two-step ritual that is easy to half-complete. A single command that writes both atomically keeps the board coherent.

## Behavior

### Required inputs

`--task` and `--title` are the only required flags.

#### REQ: slug-and-title-required

`--task` and `--title` MUST both be supplied. Missing either MUST exit `2` (InvalidArgs) with a message naming the missing flag(s).

#### REQ: slug-format

The value of `--task` MUST satisfy the task slug-format rule: lowercase, hyphen-separated, URL-safe. Invalid slugs MUST exit `2`.

### Created artifacts

The command produces exactly two changes in the working tree:

1. A new file `tasks/<slug>/README.md` with required sections populated.
2. A new row on `tasks/README.md` referencing the slug with status `planning`.

#### REQ: planning-only

`task new` MUST set status to `planning`. Creating tasks in other statuses is out of scope (see [parent cli/task#req:no-lifecycle-in-mvp](../README.md#req-no-lifecycle-in-mvp)). Callers cannot override this via a flag.

#### REQ: atomic-board-update

The new file and the new board row MUST be written as a pair. If either write fails, the command MUST leave the working tree unchanged (no orphan files, no orphan rows). The implementation MAY achieve atomicity by writing to a temp path and renaming, or by rolling back the file write on board-update failure.

### Collision handling

`tasks/<slug>/` existing is a hard conflict. The command MUST NOT overwrite.

#### REQ: no-clobber

If `tasks/<slug>/` already exists OR if the board already has a row with that slug, the command MUST exit `1` (Conflict) with a message naming the collision. No `--force` exists for `task new`.

### Dependencies

`--depends-on` accepts a comma-separated list of task slugs that the new task depends on.

#### REQ: depends-on-validation

Every value supplied to `--depends-on` MUST resolve to an existing task (either in the board or as a directory under `tasks/`). Unresolved slugs MUST exit `3` (NotFound) BEFORE any file is written.

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | Task created (file + board row written) |
| `1` | Slug collision — directory or board row already exists |
| `2` | Missing `--task`/`--title`, invalid flag value, bad slug |
| `3` | `--depends-on` names a non-existent task |
| `10` | Unexpected I/O failure (partial write scenario — the atomicity guarantee requires rollback if this occurs) |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [task](../../../task/README.md) | Defines required task file sections, board format, and allowed status values. `task new` emits a task that conforms. |
| [CLI Task group](../README.md) | Inherits the MVP "no lifecycle" scope rule. |

## Acceptance Criteria

### AC: creates-file-and-board-row

**Requirements:** cli/task/new#req:slug-and-title-required, cli/task/new#req:atomic-board-update

`specscore task new --task my-task --title "My Task"` creates `tasks/my-task/README.md` with the required sections and appends a row to `tasks/README.md` with status `planning`. Both files are written atomically.

### AC: collision-exits-1

**Requirements:** cli/task/new#req:no-clobber

Running the command twice with the same slug exits `1` on the second run. No state is mutated on the second run.

### AC: missing-dep-exits-3

**Requirements:** cli/task/new#req:depends-on-validation

`specscore task new --task t --title T --depends-on does-not-exist` exits `3` with a message naming the missing dependency. Neither the new file nor the new board row is created.

### AC: status-fixed-to-planning

**Requirements:** cli/task/new#req:planning-only

Every task created by this command has status `planning`. The command exposes no flag for overriding this.

## Outstanding Questions

- Should `task new` accept a `--commit` / `--push` flag pair for parity with `feature new`, enabling one-shot agent workflows that create and push a task in a single call?
- Should the board row be inserted at a sorted position (by slug or by creation date) rather than appended, to keep the board diff-friendly across concurrent creators?

---
*This document follows the https://specscore.md/feature-specification*
