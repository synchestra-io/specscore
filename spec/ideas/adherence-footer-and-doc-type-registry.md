# Idea: Adherence Footer and Document Types Registry

**Status:** Specified
**Date:** 2026-04-20
**Owner:** alexander.trakhimenok
**Promotes To:** adherence-footer, document-types-registry
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we let every SpecScore document declare — in one machine-verifiable, human-readable line — which format spec it conforms to, so the rule is defined once (not copy-pasted into every feature README), lint enforces it per-document-type, and the spec catalog is a single source of truth?

## Context

Today `REQ: adherence-footer` lives in `spec/features/feature/README.md` and hardcodes the URL `https://specscore.md/feature-specification`. The lint checker `pkg/lint/adherence_footer.go` hardcodes the same constant. One exception exists: the recently-extracted `spec/features/plan/plans-index/` sub-feature defines `plans-index-specification`.

Two distinct layers use adherence footers, and conflating them has caused confusion:

- **Meta-spec layer** — files under `spec/features/*/README.md` are Feature READMEs. They correctly adhere to `feature-specification`. No change needed at this layer. The file `spec/features/plan/README.md` describes the Plan document type, but the file *itself* is a Feature — its footer points at `feature-specification`, not `plan-specification`.
- **Consumer-repo layer** — files like `spec/plans/my-plan/README.md`, `spec/plans/README.md`, `spec/ideas/<slug>.md` should adhere to per-type URLs (`plan-specification`, `plans-index-specification`, `idea-specification`). Today no feature specifies these URLs for consumers; lint can only enforce one URL.

Side effect noticed during triage: `spec/features/feature/README.md` references `../outstanding-questions/README.md`, but that directory does not exist. A phantom feature. A registry surfaces this gap and forces resolution.

## Recommended Direction

Introduce two new sibling features as the mechanism, and elevate the feature index into the canonical registry.

**`adherence-footer/`** (Meta-Kind) defines the footer mechanism once: syntax (italic line after `---`), bare-URL rule, lint behavior, unversioned-URL policy. Every Document-Kind feature's local `REQ: adherence-footer` shrinks to two lines: "Consumers of this document type MUST end their files with an adherence footer per the Adherence Footer feature. The URL MUST be `https://specscore.md/{type}-specification`." No semantic duplication.

**`document-types-registry/`** (Meta-Kind) specifies that `spec/features/README.md` is the canonical registry. Adds required columns: `Kind` (Document / Index / Structure / Meta), `URL`, `Consumer Path` (where instances live), `Index` (back-reference for Document Kinds that have an Index counterpart). Lint cross-checks every Document/Index-Kind feature's declared URL against the registry row.

**Structural rule**: Document-Kind and Index-Kind features MUST live at the top level of `spec/features/`. This makes terminal-slug the URL stem (trivial, no nesting-rule to invent). As a consequence, `plan/plans-index/` is flattened to a top-level sibling `plans-index/` cross-referenced with `plan/` via Dependencies and Interaction tables.

The lint `pkg/lint/adherence_footer.go` is refactored to read `Kind`, `URL`, and `Consumer Path` from the registry, then walk consumer paths and check the correct footer URL per Kind. Structure and Meta Kinds are not walked.

## Alternatives Considered

**Convention-over-configuration** — no registry; derive URL from the feature's path (`spec/features/plan/` → `plan-specification`). Rejected because (a) the `plans-index` precedent flattens nested slugs in a way that is not derivable from path structure, (b) a registry provides first-class discoverability of "what doc types does SpecScore define?" that a pure convention does not, (c) lint cross-checking is harder without an explicit target.

**Single "Adherence Footer" feature with no registry** — mechanism defined once, each Document-Kind feature declares its own URL in its body without a central index. Rejected because without a registry there is no single place to answer "what are the canonical doc types?", no cross-check target for lint, and the doc-type taxonomy stays implicit — we lose the reviewer leverage of seeing the full list in one table.

**Keep the status quo and only fix consumer lint** — extend lint to know per-path expected URLs without introducing a registry feature, burying the mapping inside Go code. Rejected because the mapping table exists somewhere regardless; hiding it in source code is strictly worse than making it a first-class spec artifact. The registry pays for itself the first time a URL changes or a doc type is added.

## MVP Scope

One cycle. Scope is the two new Meta features (`adherence-footer` and `document-types-registry`), the structural flattening of `plans-index` to a top-level sibling, the shrinkage-or-addition of each Document-Kind feature's `REQ: adherence-footer` to the two-line form, and the lint refactor. The shared `index/` Meta feature is out of scope and lives in the dependent Idea `shared-index-feature`. Gate on lint transitioning `warn` → `error` after all doc-type features are updated. No work on the publishing pipeline or dual-audience URLs.

## Not Doing (and Why)

- **Shared `index/` feature** — split into a separate Idea (`shared-index-feature`). Better shipped after this Idea, once there is real data on which index shapes actually overlap.
- **Publishing pipeline to specscore.md** — separate ideation. The adherence rule is valid even while URLs 404.
- **Dual-audience URL serving (MD vs HTML)** — part of the publishing ideation. Parked.
- **Migrating existing meta-spec READMEs away from `feature-specification`** — they are Features. Their footer is already correct. The layer distinction is load-bearing; violating it would silently corrupt the meta-spec.
- **Auto-rewriting wrong footer URLs in `specscore lint --fix`** — masks semantic bugs (a wrong URL may mean the document is mis-typed). `--fix` inserts missing footers only; mismatched URLs remain hard errors.
- **Consumer-repo custom doc type extension** — out of scope for v1. Revisit when a real consumer needs it.
- **YAML frontmatter for feature metadata** — orthogonal open question in `feature/README.md`.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | The feature index (`spec/features/README.md`) can absorb `Kind`, `URL`, `Consumer Path`, `Index` columns without becoming unreadable for humans. | Produce the extended registry table populated for all ~10 current features before committing to the column set. Acceptance: table fits in standard terminal width and a reviewer scanning it can answer "what are all SpecScore doc types?" in under 10 seconds. |
| Must-be-true | The four-category Kind taxonomy (Document, Index, Structure, Meta) classifies every current feature cleanly — no feature is ambiguous or needs a fifth Kind. | Classify every existing feature in `spec/features/*` and confirm no ties. Any feature that resists classification is a signal the taxonomy is wrong. |
| Must-be-true | Document-Kind and Index-Kind features can be constrained to top-level without breaking the `plans-index` precedent or any in-flight work. | Enumerate existing nested Document/Index-Kind features (today: only `plan/plans-index`) and confirm flattening is a strict improvement. |
| Should-be-true | No real consumer needs to register custom SpecScore document types in the next 6 months. | Ask the 2–3 early adopters whether they are thinking in terms of custom doc types. If yes, bump the extension mechanism into scope. |
| Should-be-true | Lint can migrate in a warn-then-error transition without breaking existing CI for consumer repos already using `specscore lint`. | Review consumer repo CI expectations; add a `--strict` flag if gradual migration is needed. |
| Might-be-true | The registry stays small enough (N < 30 doc types) that a hand-maintained table is preferable to generated output. | If registry grows beyond 30 rows in year 1, reconsider generation. |

## SpecScore Integration

- **New Features this would create:** `adherence-footer` (Meta Kind), `document-types-registry` (Meta Kind), top-level `plans-index` (Index Kind — reparented from `plan/plans-index`).
- **Existing Features affected:** `feature` (local REQ: adherence-footer shrinks and references the new feature; index-completeness REQ gains new columns), `plan`, `idea`, `task`, `scenario`, `requirement`, `acceptance-criteria`, `source-references`, `project-definition` (each gains or shrinks a two-line `REQ: adherence-footer` pointing at its consumer URL), root `spec/features/README.md` (new required columns populated), `pkg/lint/adherence_footer.go` (refactored).
- **Dependencies:** None. This Idea is the foundation for `shared-index-feature`.

## Open Questions

- Is `Kind` + `URL` + `Consumer Path` + `Index` the right minimum column set, or do we also need `Spec Status` (mirrored from the feature's own Status) and `Instance File Pattern` (for lint globbing)?
- `outstanding-questions` phantom feature — create it as a Structure Kind, or remove the references from `feature/README.md`? The registry work forces the decision.
- Should the shrunk two-line form of each Document-Kind feature's `REQ: adherence-footer` be enforced by a template rule (lint checks for an exact two-line shape), or remain convention?
- Does the top-level-only rule for Document/Index Kinds apply only to the SpecScore meta-spec, or also to consumer repos that define custom doc-type-like features? Current position: meta-spec only; consumer repos have no top-level-only constraint.
- Are `requirement` and `acceptance-criteria` Structure Kind or Document Kind? They describe embedded constructs inside a Feature README, but their format pages have their own adherence-footer-able URL (for educational linking). Current position: Structure Kind — the URL exists, but consumers don't write standalone files.

---
*This document follows the https://specscore.md/idea-specification*
