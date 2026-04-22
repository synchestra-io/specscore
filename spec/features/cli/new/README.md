# Feature: New

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fnew) — graph, discussions, approvals

**Status:** In Progress

## Summary

`specscore new` is the scaffolding surface for SpecScore artifacts that are not features. It currently hosts `new idea`, which creates a lint-clean Idea file at `spec/ideas/<slug>.md`. Future artifact kinds (Plan, Task-template, Decision) will land under this group.

## Problem

Feature scaffolding already lives at [feature new](../feature/new/README.md) because features are the most common artifact. Other SpecScore artifacts (Ideas, Plans, Tasks, Decisions) have their own required sections and linting rules. A shared group for these non-feature scaffolders keeps the CLI surface discoverable (`specscore new <kind>`) without cluttering each type's own command group.

## Contents

| Directory | Description |
|---|---|
| [idea/](idea/README.md) | Scaffold a new Idea artifact at `spec/ideas/<slug>.md` |

### idea

Creates a lint-clean Idea skeleton with every required section, HTML-comment prompts describing what belongs in each, and either flag-supplied or interactively prompted content for the core fields (title, HMW, owner, not-doing entries).

## Behavior

### Scope of this group

Commands under `specscore new` produce a single artifact each, targeting its canonical path. They MUST NOT create features — `feature new` is the entry point for features, and this split keeps the two flows discoverable.

#### REQ: not-for-features

No subcommand of `specscore new` may create a feature directory or mutate `spec/features/`. Adding a new artifact kind here MUST use a dedicated canonical path (`spec/ideas/`, `spec/plans/`, `spec/decisions/`, etc.).

### Lint-clean output

Every artifact scaffolded under this group MUST pass `specscore spec lint` on first creation.

#### REQ: lint-clean-on-create

A file produced by any `specscore new <kind>` command MUST satisfy every lint rule applicable to its document kind. Commands that cannot produce a lint-clean artifact for the given inputs (e.g., a bad slug) MUST exit with a failure code BEFORE writing the file.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [idea](../../idea/README.md) | Defines the required shape of an Idea. `new idea` produces instances. |
| [feature new](../feature/new/README.md) | Sibling scaffolder for features. Lives in the `feature` group because features are queried there too. |
| [CLI](../README.md) | Inherits shared exit-code contract, `--project`, `--format`. |

## Outstanding Questions

- Should `new plan`, `new task`, and `new decision` land here once those artifact kinds stabilize, or should they live in their own groups (`specscore plan new`, `specscore task new`, `specscore decision new`)? `task new` currently lives under the `task` group — is that a precedent for or against consolidating here?
- Should `new` support a `--from-template <path>` flag that scaffolds from a project-local template file, enabling teams to codify conventions beyond the stock output?

---
*This document follows the https://specscore.md/feature-specification*
