# Feature: Spec

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fspec) — graph, discussions, approvals

**Status:** In Progress

## Summary

`specscore spec` commands validate and query the specification tree as a whole. Today the only subcommand is `spec lint`, which runs the full rule suite against a project.

## Problem

A specification tree is only as useful as its structural consistency. Without a central linter, conventions drift — some features ship without OQ sections, some Ideas live in the wrong directory, some adherence footers get stripped by copy-paste. A dedicated group for tree-wide validation gives authors and CI a single entry point.

## Contents

| Directory | Description |
|---|---|
| [lint/](lint/README.md) | Validate the spec tree against structural conventions |

### lint

Scans the spec tree and reports violations of structural conventions (README presence, Outstanding Questions sections, heading levels, feature references, internal links, index entries, adherence footers, Idea-specific rules). `--fix` applies autofixes for rules that support them.

## Behavior

### Scope of this group

Commands under `specscore spec` operate on the specification tree as a whole — not on individual features. Per-feature queries live in the [feature](../feature/README.md) group.

#### REQ: tree-level-only

`specscore spec` subcommands MUST operate at the project-tree level. Per-feature validation (e.g., lint a single feature) remains out of scope for this group; callers achieve it by running `spec lint` and filtering by path.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [CLI](../README.md) | Inherits shared exit-code contract and `--project` semantics. |

## Outstanding Questions

- Should this group grow a `spec search` command for full-text search across the spec tree, complementing the structural queries in the `feature` group?
- Should `spec lint` also live as `specscore lint` (top-level alias), matching the convention in many single-command tools (`go vet`, `go fmt`)? Or does the `spec` scope justify the prefix?

---
*This document follows the https://specscore.md/feature-specification*
