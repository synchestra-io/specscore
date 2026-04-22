# Feature: Feature Refs

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Ffeature%2Frefs) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore feature refs <feature_id>` returns the inverse of `feature deps`: every feature in the project that names the target as a dependency in its own `## Dependencies` section. `--transitive` follows the chain.

## Synopsis

```
specscore feature refs <feature_id> [--transitive] [--fields <names>] [--format <yaml|json|text>] [--project <path>]
```

## Problem

Impact analysis and refactor safety need the reverse graph: "if I change feature X, what else depends on it?" Without a dedicated command, the only option is to grep every feature's `## Dependencies` section — fragile and slow on large repos.

## Behavior

### Direct references

Default output lists features whose `## Dependencies` section contains the target feature ID.

#### REQ: scans-all-features

The command MUST scan every feature README in the project and report those whose `## Dependencies` section lists the target. Features without a `## Dependencies` section are simply not candidates — they never appear in the output.

### Transitive references

With `--transitive`, the command MUST walk the reverse graph recursively. Dedup and cycle safety MUST match the rules in [feature deps](../deps/README.md#req-transitive-no-duplicates).

#### REQ: symmetric-transitive-rules

Transitive refs MUST apply the same deduplication and cycle-safety guarantees as transitive deps. Consumers can rely on the two commands having symmetric behavior when traversing either direction of the graph.

### Empty result

A target with no references MUST exit `0` with an empty list — not `3`. "Nothing references this" is a valid answer.

## Parameters

| Name | Required | Description |
|---|---|---|
| `feature_id` | Yes | Feature whose referrers to find. |

## Exit codes

| Code | Condition |
|---|---|
| `0` | References listed (even if empty) |
| `2` | Missing `feature_id`, invalid flag value |
| `3` | `feature_id` not found |
| `10` | Unexpected I/O failure while scanning features |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [feature](../../../feature/README.md) | Defines the `## Dependencies` section convention this command scans. |
| [deps](../deps/README.md) | Forward edge in the dependency graph; `refs` is the reverse. |

## Acceptance Criteria

### AC: direct-refs-listed

**Requirements:** cli/feature/refs#req:scans-all-features

`specscore feature refs <f>` in a project where features `x` and `y` both list `<f>` under `## Dependencies` returns exactly `x` and `y`.

### AC: transitive-matches-deps

**Requirements:** cli/feature/refs#req:symmetric-transitive-rules

For any feature `f` in a project, every feature returned by `refs f --transitive` lists `f` somewhere in its own transitive deps. The two commands agree on reachability.

### AC: empty-refs-is-success

**Requirements:** cli/feature/refs#req:scans-all-features

`refs <f>` for a feature no one references exits `0` with empty output — not `3`.

## Outstanding Questions

- Should `refs` also scan Plans and Ideas for references (via their own "affects features" metadata), or is feature-to-feature scope sufficient for the MVP?

---
*This document follows the https://specscore.md/feature-specification*
