# Feature: Plans Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fplans-index) — graph, discussions, approvals

**Status:** Draft

## Summary

The plans index — `spec/plans/README.md` — is the canonical entry point for understanding what implementation work is planned, in flight, approved, or superseded in a spec repository. It inherits the shared shape of every SpecScore index from the [Index feature](../index/README.md) and declares only the plans-specific overrides: the Contents column set, a Recently Closed section, sub-plan indentation, and an optional per-plan summary convention.

## Problem

Plans accumulate silently without a structured aggregation: readers cannot answer "what's the state of our compute migration?" without opening every plan's README. Tooling that aggregates plan status across repos has no stable shape to parse. The Index feature solves the generic shape; this feature declares the plans-specific additions — the Contents column set, the Recently Closed surface, and plan-nesting semantics — that make the index useful for Plan-kind work.

## Behavior

This feature inherits all shared rules from the [Index feature](../index/README.md): canonical location at `spec/plans/README.md`, the minimum-required-sections rule, completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. Those rules are NOT re-stated here per [Index#req:overrides-only](../index/README.md#req-overrides-only). What follows is the plans-index's own overrides.

### Section structure

The plans index adds one per-domain section, **Recently Closed**, between the list section and Outstanding Questions:

1. `## Contents` — the list-holding section (satisfies [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum)).
2. `## Recently Closed` — plans whose status indicates completion, supersession, or abandonment.
3. `## Outstanding Questions` — the shared terminal section.

Repositories MAY prepend optional header prose and an optional `## Conventions` section for project-specific naming or structural conventions beyond the SpecScore baseline.

#### REQ: required-sections

The plans index MUST contain, in order: `## Contents`, `## Recently Closed`, and `## Outstanding Questions`. Header prose and `## Conventions` are OPTIONAL. This REQ is a per-domain refinement of [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum) — it adds `## Recently Closed` between the list section and the shared terminal section.

### Contents table

The Contents section holds a table with one row per plan. Required columns:

| Column | Meaning |
|---|---|
| Plan | Link to the plan directory — e.g. a plan named `my-plan` links to `my-plan/README.md` |
| Status | `draft` \| `in_review` \| `approved` — or project-extended statuses like `paused`, `superseded`, `completed` |
| Features | Comma-separated list of features the plan touches |
| Effort | `S` \| `M` \| `L` \| `XL` — or `-` if absent |
| Impact | `low` \| `medium` \| `high` \| `critical` — or `-` if absent |
| Author | Handle or name of the plan author |
| Approved | Approval date (`YYYY-MM-DD`) — or `-` if not approved |

#### REQ: contents-columns

The Contents table MUST include columns for Plan, Status, Features, Effort, Impact, Author, and Approved. When ROI metadata (Effort, Impact) or approval metadata is absent, the cell MUST display `-`.

#### REQ: sub-plan-indentation

Sub-plans (plans nested inside a parent plan directory) MUST be indented in the Contents table with `&ensp;` and their link paths MUST include the parent directory (e.g., `chat-feature/chat-infrastructure/`).

### Per-plan summaries

For plans with nontrivial scope, a short paragraph summary MAY appear below the Contents table, introduced by a `### {plan-slug}` subheading. Summaries help readers decide whether to open the plan without mandating that every plan provide one.

#### REQ: summary-optional

Per-plan summary subsections in the index are OPTIONAL. When present, they MUST use `### {plan-slug}` as the subheading.

### Recently Closed section

The Recently Closed section surfaces the most recently closed plans (status indicating `completed`, `superseded`, or project-equivalent). It is short — the N most recent, where N is project-configurable with a default of 5.

#### REQ: recently-closed-present

The plans index MUST include a `## Recently Closed` section. The section MAY be empty (explicitly noting "None at this time"), but MUST be present.

#### REQ: recently-closed-shape

Entries in Recently Closed MUST include the plan link, its terminal status, and the closure date. Older entries beyond the project's configured N are dropped from the list; git history preserves them.

### Adherence footer

#### REQ: adherence-footer

Every plans-index document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/plans-index-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Plan](../plan/README.md) | The plans index aggregates plan metadata. Individual plan documents are specified by the Plan feature; this feature specifies only the index. Plans-Index is a top-level sibling of Plan (not nested) so its specification URL stays flat. |
| [Index](../index/README.md) | Plans-Index inherits its shared shape from the Index Meta feature: location, minimum required sections, completeness, and adherence-footer delegation. This feature declares only the plans-specific overrides (Recently Closed, Contents column set, sub-plan indentation, optional per-plan summaries). |
| [Adherence Footer](../adherence-footer/README.md) | Plans-Index delegates its footer mechanism to the shared Adherence Footer feature; the local `REQ: adherence-footer` above uses the two-line delegation form. |
| [Document Types Registry](../document-types-registry/README.md) | Plans-Index is an Index-Kind entry in the registry with URL `plans-index-specification` and Consumer Path `spec/plans/README.md`. |

## Acceptance Criteria

### AC: plan-specific-sections

**Requirements:** plans-index#req:required-sections

The plans index adds `## Recently Closed` between the shared list section and the shared Outstanding Questions section. Section order is enforced.

### AC: contents-table

**Requirements:** plans-index#req:contents-columns, plans-index#req:sub-plan-indentation

The Contents table includes all plan-specific columns with `-` for absent metadata, and indents sub-plans with `&ensp;`.

### AC: recently-closed

**Requirements:** plans-index#req:recently-closed-present, plans-index#req:recently-closed-shape

A Recently Closed section exists (even if empty) and its entries include plan link, terminal status, and closure date.

### AC: adherence-footer

**Requirements:** plans-index#req:adherence-footer

Every plans index ends with a footer containing the bare URL `https://specscore.md/plans-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the default N for Recently Closed (currently 5) be defined here or be purely project-configurable with no baseline? Current position: project-configurable, default 5.
- Should the index support deeper sub-plan indentation (`&ensp;&ensp;` and beyond) for three-or-more levels of nesting? Current position: one level of indentation sufficient for current SpecScore patterns.
- Should status vocabulary in the Status column be constrained to the SpecScore canonical set (`draft | in_review | approved`), or explicitly permit project-extended values (`paused`, `superseded`, `completed`)? Current position: permit project-extended; tooling treats anything beyond the canonical set as project-specific.

---
*This document follows the https://specscore.md/feature-specification*
