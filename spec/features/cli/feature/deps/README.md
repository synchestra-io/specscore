# Feature: Feature Deps

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature%2Fdeps) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature deps <feature_id>` reads the `## Dependencies` section of the named feature's README and lists the features it depends on. `--transitive` follows the chain recursively.

## Synopsis

```
specscore feature deps <feature_id> [--transitive] [--fields <names>] [--format <yaml|json|text>] [--project <path>]
```

## Problem

Feature dependencies are declared in free-form Markdown but consumed by tooling that builds dependency graphs (planning, impact analysis, refactor safety). A structured query that returns direct or transitive dependencies — without the caller parsing Markdown — is the cheapest integration point.

## Behavior

### Direct dependencies

The default output lists only the features explicitly listed in the target feature's `## Dependencies` section.

#### REQ: reads-dependencies-section

Dependencies MUST be read from the `## Dependencies` section of `spec/features/<feature_id>/README.md`. Features without a `## Dependencies` section MUST return an empty list (exit `0`), not an error.

### Transitive dependencies

With `--transitive`, the command MUST walk the dependency graph recursively and return every feature reachable from the target via dependency edges.

#### REQ: transitive-no-duplicates

Transitive output MUST deduplicate. A feature reachable via multiple paths MUST appear exactly once in the result.

#### REQ: transitive-cycle-safe

If the dependency graph contains a cycle, transitive resolution MUST terminate (no infinite loop) and MUST NOT error. Cycles are a spec-authoring smell; lint is responsible for flagging them.

### Output format

Default text output lists one feature ID per line. YAML / JSON output is a structured list matching the shape defined in [feature list](../list/README.md).

## Parameters

| Name | Required | Description |
|---|---|---|
| `feature_id` | Yes | Feature whose dependencies to resolve. |

## Exit codes

| Code | Condition |
|---|---|
| `0` | Dependencies listed (even if empty) |
| `2` | Missing `feature_id`, invalid flag value |
| `3` | `feature_id` not found |
| `10` | Unexpected I/O failure |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../../feature/README.md) | Defines the `## Dependencies` section convention this command reads. |
| [refs](../refs/README.md) | Inverse query — which features reference the given feature. |

## Acceptance Criteria

### AC: direct-deps-listed

**Requirements:** cli/feature/deps#req:reads-dependencies-section

`specscore feature deps <f>` for a feature whose README lists `- a` and `- b` under `## Dependencies` returns exactly `a` and `b` (and nothing else) in the default text output.

### AC: transitive-dedup-and-terminate

**Requirements:** cli/feature/deps#req:transitive-no-duplicates, cli/feature/deps#req:transitive-cycle-safe

`deps --transitive` on a graph where `a → b, a → c, b → c` returns `b` and `c` exactly once. On a cycle `a → b → a`, the command returns `b` (and `a` if the implementation reports the root on self-cycle) without hanging.

## Outstanding Questions

- Should `--transitive` have a depth limit (`--depth 2`) for large graphs, or is the cycle-safe full walk always acceptable?

---
*This document follows the https://specscore.md/feature-specification*
