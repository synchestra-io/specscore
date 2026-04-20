# Feature: Scenarios Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fscenarios-index) — graph, discussions, approvals

**Status:** Draft

## Summary

The scenarios index is the per-feature aggregation of test scenarios living under `{feature}/_tests/`. One scenarios index exists per feature that has scenarios, at the path `{feature-path}/_tests/README.md`. It inherits the shared shape of every SpecScore index from the [Index feature](../index/README.md) and declares only scenario-specific overrides: a compact Scenario / Validates column pair that traces each scenario back to the requirements and acceptance criteria it exercises.

Unlike the other Index-Kind features (plans-index, ideas-index, features-index), scenarios indexes are **per-parent-feature** rather than repo-wide. Each feature that carries scenarios owns its own scenarios index inside its `_tests/` directory.

## Problem

A feature with ten scenarios is hard to navigate without an index: readers opening `_tests/` see a directory of markdown files with no hint of which scenario validates what. Reviewers checking coverage cannot answer "does every REQ have a scenario?" without opening each file. Across the SpecScore meta-spec, nine features already maintain an informal `_tests/README.md` with the same shape — a Scenario / Validates table. This feature pins that convention down, adds the adherence footer every other Index-Kind carries, and makes the index a lint-checkable artifact.

## Behavior

This feature inherits all shared rules from the [Index feature](../index/README.md): canonical location (adapted — see below), minimum required sections (list section + Outstanding Questions), completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. Those rules are NOT re-stated here per [Index#req:overrides-only](../index/README.md#req-overrides-only). What follows is the scenarios-index's own overrides.

### Location — per-feature, not repo-wide

Unlike plans-index, ideas-index, and features-index (each of which is a single repo-wide file), a scenarios index lives **inside each feature's `_tests/` directory**. There are as many scenarios indexes as there are features with scenarios.

#### REQ: per-feature-location

Each feature that has scenarios MUST carry a scenarios index at `{feature-path}/_tests/README.md`, where `{feature-path}` is the directory of the parent feature relative to `spec/features/`. Examples: `spec/features/feature/_tests/README.md`, `spec/features/plan/_tests/README.md`. This refines [Index#req:index-location](../index/README.md#req-index-location): the canonical location template is per-feature, not repo-wide.

A feature with no scenarios MAY omit the `_tests/` directory entirely. When `_tests/` exists and contains scenario files, its `README.md` MUST follow this specification.

### Section structure

A scenarios index is minimal — typically a title line, one table, and Outstanding Questions. No additional sections are required.

#### REQ: list-section-heading

The list-holding section MAY use the heading `## Scenarios`, OR the title line itself (`# Scenarios: {Feature}`) MAY directly precede the table without a separate section heading. Both forms are valid because the scope of a scenarios index is narrow enough that a dedicated list-section subheading is often redundant with the title.

### Index table

The scenarios index holds a table with one row per scenario in the parent feature's `_tests/` directory. Required columns:

| Column | Meaning |
|---|---|
| Scenario | Relative link to the scenario file (`[slug](slug.md)`) |
| Validates | Comma-separated list of REQ or AC references the scenario validates, formatted as `[feature-id#req:slug](...)` or `[feature-id#ac:slug](...)` |

Projects MAY append additional columns (e.g., `Priority`, `Last Updated`). The two required columns MUST remain in the leading positions.

#### REQ: index-columns

The Index table MUST include columns for Scenario and Validates, in that order. The Validates column MUST NOT be empty — a scenario with no validation target is a scenario-lint error independently of this index.

### Cross-feature scenarios

A scenario MAY validate REQs or ACs from features other than its parent (cross-feature scenarios per [Scenario#req:cross-feature-validates](../scenario/README.md#req-cross-feature-validates)). The scenarios index lists such scenarios normally; the Validates column captures all referenced features regardless of parent.

### Adherence footer

#### REQ: adherence-footer

Every scenarios-index document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/scenarios-index-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Scenario](../scenario/README.md) | The scenarios index aggregates scenario metadata for one feature. Individual scenario documents are specified by the Scenario feature; this feature specifies only the per-feature index. Scenarios-Index is a top-level sibling of Scenario (not nested) so its specification URL stays flat. |
| [Index](../index/README.md) | Scenarios-Index inherits its shared shape from the Index Meta feature: minimum required sections, completeness, and adherence-footer delegation. This feature declares only the scenarios-specific overrides (per-feature location, minimal column set). |
| [Adherence Footer](../adherence-footer/README.md) | Scenarios-Index delegates its footer mechanism to the shared Adherence Footer feature; the local `REQ: adherence-footer` above uses the two-line delegation form. |
| [Document Types Registry](../document-types-registry/README.md) | Scenarios-Index is an Index-Kind entry in the registry with URL `scenarios-index-specification` and Consumer Path `spec/features/**/_tests/README.md`. |
| [Feature](../feature/README.md) | Each feature that carries scenarios owns a scenarios index in its own `_tests/` directory. |

## Acceptance Criteria

### AC: per-feature-location

**Requirements:** scenarios-index#req:per-feature-location

Each feature with scenarios has a `_tests/README.md` that follows this specification. A feature with no scenarios may omit `_tests/` entirely; an existing `_tests/` directory with scenario files but no conforming README fails validation.

### AC: index-table

**Requirements:** scenarios-index#req:index-columns

The Index table includes columns Scenario and Validates in the required order. Every listed scenario has a non-empty Validates cell.

### AC: adherence-footer

**Requirements:** scenarios-index#req:adherence-footer

Every scenarios index ends with a footer containing the bare URL `https://specscore.md/scenarios-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the scenarios index enforce a completeness-style rule (every `.md` file in `_tests/` MUST have a row), or is listing opt-in? Current position: completeness-enforced (mirrors [Index#req:completeness](../index/README.md#req-completeness)). Scenario files exist to be tracked; an unlisted scenario is a coverage blind spot.
- Should the Validates column support grouping (e.g., one row per REQ/AC with multiple scenarios) instead of one row per scenario? Current position: one row per scenario; grouping is a reviewer concern, not a spec format.
- Should cross-feature scenarios appear in the scenarios indexes of all features they touch, or only the parent feature's index (the current position per [Scenario#req:cross-feature-placement](../scenario/README.md#req-cross-feature-placement))?

---
*This document follows the https://specscore.md/feature-specification*
