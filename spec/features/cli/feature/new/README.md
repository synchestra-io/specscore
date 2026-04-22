# Feature: Feature New

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature%2Fnew) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature new` scaffolds a new feature directory containing a README with every required section. It optionally stages a git commit (`--commit`) or commits and pushes atomically (`--push`) so feature creation can be wired directly into agent workflows.

## Synopsis

```
specscore feature new --title <title> [--slug <slug>] [--parent <id>] [--status <status>] [--description <text>] [--depends-on <ids>] [--commit | --push] [--format <yaml|json|text>] [--project <path>]
```

## Problem

Creating a new feature by hand means remembering every required section, the correct status values, the Dependencies subsection format, the hub view link template, and the adherence footer. Every manual authoring slip becomes a lint violation. A scaffolding command that always emits a lint-clean feature removes that tax.

## Behavior

### Required inputs

`--title` is the only required flag. All other metadata has defaults.

#### REQ: title-required

`--title` MUST be supplied. Absence MUST exit `2` (InvalidArgs) with a message naming the missing flag.

#### REQ: slug-derivation

When `--slug` is omitted, the slug MUST be derived from `--title` by lowercasing, replacing whitespace and underscores with hyphens, and dropping characters outside `[a-z0-9-]`. The derived slug MUST satisfy [feature#req:slug-format](../../../feature/README.md#req-slug-format).

### Output artifact

The command creates `spec/features/<parent>/<slug>/README.md` containing every section required by [feature#req:required-sections](../../../feature/README.md#req-required-sections). When `--parent` is supplied, the new feature is nested under the named parent; otherwise it is created at the top level.

#### REQ: lint-clean

The generated README MUST pass `specscore spec lint` immediately — including the adherence footer, OQ section, hub view link (when `hub.host` is configured), and all structural conventions.

#### REQ: status-default

`--status` MUST accept one of `Draft`, `In Progress`, `Stable`, `Deprecated`. Default is `Draft`. Any other value MUST exit `2`.

### Dependencies

`--depends-on` accepts a comma-separated list of feature IDs and writes them into the new feature's `## Dependencies` section.

#### REQ: depends-on-validation

Every value supplied to `--depends-on` MUST resolve to an existing feature in the project. If any ID does not resolve, the command MUST exit `3` (NotFound) BEFORE any file is written.

### Git integration

Two mutually exclusive modes commit the scaffolded files:

| Flag | Behavior |
|---|---|
| `--commit` | Stage and commit the new feature on the current branch. Do not push. |
| `--push` | Imply `--commit`; additionally push to the remote. Atomic — a push failure rolls back no local state, but the command reports `1` (Conflict) so callers know to reconcile. |

#### REQ: git-requires-repo

`--commit` and `--push` MUST fail with exit `10` if the project is not a git repository. The working tree MUST NOT be mutated when the preflight check fails.

#### REQ: push-conflict-exit-code

A `--push` that fails because the remote has diverged MUST exit `1` (Conflict). A generic push failure (network, auth) MUST exit `10`.

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | Feature scaffolded (and optionally committed/pushed) |
| `1` | `--push` failed due to remote conflict |
| `2` | Missing `--title`, invalid flag value, bad slug, collision with existing feature |
| `3` | `--parent` or `--depends-on` names a non-existent feature |
| `10` | Not a git repository (with `--commit`/`--push`), or unexpected I/O / git failure |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../../feature/README.md) | Source of truth for required sections, slug format, and status values. `feature new` emits a README that conforms. |
| [adherence-footer](../../../adherence-footer/README.md) | The generated README's footer uses the canonical `feature-specification` URL. |
| [CLI Feature group](../README.md) | Inherits `--format`, `--project`. |

## Acceptance Criteria

### AC: new-feature-is-lint-clean

**Requirements:** cli/feature/new#req:lint-clean

`specscore feature new --title "My Feature"` creates `spec/features/my-feature/README.md` containing every required section. `specscore spec lint` immediately afterwards reports no new violations.

### AC: missing-title-exits-2

**Requirements:** cli/feature/new#req:title-required

`specscore feature new` with no `--title` exits `2` with a message naming the missing flag. No directory is created.

### AC: missing-parent-exits-3

**Requirements:** cli/feature/new#req:depends-on-validation

`specscore feature new --title X --parent does-not-exist` exits `3` with a message naming the missing parent. No directory is created.

### AC: push-conflict-returns-1

**Requirements:** cli/feature/new#req:push-conflict-exit-code

When `--push` is used against a remote that has diverged from the local branch, the command exits `1` (Conflict). The local scaffolded files remain so the author can reconcile manually.

## Outstanding Questions

- Should `feature new` accept a `--dry-run` flag that prints the intended paths and README content to stdout without writing, so agents can preview before committing?
- Should the scaffolded README's `## Source Ideas` field be auto-populated when `--from-idea <slug>` is supplied, tightening the Idea → Feature promotion flow?

---
*This document follows the https://specscore.md/feature-specification*
