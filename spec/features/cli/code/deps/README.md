# Feature: Code Deps

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fcode%2Fdeps) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore code deps` scans source files for `specscore:` annotations and bare SpecScore URLs in comments, then lists the SpecScore resources (features, plans, docs) those files depend on.

## Synopsis

```
specscore code deps [--path <glob>] [--type <feature|plan|doc>]
```

## Problem

Authors and CI need to answer "what specs does this code claim to implement?" without reading every comment. Without a dedicated query command, the only options are grep (fragile, misses URL-only references) or opening files one by one.

## Behavior

### Inputs

The command operates on the working tree under the project root. It reads files matching `--path` and extracts resource references from comments.

#### REQ: path-glob

`--path` MUST accept a double-star glob (e.g., `pkg/**/*.go`, `internal/cli/*.go`). The default value `**/*` matches every file in the tree.

#### REQ: type-filter

`--type` MAY be one of `feature`, `plan`, or `doc`. When set, results MUST be filtered to the given resource type. When omitted, results include all types. Any other value is a `2` (InvalidArgs) error.

### Sources matched

The scanner recognizes two forms of reference in comments:

1. `specscore:` annotations as defined by the [source-references](../../../source-references/README.md) feature.
2. Bare `https://specscore.md/...` URLs.

#### REQ: both-forms

The scanner MUST detect both forms. A reference that appears in either form MUST be reported exactly once per source file.

### Output

Output lists, per source file, the resources it depends on.

#### REQ: stable-output

Output MUST be stable for the same working tree across runs — files sorted, references within a file sorted. This makes the output safe to diff in CI.

## Parameters

None. All inputs are flags.

## Exit codes

Standard CLI exit codes (see [parent](../../README.md#shared-exit-code-contract)). The ones this command can return:

| Code | Condition |
|---|---|
| `0` | Scan completed (zero or more references reported) |
| `2` | `--type` or `--path` is invalid |
| `10` | Unexpected I/O failure while scanning |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [source-references](../../../source-references/README.md) | Defines the annotation syntax this command parses. |

## Acceptance Criteria

### AC: path-filter-works

**Requirements:** cli/code/deps#req:path-glob

Supplying `--path pkg/**/*.go` restricts the scan to Go files under `pkg/`. Files outside the glob are not opened.

### AC: type-filter-works

**Requirements:** cli/code/deps#req:type-filter

Supplying `--type feature` causes the output to contain only feature references; `plan` and `doc` references are suppressed. An invalid `--type` value exits `2`.

### AC: both-annotation-and-url-detected

**Requirements:** cli/code/deps#req:both-forms

A file containing both a `specscore:` annotation and a bare `https://specscore.md/...` URL in comments reports the referenced resources without duplication.

## Outstanding Questions

- Should `--path` accept a comma-separated list of globs for unions (e.g., `pkg/**/*.go,internal/**/*.go`) or should the current single-glob behavior stay, requiring callers to run the command twice?

---
*This document follows the https://specscore.md/feature-specification*
