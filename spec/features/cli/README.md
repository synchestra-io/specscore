# Feature: CLI

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli) — graph, discussions, approvals

**Status:** In Progress

## Summary

The `specscore` CLI is the reference tooling for working with SpecScore specification repositories. It validates specs, queries the feature tree, inspects source-to-spec links, manages tasks, scaffolds new artifacts, and reports its own identity.

This feature is the umbrella for command-level specifications. Each command or command group owns a child feature directory with its own contract — flags, output, exit codes, and behavior.

## Problem

Commands that grew inside the codebase without written specs accrete inconsistencies: output formats drift between releases, exit codes become arbitrary, flag names repeat the same idea in different words, and script authors cannot tell which behaviors are guaranteed versus incidental. A CLI that is the primary interface to a specification framework should itself be specified — both to pin its own contract and to dogfood the format on its own tooling.

## Contents

| Directory | Description |
|---|---|
| [code/](code/README.md) | Source-code → SpecScore relationship queries |
| [feature/](feature/README.md) | Feature tree queries and scaffolding |
| [new/](new/README.md) | Scaffolding for non-feature artifacts (ideas, …) |
| [spec/](spec/README.md) | Specification-tree validation and search |
| [task/](task/README.md) | Task board management |
| [version/](version/README.md) | CLI version reporting |

### code

Queries relationships from source files to SpecScore resources. Scans `specscore:` annotations and URLs embedded in source comments and reports the features, plans, or docs those files depend on. Read-only.

### feature

Queries the feature tree: list every feature, inspect a feature's metadata and section TOC, view the hierarchy as a tree, and follow dependency / reference chains. Also hosts `feature new`, which scaffolds a new feature directory with a lint-clean README.

### new

Scaffolds non-feature SpecScore artifacts. Currently hosts `new idea`, which creates a lint-clean Idea file at `spec/ideas/<slug>.md`. Future artifact kinds (plan, task, decision, …) will slot in as additional subcommands here.

### spec

Validates the specification tree. Hosts `spec lint`, which runs the full checker suite (structural conventions, adherence footers, OQ sections, index completeness, and Idea-specific rules) and optionally applies autofixes. Reports violations with severity levels.

### task

Manages the project task board at `tasks/README.md` and individual task files under `tasks/<slug>/README.md`. Supports listing, inspecting, and creating tasks. Task status transitions and claim/release semantics are not part of the MVP surface — this group is the minimum needed to read and seed a board.

### version

Reports the CLI's build identity. `specscore version` prints the full human-readable line; `specscore --version` (and `-v`) prints the bare semver for scripts. See [version/README.md](version/README.md) for the full contract.

## Behavior

### Command-naming conventions

Commands follow a `specscore <resource> <action>` pattern with singular nouns and verb subcommands, matching the style of `gh`, `kubectl`, and `docker`.

#### REQ: singular-resource-names

Resource names in command paths MUST be singular (`feature`, `task`, `idea`), never plural. The resource name identifies a *type*; pluralization is an output-shape concern, not a command-name one.

#### REQ: verb-subcommands

Every action MUST be an explicit subcommand verb (`list`, `info`, `new`, `deps`, `refs`, `tree`, `lint`). A bare resource name (e.g., `specscore feature`) MUST show help — it MUST NOT perform an implicit default action like listing.

#### REQ: prefer-new-over-create

Commands that create new artifacts MUST use the verb `new`, never `create`. This matches `gh issue new`, `gh pr new`, and synchestra's conventions.

### Shared exit-code contract

Every `specscore` command MUST observe the following exit-code contract. These codes match the constants exported by [`pkg/exitcode`](../../../pkg/exitcode), which the CLI uses uniformly.

| Exit code | Meaning |
|---|---|
| `0` | Success |
| `1` | Conflict (concurrent modification, stale read) |
| `2` | Invalid arguments (missing required flag, bad flag value, malformed input) |
| `3` | Resource not found |
| `4` | Invalid state transition |
| `10` | Unexpected / catch-all runtime error |

Exit codes `5–9` and `11–19` are reserved for future standard codes and MUST NOT be used by individual commands.

#### REQ: standard-exit-codes

Commands MUST map errors to the standard code with the matching semantics. A command that has no notion of "conflict" or "invalid state transition" simply never returns those codes; it does not repurpose them.

#### REQ: error-on-stderr

On any non-zero exit, a human-readable explanation MUST be written to stderr. stdout MUST remain free of error prose so that pipelines consuming structured output (YAML/JSON) are not corrupted by error messages.

### Output format conventions

Most read commands support `--format` for selecting between `text`, `yaml`, and `json`. Some also support `md` (task list).

#### REQ: yaml-default-for-structured

Read commands that return structured data (feature info, feature list, task info, task list, feature deps, feature refs, feature tree) MUST default to YAML output. `--format json` and `--format text` MUST be accepted as alternatives where documented on the individual command.

#### REQ: stable-yaml-keys

YAML and JSON output keys are part of the command's contract. Renaming or removing a key is a breaking change and MUST follow the deprecation path (announce in release notes, keep the old key for at least one release cycle). Adding new keys is always allowed.

### Shared flags

Several flags appear across multiple commands with identical semantics:

| Flag | Semantics |
|---|---|
| `--project` | Path to the project root. Autodetected from `cwd` (walks up until finding `specscore-spec-repo.yaml`) when omitted. |
| `--format` | Output format. Allowed values vary by command (always a subset of `yaml`, `json`, `text`, `md`). |
| `-h`, `--help` | Print help and exit `0`. Provided by cobra; commands MUST NOT override it. |

#### REQ: project-autodetect

When `--project` is not supplied, commands MUST autodetect the project root by searching upward from the current working directory for `specscore-spec-repo.yaml`. If no project is found, commands MUST exit `3` (NotFound) with a clear message.

## Outstanding Questions

- The MVP task surface (list, info, new) does not include status transitions (claim, release, status update). Should those land as part of this feature or in a future `cli/task/status/…` expansion that tracks a `task-lifecycle` feature spec?
- Several commands (`feature deps`, `feature refs`, `feature tree`, `feature list`) share a `--fields` flag with overlapping semantics. Should that flag be promoted to a shared-flag REQ in this parent feature, or stay documented per-command until the semantics fully converge?
- Commands currently do not emit a machine-readable error envelope on non-zero exit — error details go to stderr as free prose. Should stderr output for structured formats (`--format json`, `--format yaml`) also be structured (JSON/YAML error object), or is free prose + exit code sufficient for the CLI's callers?

---
*This document follows the https://specscore.md/feature-specification*
