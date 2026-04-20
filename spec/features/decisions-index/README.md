# Feature: Decisions Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fdecisions-index) — graph, discussions, approvals

**Status:** Draft
**Source Ideas:** decision-and-decisions-index

## Summary

The decisions index — `spec/decisions/README.md` — is the canonical entry point for reviewing the architectural, product, and process choices made in a SpecScore repository. It inherits the shared shape of every SpecScore index from the [Index feature](../index/README.md) and declares only the decisions-specific overrides: the Index column set, the numeric ordering rule, and the exclusion of Superseded and Deprecated Decisions from the active index.

Superseded and Deprecated Decisions live at `spec/decisions/archived/README.md` in a chronological list (oldest first), per the [Decision feature's](../decision/README.md) archival rule. This feature covers only the active index; the archived index shape is specified within the Decision feature.

## Problem

Without a structured active index, Decisions accumulate silently — one numbered file per choice, invisible to anyone who did not author it. The [Index feature](../index/README.md) solves the generic shape (location, completeness, adherence-footer, required-sections minimum); this feature declares the decisions-specific additions — the column set, the numeric ordering, and the archived-exclusion — that make the index meaningful for the Decision lifecycle.

## Behavior

This feature inherits all shared rules from the [Index feature](../index/README.md): canonical location at `spec/decisions/README.md`, the minimum-required-sections rule (list section + Outstanding Questions), completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. Those rules are NOT re-stated here per [Index#req:overrides-only](../index/README.md#req-overrides-only). What follows is the decisions-index's own overrides.

### Section structure

The decisions index adds no per-domain sections beyond the shared minimum. It contains exactly:

1. `## Decisions` — the list-holding section (satisfies [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum)).
2. `## Outstanding Questions` — the shared terminal section.

#### REQ: list-section-heading

The list-holding section MUST use the heading `## Decisions`. This REQ pins the per-domain heading name; the section's required position is inherited from [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum).

### Index table

The Decisions section holds a table with one row per active Decision. Required columns:

| Column | Meaning |
|---|---|
| # | Zero-padded four-digit Decision number, linked to the Decision file — e.g. `[0007](0007-postgres-over-mongo.md)` |
| Decision | Plain-text title of the Decision (the part after `# Decision:` in the target file) |
| Status | One of `Proposed` \| `Accepted` — the Decision's current `**Status:**`. `Superseded` and `Deprecated` Decisions are NEVER listed here. |
| Date | The Decision's `**Date:**` field (`YYYY-MM-DD`) |
| Tags | Comma-separated list from the Decision's `**Tags:**` field — or `—` when empty |
| Affected | Comma-separated list of Feature slugs from the Decision's `## Affected Features` section — or `—` when empty |

#### REQ: index-columns

The Decisions table MUST include columns for #, Decision, Status, Date, Tags, and Affected, in that order. Empty Tags and Affected cells MUST display `—`.

#### REQ: status-excludes-archived

The Decisions table MUST NOT list any Decision whose `**Status:**` is `Superseded` or `Deprecated`. Archived Decisions MUST appear only in `spec/decisions/archived/README.md`. This REQ refines the shared [Index#req:completeness](../index/README.md#req-completeness): the universe for completeness is active Decisions only. Transitioning a Decision to `Superseded` or `Deprecated` triggers its removal from this index.

#### REQ: numeric-ordering

Rows in the Decisions table MUST appear in ascending order by the `#` column. A table out of numeric order is a validation error (auto-fixable by `specscore lint --fix`). This differs from `ideas-index` (which does not prescribe ordering) because numeric IDs give Decisions a natural, stable sort key that matches how reviewers scan ADRs in the wild — oldest first, most recent at the bottom.

### Archived index (separate file)

The archived decisions index at `spec/decisions/archived/README.md` is shaped by the [Decision feature's](../decision/README.md) archival rules, not this feature. It is a chronological list (oldest first) of Superseded and Deprecated Decisions, with each entry of the form:

```markdown
- YYYY-MM-DD — [NNNN-slug](NNNN-slug.md) — <Status> — <reason or successor reference>
```

The `reason or successor reference` is the successor Decision ID (`→ D-0042`) for Superseded entries, or the first paragraph of the Decision's `## Context` for Deprecated entries.

#### REQ: archived-index-chronological

Entries in `spec/decisions/archived/README.md` MUST appear in chronological order by each Decision's `**Date:**` field (oldest first). Ordering violations are a validation error (auto-fixable by `specscore lint --fix`). The decisions-index feature enforces this REQ here rather than in the Decision feature because it is an index-shaping rule, not a Decision-authoring rule.

### Adherence footer

#### REQ: adherence-footer

Every decisions-index document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/decisions-index-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Decision](../decision/README.md) | The decisions index aggregates active Decision metadata. Individual Decision documents are specified by the Decision feature; this feature specifies the active index shape and (via `REQ: archived-index-chronological`) the archived-index ordering rule. |
| [Index](../index/README.md) | Decisions-Index inherits its shared shape from the Index Meta feature: location, minimum required sections, completeness, and adherence-footer delegation. This feature declares only the decisions-specific overrides (column set, numeric ordering, archived-exclusion). |
| [Adherence Footer](../adherence-footer/README.md) | Decisions-Index delegates its footer mechanism to the shared Adherence Footer feature; the local `REQ: adherence-footer` above uses the two-line delegation form. |
| [Document Types Registry](../document-types-registry/README.md) | Decisions-Index is an Index-Kind entry in the registry with URL `decisions-index-specification` and Consumer Path `spec/decisions/README.md`. |

## Acceptance Criteria

### AC: list-section-heading

**Requirements:** decisions-index#req:list-section-heading

The list-holding section uses the heading `## Decisions`. Other headings are rejected.

### AC: index-table

**Requirements:** decisions-index#req:index-columns, decisions-index#req:status-excludes-archived, decisions-index#req:numeric-ordering

The Decisions table includes all required columns (#, Decision, Status, Date, Tags, Affected) in order, excludes Superseded and Deprecated Decisions entirely, and appears in ascending numeric order by `#`.

### AC: archived-index

**Requirements:** decisions-index#req:archived-index-chronological

The archived index at `spec/decisions/archived/README.md` lists every Superseded and Deprecated Decision in chronological order by `**Date:**` (oldest first). Ordering violations are rejected by lint and auto-fixable.

### AC: adherence-footer

**Requirements:** decisions-index#req:adherence-footer

Every decisions index ends with a footer containing the bare URL `https://specscore.md/decisions-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the Index table gain a `Source Idea` column linking Decisions back to originating Ideas, or is the relationship niche enough that cross-referencing from within the Decision file is sufficient? Current position: no column; the link is one-to-optional-one and visible inside the Decision itself.
- Should the Affected column render as links to the Feature READMEs rather than plain slugs, at the cost of a longer, wrapped cell? Current position: plain slugs; brevity wins until a consumer requests linked rendering.
- Should `lint --fix` auto-reorder rows to numeric ascending order, or only flag violations? Current position: auto-fix; numeric reordering is mechanical and safe.

---
*This document follows the https://specscore.md/feature-specification*
