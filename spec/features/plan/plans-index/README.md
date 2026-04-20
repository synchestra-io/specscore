# Feature: Plans Index

**Status:** Draft

## Summary

The plans index — `spec/plans/README.md` — is the canonical entry point for understanding what implementation work is planned, in flight, approved, or superseded in a spec repository. It lists every plan with enough metadata for a reader to scan for relevance and click through for details, and provides a stable structure that tooling can parse. This sub-feature specifies the index's location, structure, required columns, and adherence footer.

## Problem

Without a structured index, plans accumulate silently — each in its own directory, indistinguishable from outside. A reader cannot answer "what's the state of our compute migration?" without opening every plan's README. Tooling that aggregates plan status across repos has no stable shape to parse. Human-written indexes drift in column choice, status vocabulary, and section order, making cross-repo comparison painful.

## Behavior

### Location

Every spec repository that contains plans MUST have a plans index at a single canonical path.

#### REQ: location

The plans index MUST live at `spec/plans/README.md` in each spec repository that contains plans.

### Section structure

The index contains three required sections in this order:

1. **Contents** — the canonical list of plans with standard columns.
2. **Recently Closed** — plans whose status indicates completion, supersession, or abandonment.
3. **Outstanding Questions** — the standard question-lifecycle section shared by all SpecScore documents.

Repositories MAY prepend optional header prose (a short statement of what this plans directory covers) and an optional Conventions section for project-specific naming or structural conventions beyond the SpecScore baseline.

#### REQ: required-sections

The plans index MUST contain, in order: a `## Contents` section, a `## Recently Closed` section, and an `## Outstanding Questions` section. Header prose and `## Conventions` sections are OPTIONAL.

### Contents table

The Contents section holds a table with one row per plan. Required columns:

| Column | Meaning |
|---|---|
| Plan | Link to the plan directory (e.g., `[my-plan](my-plan/README.md)`) |
| Status | `draft` \| `in_review` \| `approved` — or project-extended statuses like `paused`, `superseded`, `completed` |
| Features | Comma-separated list of features the plan touches |
| Effort | `S` \| `M` \| `L` \| `XL` — or `-` if absent |
| Impact | `low` \| `medium` \| `high` \| `critical` — or `-` if absent |
| Author | Handle or name of the plan author |
| Approved | Approval date (`YYYY-MM-DD`) — or `-` if not approved |

#### REQ: contents-columns

The Contents table MUST include columns for Plan, Status, Features, Effort, Impact, Author, and Approved. When ROI metadata (Effort, Impact) or approval metadata is absent, the cell MUST display `-`.

#### REQ: completeness

Every plan directory under `spec/plans/` MUST have a corresponding row in the Contents table. An unlisted plan is a validation error.

#### REQ: sub-plan-indentation

Sub-plans (plans nested inside a parent plan directory) MUST be indented in the Contents table with `&ensp;` and their link paths MUST include the parent directory (e.g., `chat-feature/chat-infrastructure/`).

### Per-plan summaries

For plans with nontrivial scope, a short paragraph summary MAY appear below the Contents table, introduced by a `### {plan-slug}` subheading. Summaries help readers decide whether to open the plan without mandating that every plan provide one.

#### REQ: summary-optional

Per-plan summary subsections in the index are OPTIONAL. When present, they MUST use `### {plan-slug}` as the subheading.

### Recently Closed section

The Recently Closed section surfaces the most recently closed plans (status indicating `completed`, `superseded`, or project-equivalent). It is short — the N most recent, where N is project-configurable with a default of 5.

Each row identifies the plan, its closure status, and the closure date.

#### REQ: recently-closed-present

The plans index MUST include a `## Recently Closed` section. The section MAY be empty (explicitly noting "None at this time"), but MUST be present.

#### REQ: recently-closed-shape

Entries in Recently Closed MUST include the plan link, its terminal status, and the closure date. Older entries beyond the project's configured N are dropped from the list; git history preserves them.

### Adherence footer

Every plans index ends with a footer referencing this SpecScore plans-index specification.

#### REQ: adherence-footer

Every plans index MUST end with a footer that references this SpecScore plans-index specification. The footer is a single italic line preceded by a horizontal rule (`---`), containing the bare URL `https://specscore.md/plans-index-specification` (trailing slash optional). The URL MUST appear in bare form (not wrapped in Markdown link syntax) so it remains clickable in terminals and renderers that detect URLs but do not parse Markdown.

The recommended form is:

```markdown
---
*This index follows the https://specscore.md/plans-index-specification*
```

Authors MAY reword the surrounding prose, but the specification URL MUST be present. Tooling matches on the URL, not the sentence.

The specification URL is unversioned by design. SpecScore commits to additive-only evolution of the plans-index format; breaking changes are not planned. If a breaking change ever becomes necessary, the canonical URL will be reissued rather than branched into parallel versions.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Plan](../README.md) | The plans index aggregates plan metadata. Individual plan documents are specified by the parent plan feature; this sub-feature specifies only the index. |

## Acceptance Criteria

### AC: index-location-and-structure

**Requirements:** plan/plans-index#req:location, plan/plans-index#req:required-sections

The plans index lives at `spec/plans/README.md` and contains the required sections (`Contents`, `Recently Closed`, `Outstanding Questions`) in order.

### AC: contents-table

**Requirements:** plan/plans-index#req:contents-columns, plan/plans-index#req:completeness, plan/plans-index#req:sub-plan-indentation

The Contents table includes all standard columns with `-` for absent metadata, lists every plan under `spec/plans/`, and indents sub-plans with `&ensp;`.

### AC: recently-closed

**Requirements:** plan/plans-index#req:recently-closed-present, plan/plans-index#req:recently-closed-shape

A Recently Closed section exists (even if empty) and its entries include plan link, terminal status, and closure date.

### AC: adherence-footer

**Requirements:** plan/plans-index#req:adherence-footer

Every plans index ends with a footer containing the bare URL `https://specscore.md/plans-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the default N for Recently Closed (currently 5) be defined here or be purely project-configurable with no baseline? Current position: project-configurable, default 5.
- Should the index support deeper sub-plan indentation (`&ensp;&ensp;` and beyond) for three-or-more levels of nesting? Current position: one level of indentation sufficient for current SpecScore patterns.
- Should status vocabulary in the Status column be constrained to the SpecScore canonical set (`draft | in_review | approved`), or explicitly permit project-extended values (`paused`, `superseded`, `completed`)? Current position: permit project-extended; tooling treats anything beyond the canonical set as project-specific.

---
*This document follows the https://specscore.md/feature-specification*
