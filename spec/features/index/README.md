# Feature: Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Findex) — graph, discussions, approvals

**Status:** Draft
**Source Ideas:** shared-index-feature

## Summary

The Index feature specifies the **shared shape** of every Index-Kind feature in SpecScore — the structural rules that apply equally to the plans index, the ideas index, and any future index (tasks, features, proposals). It defines what every Index-Kind feature has in common: where the canonical index file lives, which sections it must contain, the completeness rule that tooling enforces, and the adherence-footer delegation every index inherits.

This is a Meta-Kind feature. Consumers do not write files of this type. Instead, every Index-Kind feature (`plans-index`, `ideas-index`, …) delegates its shared rules here and carries only its domain-specific overrides — columns, per-domain sections, status vocabulary, and any domain-specific REQs.

## Problem

Before this feature existed, each Index-Kind feature re-specified the same structural rules from scratch. `plans-index` and `ideas-index` independently declared "must live at a single canonical path", "must have these required sections", "must list every instance", and "must carry an adherence footer". The duplication was ~65% of REQ bodies. A change to the shared rules — say, adding a "must have a Last Updated line" — would require editing every Index-Kind feature.

Worse, the duplication invited drift. One index might accept a looser completeness rule; another might reword the adherence-footer delegation. Reviewers and tooling would have no single definition to check against.

This feature solves the duplication once. Each Index-Kind feature becomes a thin overrides page that declares only what differs from the shared shape.

## Behavior

### Shared structural rules

Every Index-Kind feature inherits these rules from this feature. Index-Kind features MUST NOT re-specify them; they delegate via their Interaction table.

#### REQ: index-location

Every Index-Kind feature's consumer instance MUST live at `spec/{domain}/README.md`, where `{domain}` is the directory that holds the Document-Kind instances the index aggregates. Examples: `spec/plans/README.md`, `spec/ideas/README.md`. The Document Types Registry captures this path as the `Consumer Path` column for the Index-Kind row.

#### REQ: required-sections-minimum

Every index MUST include, at minimum, these two sections in order:

1. A **list-holding section** — a section whose heading is determined by the per-domain Index-Kind feature (e.g., `## Contents` for plans-index, `## Index` for ideas-index), containing the table of instances.
2. An **Outstanding Questions** section — the standard question-lifecycle section shared by all SpecScore documents.

Per-domain Index-Kind features MAY add additional required sections between the list section and Outstanding Questions (e.g., plans-index adds `## Recently Closed`). Sections MUST appear in the order declared by the per-domain feature; Outstanding Questions MUST always be last.

#### REQ: completeness

Every Document-Kind instance under the index's consumer-path domain MUST have a corresponding row in the index's table. An unlisted instance is a validation error. The per-domain Index-Kind feature MAY add exclusions (e.g., ideas-index's `REQ: status-excludes-archived` excludes Archived Ideas).

#### REQ: adherence-footer-delegation

Every Index-Kind feature's local `REQ: adherence-footer` MUST take the delegation form specified by the [Adherence Footer feature](../adherence-footer/README.md): two lines referencing that feature and declaring the per-index URL (`https://specscore.md/{slug}-specification`). Index-Kind features MUST NOT re-specify the footer mechanism in prose.

#### REQ: registry-cross-references

Every Index-Kind feature's Interaction table MUST include three rows cross-referencing:

1. The Document-Kind feature it aggregates (e.g., plans-index references Plan; ideas-index references Idea).
2. The [Adherence Footer feature](../adherence-footer/README.md) — because the index delegates its footer mechanism there.
3. The [Document Types Registry](../document-types-registry/README.md) — the registry entry for the Index-Kind row.

### Per-domain overrides

Index-Kind features declare, in their own README:

- **Per-domain sections** beyond list + Outstanding Questions (e.g., `## Recently Closed`, `## Conventions`).
- **Table columns** — entirely per-domain; this feature does not prescribe a column set.
- **Status vocabulary** — which status values the table's Status column accepts.
- **Domain-specific REQs** — e.g., plans-index's sub-plan indentation, ideas-index's archived-exclusion, optional per-item summaries.

#### REQ: overrides-only

A per-domain Index-Kind feature's Behavior section MUST declare only overrides and domain-specific rules — never re-specify the five shared rules above. A per-domain feature that duplicates an inherited REQ is a lint error (the shared REQ is the single source of truth).

### Deleting an Index-Kind feature

If a per-domain Index-Kind feature's overrides collapse to nothing (no extra sections, no unique columns beyond the shared minimum, no domain-specific REQs), the feature SHOULD be deleted and its registry row updated to describe the index directly via this shared feature. This is the signal that the abstraction has fully absorbed the domain.

#### REQ: delete-empty-overrides

An Index-Kind feature whose Behavior section contains only inherited REQ references and no domain-specific content SHOULD be deleted. Lint MAY emit a warning in this case; it MUST NOT auto-delete (structural removal is a human decision).

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Adherence Footer](../adherence-footer/README.md) | Every Index-Kind feature's `REQ: adherence-footer` delegates to Adherence Footer. This feature enforces that delegation form via `REQ: adherence-footer-delegation`. |
| [Document Types Registry](../document-types-registry/README.md) | Index-Kind rows in the registry use this feature as their shared-shape source. The registry's `REQ: kind-cross-check` enforces that Index-Kind features exist at the top level of `spec/features/`. |
| [Plans Index](../plans-index/README.md) | plans-index inherits from this feature. Its overrides: a `Recently Closed` section, a Contents column set (Plan / Status / Features / Effort / Impact / Author / Approved), a `REQ: sub-plan-indentation` for nested plans, and an optional `summary-optional` per-plan paragraph. |
| [Ideas Index](../ideas-index/README.md) | ideas-index inherits from this feature. Its overrides: an Index column set (Idea / Status / Date / Owner / Promotes To) and a `REQ: status-excludes-archived` that keeps Archived Ideas out of the active index. |
| [Features Index](../features-index/README.md) | features-index inherits from this feature. Its overrides: an Index column set (Feature / Status / Description) and a `REQ: top-level-only` that excludes sub-features (they appear in their parent Feature's Contents section). |

## Acceptance Criteria

### AC: shared-rules-enforced

**Requirements:** index#req:index-location, index#req:required-sections-minimum, index#req:completeness, index#req:adherence-footer-delegation, index#req:registry-cross-references

Every Index-Kind feature inherits the five shared rules: canonical location, minimum required sections (list section + Outstanding Questions, with Outstanding Questions last), completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. A feature missing any is rejected.

### AC: no-duplicate-rules

**Requirements:** index#req:overrides-only

A per-domain Index-Kind feature's Behavior section declares overrides only — no duplicated shared REQs. Duplication is a lint error because it creates drift risk.

### AC: empty-overrides-deletable

**Requirements:** index#req:delete-empty-overrides

When an Index-Kind feature's overrides become empty (no per-domain sections, no unique columns, no domain-specific REQs), the feature SHOULD be deleted. Lint may emit a warning; deletion remains a human decision.

## Outstanding Questions

- Should the shared feature prescribe a minimum or default column set (e.g., always at least `{Item link}` + `Status`) or leave column composition entirely per-domain? Current position: entirely per-domain — column semantics vary too much across domains to share any.
- Should Outstanding-Questions-last be enforced by lint, or treated as convention? Current position: lint-enforced, because section order is a determinant of reviewer navigation.
- How should cross-repo Index-Kind features be handled when consumer repos define custom Document Kinds with their own indexes? This requires the Document Types Registry's consumer-extension mechanism (currently out of scope for v1).

---
*This document follows the https://specscore.md/feature-specification*
