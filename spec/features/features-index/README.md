# Feature: Features Index

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Ffeatures-index) — graph, discussions, approvals

**Status:** Draft

## Summary

The features index — `spec/features/README.md` — is the canonical entry point for understanding the feature surface of a SpecScore repository. It inherits the shared shape of every SpecScore index from the [Index feature](../index/README.md) and declares only the features-specific override: the column set (Feature, Status, Description). Consumer repos that use SpecScore carry a features-index that follows this spec; no extra metadata is required.

In the SpecScore meta-spec repo specifically, the same file doubles as the [Document Types Registry](../document-types-registry/README.md) — an overlay that adds four columns (Kind, URL, Consumer Path, Index) capturing SpecScore's internal doc type taxonomy. The overlay applies only to SpecScore's own meta-spec; consumer repos do not need the registry columns.

## Problem

A SpecScore repo can have dozens of feature specifications, nested and flat, with varied status. Readers, reviewers, and tooling need a single place to answer "what features exist and what's their status?". Without a pinned shape for this listing, each repo drifts on column choice, row ordering, and section layout, making cross-repo comparison painful. The [Index feature](../index/README.md) solves the generic shape; this feature declares the features-specific columns that make the listing useful for the Feature-kind surface.

## Behavior

This feature inherits all shared rules from the [Index feature](../index/README.md): canonical location at `spec/features/README.md`, the minimum-required-sections rule (list section + Outstanding Questions), completeness, adherence-footer delegation form, and the three registry cross-references in the Interaction table. Those rules are NOT re-stated here per [Index#req:overrides-only](../index/README.md#req-overrides-only). What follows is the features-index's own overrides.

### Section structure

The features index adds no per-domain sections beyond the shared minimum. It contains, at minimum:

1. A list-holding section — the Contents table.
2. `## Outstanding Questions` — the shared terminal section.

Repositories MAY prepend optional header prose introducing the feature surface, and MAY append optional sections beyond Outstanding Questions only with project-specific rationale. The canonical example in this repo adds a `## Feature Hierarchy` ASCII diagram and an `## Integration with Orchestration Tools` narrative section between Contents and Outstanding Questions.

#### REQ: list-section-heading

The list-holding section MUST be the first substantive section after optional header prose. Its heading MAY be `## Contents`, `## Features`, or an equivalent project-appropriate label. Heading choice is not constrained because `spec/features/README.md` often presents the feature surface as the whole document's subject and may omit a dedicated section heading for the table.

### Index table

The features index holds a table with one row per top-level feature. Required columns:

| Column | Meaning |
|---|---|
| Feature | Relative link to the feature's README (target ends in `/README.md`) |
| Status | The feature's declared `**Status:**` (`Draft` \| `In Progress` \| `Stable` \| `Deprecated` \| `Conceptual`) |
| Description | One-line human summary of the feature |

Projects MAY append additional columns after Description (e.g., `Owner`, `Last Updated`, or — in the SpecScore meta-spec — the four registry columns). The three required columns MUST appear in the order above.

#### REQ: index-columns

The Index table MUST include columns for Feature, Status, and Description, in that order. Additional columns are permitted; the three required columns MUST remain in the leading positions.

### Nesting and sub-features

The features index lists **top-level features only** — directories directly under `spec/features/`. Sub-features (feature directories nested inside a parent feature) are NOT listed in the root features index; they appear in their parent Feature README's `## Contents` section per [Feature#req:contents-when-children](../feature/README.md#req-contents-when-children).

#### REQ: top-level-only

The Index table MUST list every top-level feature under `spec/features/` and MUST NOT list sub-features. Sub-features are surfaced through their parent Feature's own Contents section. This is a per-domain refinement of [Index#req:completeness](../index/README.md#req-completeness): the completeness universe is top-level features only.

### Adherence footer

#### REQ: adherence-footer

Every features-index document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/features-index-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | The features index aggregates top-level Feature metadata. Individual Feature documents are specified by the Feature feature; this feature specifies only the index. Features-Index is a top-level sibling of Feature (not nested) so its specification URL stays flat. |
| [Index](../index/README.md) | Features-Index inherits its shared shape from the Index Meta feature: location, minimum required sections, completeness, and adherence-footer delegation. This feature declares only the features-specific overrides (column set, top-level-only). |
| [Adherence Footer](../adherence-footer/README.md) | Features-Index delegates its footer mechanism to the shared Adherence Footer feature; the local `REQ: adherence-footer` above uses the two-line delegation form. |
| [Document Types Registry](../document-types-registry/README.md) | In the SpecScore meta-spec specifically, the features-index document is ALSO the Document Types Registry. The registry is a Meta overlay that adds four columns (Kind, URL, Consumer Path, Index) on top of this feature's base shape. Consumer repos carry only the features-index shape; they do not carry the registry. |

## Acceptance Criteria

### AC: index-table

**Requirements:** features-index#req:index-columns, features-index#req:top-level-only

The Index table includes columns Feature, Status, Description in the required order, lists every top-level feature under `spec/features/`, and excludes sub-features.

### AC: adherence-footer

**Requirements:** features-index#req:adherence-footer

Every features index ends with a footer containing the bare URL `https://specscore.md/features-index-specification`. Lint matches on the URL, not the prose.

## Outstanding Questions

- Should the features index sort feature rows by a stable rule (alphabetical, Status then name, creation order), or is sort order a project choice? Current position: project choice.
- When a feature is `Deprecated`, should it continue to appear in the main Index table or move to a dedicated "Deprecated" sub-section? Current position: remain in the main table with `Deprecated` in the Status cell.
- Should the features index support grouping features into categories (e.g., via intermediate `### Category` subheadings between Contents and the table)? Current position: no — flat listing keeps the shape consistent across repos; narrative grouping belongs in the optional header prose.

---
*This document follows the https://specscore.md/feature-specification*
