# Idea: Shared Index Feature

**Status:** Specified
**Date:** 2026-04-20
**Owner:** alexander.trakhimenok
**Promotes To:** index
**Supersedes:** —
**Related Ideas:** depends_on:adherence-footer-and-doc-type-registry

## Problem Statement

How might we define the shape of an "index README" once — so every Index-Kind feature (plans-index today, features-index and ideas-index plausibly) declares only what differs, rather than re-specifying required sections, completeness rules, and adherence-footer boilerplate in each index feature?

## Context

The `plan/plans-index/` sub-feature was extracted in commit 0107396. Its README duplicates ~80% of the structure found in `spec/features/feature/README.md` and other Feature READMEs: Status, required sections, adherence-footer rule, Outstanding Questions section, completeness requirement. If every Document-Kind feature grows its own Index-Kind sibling (features-index, ideas-index, proposals-index, tasks-index), the duplication multiplies — every URL-convention or lifecycle change requires editing N index features.

The upstream Idea `adherence-footer-and-doc-type-registry` flattens `plans-index` to a top-level sibling of `plan` and introduces an `Index` Kind in the registry. That work leaves Index-Kind features as standalone top-level artifacts, each currently bearing the full shape. This Idea is what comes next: factor the shared shape out, one cycle after the registry ships, when there is real data on which parts are truly reused.

## Recommended Direction

Introduce `index/` as a top-level Meta-Kind feature that defines the shared shape of every Index-Kind feature: required sections (`Contents`, `Recently Closed`, `Outstanding Questions`), a tiny column base set for the Contents table (link to item, Status), completeness rule (every instance under the domain's `Consumer Path` must have a row), adherence-footer mechanism (delegated to `adherence-footer`), sub-item indentation convention.

Every Index-Kind feature is reshaped to an overrides-only page. It declares:

- Which Document-Kind it indexes (via a typed link and Dependencies).
- What additional columns it requires beyond the shared base (for plans: Effort, Impact, Author, Approved).
- What Status vocabulary applies and any project-extended values.
- Any domain-specific rules (e.g., plans use `&ensp;` for sub-plan indentation; other domains may not).

If after reshaping, an Index-Kind feature has nothing beyond the delegation, the feature can be deleted entirely — the shared `index/` feature specifies it directly via the registry. This is the signal that the abstraction flattened correctly. For `plans-index` specifically, the Effort/Impact/Approved columns and sub-plan indentation likely remain, so the override page shrinks from ~145 lines to ~40 but does not disappear.

This Idea is deliberately timed to come second. We ship the registry work first, watch `plans-index` live at the top level for a cycle, then extract shared shape based on observed duplication rather than speculative design. Extracting a shared feature from N=1 is how over-abstraction starts.

## Alternatives Considered

**Ship `index/` alongside the registry work** (bundled into the first Idea). Rejected because there is only one Index-Kind feature today (`plans-index`). Extracting from N=1 risks codifying idiosyncrasies into the shared feature. Waiting until a second index is authored gives real data on what is common and what is domain-specific.

**Skip the shared feature; accept per-index duplication.** Rejected because the duplication is substantial (~80%) and grows linearly with the number of indexes. Every URL-convention or lifecycle change requires editing every Index-Kind feature. The DRY benefit is real once the second index appears.

**Make every Index-Kind feature a sub-feature of its Document-Kind parent** (reverting to the original `plan/plans-index` shape). Rejected in the registry Idea: top-level siblings trivialize the URL rule. Not re-litigating here.

## MVP Scope

One cycle after `adherence-footer-and-doc-type-registry` has shipped and at least one cycle of real use has accumulated. Scope: write the `index/` Meta-Kind feature defining shared shape; reshape `plans-index/` to overrides-only; update the registry to reflect the new shape; extend registry lint to cross-check that every Index-Kind feature delegates to `index/`.

## Not Doing (and Why)

- **Creating features-index, ideas-index, proposals-index, tasks-index speculatively** — these get authored when a real need appears. Creating all four up-front reintroduces the very duplication problem this Idea solves.
- **Changing `plans-index`'s normative content** (columns, statuses, indentation rules) — this is a pure refactor. Semantics are preserved; only the location of the rules shifts.
- **Re-specifying the adherence footer format** — inherited from `adherence-footer` feature. This Idea delegates, does not duplicate.
- **Publishing pipeline work** — separate ideation, already parked upstream.
- **Shipping before the registry Idea has stabilized** — the `Kind` taxonomy and the top-level rule must exist first. Shipping out of order risks rework.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | After `plans-index` is reshaped to overrides-only, its override page has non-trivial content — otherwise the feature should be deleted rather than kept. | Sketch the override page before committing to the refactor. If it is <15 lines and adds nothing `index/` does not already cover, delete `plans-index` and let `index/` specify plans directly. |
| Must-be-true | The shared Index shape is extractable — the duplication between `plans-index` today and any second index is genuinely structural, not coincidental. | Before promoting this Idea to Features, author a second Index-Kind feature (features-index is the natural candidate) and confirm ≥70% shape overlap with `plans-index`. If overlap is low, this Idea is wrong and per-index definition is correct. |
| Should-be-true | The delegation mechanism (Index-Kind feature "overrides" `index/`) is expressible in the SpecScore REQ model without inventing new syntax. | Prototype the `plans-index` override page using only existing SpecScore conventions (REQs, Interaction tables, Dependencies, Source Ideas). |
| Might-be-true | Consumers will ever need to define their own custom Index-Kind features in their own repos. | Monitor; out of scope for v1 either way. |

## SpecScore Integration

- **New Features this would create:** `index` (Meta Kind).
- **Existing Features affected:** `plans-index` (reshaped to overrides-only; may be deleted if nothing non-trivial remains), `document-types-registry` (lint extended to cross-check Index-Kind delegation to `index/`).
- **Dependencies:** Depends on `adherence-footer-and-doc-type-registry`. Specifically, the `Kind` taxonomy and the top-level rule for Index-Kind features must exist before this Idea is meaningful.

## Open Questions

- Does `index/` specify the columns of the Contents table as a tiny base set (first column = link to item, plus Status) that per-domain features extend, or is the column set entirely per-domain? Current lean: tiny base with everything else per-domain.
- Does the shared `index/` feature hard-require a `Recently Closed` section, or leave it as per-domain opt-in? Plans have a lifecycle where closed items matter; Ideas have `Archived`; Features have `Deprecated`. Different domains, different recency semantics.
- If `plans-index` turns out to be deletable after reshape, what is the cleanup gesture — delete the directory, or keep a stub with a redirect so historical links keep working?
- Should `index/` define a machine-readable format for the delegation (e.g., a required Interaction table row pointing at `index/`), or rely on convention alone?

---
*This document follows the https://specscore.md/idea-specification*
