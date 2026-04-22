# Feature: Spec Lint

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fcli%2Fspec%2Flint) — graph, discussions, approvals

**Status:** Stable

## Summary

`specscore spec lint` scans the specification tree and reports violations of structural conventions. Violations are categorized by severity (error, warning, info). `--fix` applies autofixes for rules that support them (adherence footers, hub view links, idea sync / index / archived-order rules).

## Synopsis

```
specscore spec lint [--fix] [--severity <error|warning|info>] [--rules <list>] [--ignore <list>] [--format <text|json|yaml>] [--project <path>]
```

## Problem

SpecScore's structural rules are useless if they are not enforced. Humans miss adherence footers, ship OQ sections that say nothing, mis-nest headings, or let indexes drift from reality. A deterministic linter catches every violation mechanically so review effort goes to content, not conventions.

## Behavior

### Rule suite

The default run enables every known rule. The canonical rule list lives in the [lint package](../../../../../pkg/lint) and is not duplicated here.

#### REQ: default-runs-all-rules

Without `--rules` or `--ignore`, `spec lint` MUST execute every registered rule against the project.

#### REQ: rules-whitelist

`--rules a,b,c` MUST restrict execution to the named rules. Any rule name not in the registered set MUST exit `2` (InvalidArgs) with a message naming the unknown rule. The command MUST fail fast — no lint run is performed when any name is invalid.

#### REQ: ignore-blacklist

`--ignore x,y` MUST exclude the named rules from execution. Unknown rule names in `--ignore` follow the same hard-error policy as `--rules`.

#### REQ: mutually-exclusive-filters

`--rules` and `--ignore` MUST NOT be combined. Supplying both MUST exit `2`.

### Severity filtering

Each rule has a built-in severity (`error`, `warning`, `info`). `--severity` sets the minimum severity reported.

#### REQ: severity-default

Default `--severity` MUST be `error`. Warnings and info-level findings are suppressed unless the caller explicitly widens the filter.

#### REQ: severity-values

`--severity` MUST accept `error`, `warning`, or `info`. Any other value MUST exit `2`. The ordering is strict: `info ⊂ warning ⊂ error` (wider filter includes narrower).

### Autofix

Rules that support autofix declare so in their registration. `--fix` applies only those fixes; rules without autofix still report violations unchanged.

#### REQ: fix-is-safe-subset

`--fix` MUST only mutate what the rule declares safe to mutate. Rules that require human judgment (e.g., wrong-URL adherence footers — possibly mis-classified documents) MUST NOT be autofixed, even with `--fix`.

#### REQ: fix-is-idempotent

Running `spec lint --fix` twice in a row MUST yield no changes on the second run. The second run exits `0` with no output beyond the standard violation report.

### Exit codes

`spec lint` signals violations through the exit code so CI can gate on it.

#### REQ: exit-1-on-violations

If any violation at or above the effective `--severity` is found, the command MUST exit `1`. Exit `0` is only returned when zero such violations exist.

### Output formats

Default output is text. `--format json` and `--format yaml` produce structured output suitable for tool consumption.

## Parameters

None. All inputs are flags.

## Exit codes

| Code | Condition |
|---|---|
| `0` | No violations at or above `--severity` |
| `1` | One or more violations at or above `--severity` |
| `2` | Invalid `--rules` / `--ignore` name, invalid `--severity`/`--format` value, conflicting filters |
| `10` | Unexpected I/O failure while scanning |

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [adherence-footer](../../../adherence-footer/README.md) | Defines the footer contract; `spec lint` enforces it. Autofix inserts missing footers. |
| [idea](../../../idea/README.md) | Declares idea-specific rules; `spec lint` runs them alongside shared rules. |
| [feature](../../../feature/README.md) | Declares the required-sections rules for feature READMEs. |
| [CLI](../../README.md) | Inherits exit-code contract and project autodetection. `spec lint`'s exit-1-on-violations convention is part of the shared contract. |

## Acceptance Criteria

### AC: clean-tree-exits-0

**Requirements:** cli/spec/lint#req:exit-1-on-violations

On a specification tree with zero error-severity violations, `specscore spec lint` exits `0` with no violation lines printed.

### AC: violations-exit-1

**Requirements:** cli/spec/lint#req:exit-1-on-violations

A tree containing at least one error-severity violation exits `1`. The violation list is printed to stdout (text format) or is the payload of the structured output.

### AC: unknown-rule-name-exits-2

**Requirements:** cli/spec/lint#req:rules-whitelist

`specscore spec lint --rules not-a-rule` exits `2` with a message naming the unknown rule. No lint run is performed.

### AC: fix-idempotent

**Requirements:** cli/spec/lint#req:fix-is-idempotent

Running `spec lint --fix` twice consecutively on the same tree yields no file changes on the second run and exits `0`.

## Outstanding Questions

- Should `spec lint` accept a path argument (`spec lint spec/features/cli/`) to lint a subtree, for faster feedback during development? Today the full tree is always scanned.
- Should `--fix` have a paired `--dry-run` that prints the intended edits without applying them, so authors can preview fixes before accepting?

---
*This document follows the https://specscore.md/feature-specification*
