# Feature: Document Types Registry

> [View in Synchestra Hub](https://hub.synchestra.io/project/features?id=specscore@synchestra-io@github.com&path=spec%2Ffeatures%2Fdocument-types-registry) — graph, discussions, approvals

**Status:** Draft
**Source Ideas:** adherence-footer-and-doc-type-registry

## Summary

The document types registry is the single canonical list of every document type SpecScore defines — Features, Plans, Ideas, Tasks, Scenarios, Indexes, and the Meta features that define SpecScore itself. The registry is a **Meta overlay on the [features-index](../features-index/README.md)**: in the SpecScore meta-spec repository, `spec/features/README.md` serves simultaneously as the features-index (base shape) AND the registry (four additional columns capturing the doc type taxonomy). This feature specifies the overlay — the additional columns, the Kind taxonomy, and the lint cross-check that prevents drift between registry rows and feature-declared URLs.

Consumer repos that use SpecScore carry a features-index following the [features-index specification](../features-index/README.md), but do NOT carry the registry overlay — the registry describes SpecScore's own internal taxonomy, not the consumer's features.

## Problem

Without a registry, the set of SpecScore document types is implicit and scattered. A reviewer who wants to answer "what are all the document types?" must walk `spec/features/` directory by directory, reading each README to infer whether it describes a standalone document, an embedded structure, an index, or SpecScore meta. Lint has no machine-readable target to cross-check adherence-footer URLs against — the URLs live only in prose inside each feature, so drift between a feature's declared URL and any other declaration is invisible. As SpecScore grows and consumer repos proliferate, the taxonomy needs a first-class home.

## Behavior

### Registry location

The registry IS the root features-index, extended. It is not a separate file.

#### REQ: registry-location

In the SpecScore meta-spec repository, the document types registry MUST live at `spec/features/README.md` — co-located with the [features-index](../features-index/README.md). The registry is the features-index's Contents table extended with the four overlay columns defined below. A separate registry file (e.g., `spec/document-types.md`) MUST NOT be created — duplication would reintroduce the drift problem this feature exists to solve.

Consumer repos do NOT carry the registry overlay. Their `spec/features/README.md` follows the features-index specification alone; the registry applies only to SpecScore's own meta-spec.

### Required columns

The registry overlay adds four columns to the features-index base (Feature, Status, Description). These four overlay columns sit between Status and Description.

#### REQ: required-columns

The Contents table in SpecScore's `spec/features/README.md` MUST include, at minimum, the features-index base columns AND the four registry overlay columns, in this combined order:

| Column | Origin | Meaning |
|---|---|---|
| Feature | features-index | Relative link to the feature's README — e.g. the feature named `foo` is linked as a Markdown link whose target is `foo/README.md` |
| Status | features-index | The feature's spec status (`Draft` \| `In Progress` \| `Stable` \| `Deprecated` \| `Conceptual`) |
| Kind | registry overlay | One of: `Document` \| `Index` \| `Structure` \| `Meta` (see [Kind taxonomy](#kind-taxonomy)) |
| URL | registry overlay | The bare specification URL, or `—` for Kinds that have none (see [REQ: url-per-kind](#req-url-per-kind)) |
| Consumer Path | registry overlay | Glob of where instances live in consumer repos (e.g. `spec/plans/**/README.md`), or `—` when no consumer instances exist |
| Index | registry overlay | Relative link to the companion Index-Kind feature, or `—` when none |
| Description | features-index | One-line human summary |

Additional columns MAY be appended after `Description` for project-level needs. The required columns MUST appear in the order above.

### Kind taxonomy

Every feature is classified into one of four Kinds. The Kind determines how lint treats the feature and its instances.

#### REQ: kind-values

The `Kind` column MUST be one of:

- **`Document`** — Consumers write standalone files of this type. Has a non-empty `URL` and a non-empty `Consumer Path`. Examples: `feature`, `plan`, `idea`, `task`, `scenario`, `project-definition`.
- **`Index`** — A canonical aggregation README for a Document Kind. Lives at a well-known consumer path (e.g. `spec/plans/README.md`). Has a non-empty `URL` and a non-empty `Consumer Path`. The Document-Kind feature the Index aggregates MUST cross-reference the Index via the `Index` column. Example: `plans-index`.
- **`Structure`** — Embedded inside another Document; not a standalone file. `Consumer Path` is `—`. `URL` MAY be present (for educational linking) or `—`. Examples: `requirement`, `acceptance-criteria`, `source-references`.
- **`Meta`** — Defines SpecScore itself. Not authored by consumers. `Consumer Path` is `—`. `URL` is `—`. Examples: `adherence-footer`, `document-types-registry`.

#### REQ: url-per-kind

For every feature row:

- `Kind = Document` or `Kind = Index` → `URL` MUST be a bare `https://specscore.md/{type}-specification` URL, matching the URL declared in the feature's local `REQ: adherence-footer` (see [Adherence Footer#req:delegation-form](../adherence-footer/README.md#req-delegation-form)).
- `Kind = Structure` → `URL` MAY be a bare specification URL or `—`. When present, the URL resolves to the feature's own specscore.md rendering (informational; consumers do not write standalone files).
- `Kind = Meta` → `URL` MUST be `—`.

#### REQ: consumer-path-per-kind

For every feature row:

- `Kind = Document` → `Consumer Path` MUST be a glob describing where instances of this type live in consumer repos (e.g. `spec/plans/**/README.md`, `spec/ideas/*.md`).
- `Kind = Index` → `Consumer Path` MUST be the single path where the canonical index lives (e.g. `spec/plans/README.md`).
- `Kind = Structure` → `Consumer Path` MUST be `—`.
- `Kind = Meta` → `Consumer Path` MUST be `—`.

### Top-level placement for Document and Index Kinds

Document-Kind and Index-Kind features are the only Kinds that consumers interact with by writing files. To keep the URL convention trivial (terminal slug = URL stem), these Kinds live at the top of the feature tree.

#### REQ: top-level-document-kinds

Features with `Kind ∈ {Document, Index}` MUST live directly under `spec/features/` — not as sub-features of another feature. A Document-Kind or Index-Kind feature at `spec/features/<parent>/<child>/` is a validation error. Structure-Kind and Meta-Kind features MAY nest freely.

### Registration and lint cross-check

The registry is the source of truth. Every Document/Index-Kind feature cross-validates against its registry row.

#### REQ: every-feature-registered

Every directory under `spec/features/` whose README exists MUST have a corresponding row in the Contents table of `spec/features/README.md`. An unregistered feature is a validation error, as is a registry row pointing to a non-existent feature directory.

#### REQ: url-cross-check

For every row where `URL` is not `—`, the URL string MUST be present in the referenced feature's local `REQ: adherence-footer` body. Drift between the registry `URL` column and the feature's declared URL is a lint error; `specscore lint --fix` MUST NOT reconcile automatically — a mismatch may indicate the document type has been mis-classified, and silent rewriting would mask the bug.

#### REQ: kind-cross-check

A Document-Kind or Index-Kind feature MUST have a `REQ: adherence-footer` in its Behavior section. A Structure-Kind or Meta-Kind feature MUST NOT have a `REQ: adherence-footer` (it has no consumer instances to footer). Lint errors on either violation.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Features Index](../features-index/README.md) | The registry is a Meta overlay on the features-index document at `spec/features/README.md`. Features-index specifies the base shape (Feature, Status, Description columns); this registry adds the four overlay columns (Kind, URL, Consumer Path, Index) that only the SpecScore meta-spec needs. |
| [Feature](../feature/README.md) | Every feature has a row in the combined features-index + registry table. Cross-checks between each feature's declared specification URL and its registry row catch drift. |
| [Adherence Footer](../adherence-footer/README.md) | The registry's `URL` column is cross-checked against each Document/Index-Kind feature's local `REQ: adherence-footer` URL. Drift is a lint error. |
| [Plans Index](../plans-index/README.md) | Plans-Index is an `Index` Kind whose `URL` is `https://specscore.md/plans-index-specification` and whose `Consumer Path` is `spec/plans/README.md`. Its row in the registry demonstrates the Index Kind conventions. |

## Acceptance Criteria

### AC: registry-shape-valid

**Requirements:** document-types-registry#req:registry-location, document-types-registry#req:required-columns

The root feature index contains a Contents table with the required columns (Feature, Status, Kind, URL, Consumer Path, Index, Description) in the required order. A tree missing any column or locating the registry elsewhere fails validation.

### AC: kind-taxonomy-enforced

**Requirements:** document-types-registry#req:kind-values, document-types-registry#req:url-per-kind, document-types-registry#req:consumer-path-per-kind

Every row's Kind is one of the four legal values, and its URL and Consumer Path cells match the rules for its Kind. A Document-Kind row with an empty URL or empty Consumer Path fails validation, as does a Meta-Kind row with a non-empty URL.

### AC: top-level-placement

**Requirements:** document-types-registry#req:top-level-document-kinds

Every Document-Kind and Index-Kind feature lives directly under `spec/features/`. A Document/Index-Kind feature nested under another feature directory fails validation.

### AC: registration-and-cross-check

**Requirements:** document-types-registry#req:every-feature-registered, document-types-registry#req:url-cross-check, document-types-registry#req:kind-cross-check

Every feature directory has a registry row; every registry row points to a real feature directory; every Document/Index-Kind feature's local adherence-footer URL matches its registry URL; Structure/Meta-Kind features do not declare a `REQ: adherence-footer`.

## Outstanding Questions

- Should the registry gain a `Spec Status` column distinct from the feature's own Status, or is the feature's Status column sufficient for both purposes?
- Should `Consumer Path` accept multiple globs (for document types whose instances live in more than one location), or is single-glob-per-type sufficient forever?
- Should consumer repos be allowed to extend the registry with custom document types (via a `spec/document-types.md` overlay or a project-definition field), and how would lint combine the base and overlay?
- When a Document-Kind feature is renamed (e.g., `plan` → `development-plan`), does the URL change (breaking every existing consumer document) or stay stable (breaking the terminal-slug-equals-URL convention)? Current lean: favor URL stability; rename-without-URL-change requires an explicit registry entry decoupling the two.

---
*This document follows the https://specscore.md/feature-specification*
