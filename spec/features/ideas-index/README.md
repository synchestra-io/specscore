# Feature: Ideas Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fideas-index) — graph, discussions, approvals

**Status:** Draft

## Summary

The ideas index — `spec/ideas/README.md` — is the canonical entry point for understanding the pre-spec thinking in a SpecScore repository. It inherits the shared shape of every SpecScore index from the [Index feature](../index/README.md) and declares only the ideas-specific overrides: the Index column set and the archived-exclusion rule that keeps Archived Ideas out of the active index.

Archived Ideas live at `spec/ideas/archived/README.md` in a different shape (chronological list, not a table) per the [Idea feature's](../idea/README.md) existing `REQ: archived-index-chronological`. This feature covers only the active index; the archived index remains specified within the Idea feature.

## Problem

Without a structured active index, Ideas accumulate silently — each as a standalone file, invisible to anyone who did not author it. The [Index feature](../index/README.md) solves the generic shape (location, completeness, adherence-footer, required-sections minimum); this feature declares the ideas-specific additions — the column set and the Archived-exclusion — that make the index meaningful for the Idea lifecycle.

## Behavior

This feature inherits all shared rules from the [Index feature](../index/README.md): canonical location at `spec/ideas/README.md`, the minimum-required-sections rule (list section + Outstanding Questions), completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. Those rules are NOT re-stated here per [Index#req:overrides-only](../index/README.md#req-overrides-only). What follows is the ideas-index's own overrides.

### Section structure

The ideas index adds no per-domain sections beyond the shared minimum. It contains exactly:

1. `## Index` — the list-holding section (satisfies [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum)).
2. `## Outstanding Questions` — the shared terminal section.

#### REQ: list-section-heading

The list-holding section MUST use the heading `## Index`. This REQ pins the per-domain heading name; the section's required position is inherited from [Index#req:required-sections-minimum](../index/README.md#req-required-sections-minimum).

### Index table

The Index section holds a table with one row per active Idea. Required columns:

| Column | Meaning |
|---|---|
| Idea | Link to the Idea file — e.g. an Idea named `my-idea` links to `my-idea.md` |
| Status | One of `Draft` \| `Under Review` \| `Approved` \| `Specified` — the Idea's current `**Status:**` value. `Archived` Ideas are NEVER listed here. |
| Date | The Idea's `**Date:**` field (`YYYY-MM-DD`) |
| Owner | The Idea's `**Owner:**` field |
| Promotes To | Comma-separated list of Feature slugs from the Idea's `**Promotes To:**` field — or `—` when empty. Managed by tooling via the Idea feature's sync mechanism. |

#### REQ: index-columns

The Index table MUST include columns for Idea, Status, Date, Owner, and Promotes To, in that order. Absent Promotes To entries MUST display `—`.

#### REQ: status-excludes-archived

The Index table MUST NOT list any Idea whose `**Status:**` is `Archived`. Archived Ideas MUST appear only in `spec/ideas/archived/README.md`. This REQ refines the shared [Index#req:completeness](../index/README.md#req-completeness): the universe for completeness is active Ideas only. Moving an Idea to Archived triggers its removal from this index.

### Adherence footer

#### REQ: adherence-footer

Every ideas-index document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/ideas-index-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Idea](../idea/README.md) | The ideas index aggregates active Idea metadata. Individual Idea documents are specified by the Idea feature; this feature specifies only the active index. Ideas-Index is a top-level sibling of Idea (not nested) so its specification URL stays flat. The Idea feature's `REQ: archived-index-chronological` separately specifies the archived index at `spec/ideas/archived/README.md`. |
| [Index](../index/README.md) | Ideas-Index inherits its shared shape from the Index Meta feature: location, minimum required sections, completeness, and adherence-footer delegation. This feature declares only the ideas-specific overrides (Index column set, Archived-exclusion). |
| [Adherence Footer](../adherence-footer/README.md) | Ideas-Index delegates its footer mechanism to the shared Adherence Footer feature; the local `REQ: adherence-footer` above uses the two-line delegation form. |
| [Document Types Registry](../document-types-registry/README.md) | Ideas-Index is an Index-Kind entry in the registry with URL `ideas-index-specification` and Consumer Path `spec/ideas/README.md`. |

## Acceptance Criteria

### AC: list-section-heading

**Requirements:** ideas-index#req:list-section-heading

The list-holding section uses the heading `## Index`. Other headings (e.g., `## Contents`) are rejected.

### AC: index-table

**Requirements:** ideas-index#req:index-columns, ideas-index#req:status-excludes-archived

The Index table includes all required columns (Idea, Status, Date, Owner, Promotes To) and excludes Archived Ideas entirely.

### AC: adherence-footer

**Requirements:** ideas-index#req:adherence-footer

Every ideas index ends with a footer containing the bare URL `https://specscore.md/ideas-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the Index table support an optional Source column showing which downstream work (plan, feature, proposal) originated from the Idea, or is `Promotes To` sufficient for all current needs?
- Should Ideas that have been promoted (Status: Specified) stay in the active index indefinitely, or should they move to a "historical" partition after their promoted Features ship? Current position: remain in active index until Archived.
- Archive policy atomicity: when an Idea is archived, should removing the active-index row and adding the archived-index entry be enforced as atomic (one commit), convention-driven, or lint-enforced?

---
*This document follows the https://specscore.md/feature-specification*
