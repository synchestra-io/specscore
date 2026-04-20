# Idea: Decision and Decisions Index

**Status:** Specified
**Date:** 2026-04-20
**Owner:** alexander.trakhimenok@gmail.com
**Promotes To:** decision, decisions-index
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we capture the moment a team chose one option over others — in a lint-clean, traceable artifact — so that future engineers (human and AI) can audit the "why we did it this way" without spelunking chat logs or commit history?

## Context

SpecScore already has artifacts for what will be built (`feature`), how it will be built (`plan`), and whether it is worth building (`idea`). It has no artifact for **which of several paths was chosen and why**.

Ideation produces alternatives; some surface inside an Idea's `Alternatives Considered` section. But that section is *pre-choice* — "here's what we thought about." It disappears into an `archived/` directory once the Idea is Specified, and it never captures choices made outside ideation (mid-implementation trade-offs, vendor constraints, reversals of earlier decisions).

The established convention for this problem is the Architecture Decision Record (ADR), popularized by Michael Nygard. ADRs are durable, immutable-once-accepted, and numbered for citation. This Idea adopts ADR discipline inside SpecScore conventions (markdown-body metadata, REQ blocks, AC blocks, adherence footer, companion index feature).

Two features are proposed together: `decision` (the Document Kind) and `decisions-index` (the Index Kind), following the same pairing pattern as `idea`/`ideas-index`, `plan`/`plans-index`, etc.

## Recommended Direction

Build a `decision` feature modeled structurally on `idea` but with ADR semantics:

- **Location:** `spec/decisions/NNNN-<slug>.md` (active) or `spec/decisions/archived/NNNN-<slug>.md` (superseded). The `NNNN-` prefix is a zero-padded sequence number — this is a deliberate break from the slug-only convention used elsewhere in SpecScore, matching universal ADR practice and giving PRs/commits a short citation (`D-0007`).
- **Immutability:** once `Status: Accepted`, the body is frozen. Course changes are made by creating a new Decision with `**Supersedes:**` pointing at the old one. The old one's status becomes `Superseded` and it moves to `archived/`.
- **Schema distinguishers from ADR-in-the-wild:**
  - `## Declined Alternatives` is a structured list (one entry per rejected option: name, one-line pitch, rejection reason) — not a prose paragraph. Future engineers asking "why didn't we just do Y?" find the answer already written.
  - `## Observed Consequences` is an append-only dated log — the section authors are *meant* to update after the decision plays out. Appends do not break immutability.
  - `**Source Idea:**` header field is optional and back-links a Decision to the Idea that triggered it, when one exists.
- **Index feature** (`decisions-index`) mirrors `ideas-index`: active index at `spec/decisions/README.md` (table: Number, Decision, Status, Date, Tags, Affected), archived index at `spec/decisions/archived/README.md` (chronological list, oldest first).

Both features declare an adherence footer per the [Adherence Footer feature](../features/adherence-footer/README.md), with URLs `https://specscore.md/decision-specification` and `https://specscore.md/decisions-index-specification`.

## Alternatives Considered

**Simpler three-section Decision with no lifecycle.** Context / Options / Choice, no status, no index, no supersede rules. Lost because it reduces a Decision to a plain doc, undermining SpecScore's core bet: lint-clean typed artifacts with traceability. Git history isn't a substitute for typed state.

**Extend `idea` to carry "finalized decision" state.** An Idea with `Status: Decided` would carry the final choice. Lost because Idea is pre-spec thinking, Archived once Specified. A decision made mid-implementation (no originating Idea) would have nowhere to live. Fails the four-trigger coverage we need.

**Slug-only filenames, no number prefix.** Consistent with every other SpecScore Document Kind. Lost because ADRs universally use numbered files for citation (`D-0007`), and the citation value is highest exactly when you're reading old Decisions — which is what this feature is for.

## MVP Scope

Ship both features together with the minimum schema to pass lint:

- `decision` README declaring: title format, header fields (Status, Date, Owner, Tags, Source Idea, Supersedes, Superseded By), required sections (Context, Decision, Rationale, Declined Alternatives, Consequences at decision time, Observed Consequences, Affected Features), status lifecycle (Proposed → Accepted → Superseded), adherence footer, acceptance criteria.
- `decisions-index` README declaring: active-index columns, archived-index chronology rule, completeness rule, adherence footer.
- `specscore new decision <title>` CLI scaffold that assigns the next number and produces a lint-clean skeleton.
- `specscore lint` rules: required sections, status values, supersede target exists, immutability check (body hash unchanged once Accepted — except for `## Observed Consequences` appends), index completeness.

No REQ↔Decision linking in MVP. No automatic migration from existing `Alternatives Considered` sections. No tag vocabulary enforcement.

## Not Doing (and Why)

- **No tag taxonomy.** Tags are free-form in MVP. Canonicalizing (`architecture` vs `arch` vs `technical`) is a second-order problem; solve it after we see real usage.
- **No REQ↔Decision linking.** Tempting, but it's a separate feature that requires a new REQ header field and lint rules. Note as future work; do not gate MVP on it.
- **No retroactive decisions.** Do not write a feature for "backfill ADRs from git history." Retrofitted ADRs are notoriously vapid; making the tool push authors toward them would damage signal-to-noise.
- **No multi-option requirement below 2.** A Decision with zero or one considered option is not a Decision — it's a design note that belongs in the Feature README. Lint enforces ≥2 entries in `Declined Alternatives` OR a mandatory "why no alternatives existed" justification. (This is the `not-doing-non-empty` discipline equivalent.)
- **No consequences-observed enforcement timeline.** We will *not* lint "Observed Consequences is empty after 6 months." That would push authors into performative updates. If the section is consistently empty across the project after a year, drop it in a revision; don't coerce usage.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Teams will write Decisions at the point of choice, not retrofit them. Retrofit ADRs are vapid and would poison the artifact's reputation. | Ship `specscore new decision` with an interactive TUI matching `new idea`'s ergonomics. Measure: first-week usage on SpecScore itself (dogfood). If we don't reach for it during real decisions, no one else will. |
| Must-be-true | Authors can tell the difference between a Decision and an Idea's `Alternatives Considered`. | Include a "When to write a Decision vs put alternatives in an Idea" subsection in the `decision` feature README, with two concrete examples. If first reviewers can't answer correctly after reading it, rewrite. |
| Must-be-true | The `NNNN-slug.md` filename break from slug-only convention is acceptable to SpecScore maintainers. | Surface this as an explicit open question on this Idea before it promotes. Alternative: use slug-only and introduce a `**Number:**` header field. Cheaper either way, but the filename form is conventional in ADR and I'd recommend adopting it. |
| Should-be-true | `Observed Consequences` will actually get updated post-hoc by at least some authors. | After SpecScore ships with ≥5 internal Decisions, check at the 6-month mark. If the section is empty across all of them, drop it in a minor revision — it's theater, not signal. |
| Might-be-true | Decisions will naturally reference REQs and vice versa, making a formal REQ↔Decision linking feature worth building later. | Observe the pattern in the first 10 Decisions. If rationale sections repeatedly cite REQ anchors, build the forward link. If not, drop the idea. |

## SpecScore Integration

- **New Features this would create:** `decision` (Document Kind), `decisions-index` (Index Kind).
- **Existing Features affected:**
  - [Adherence Footer](../features/adherence-footer/README.md) — two new document type URLs to register.
  - [Document Types Registry](../features/document-types-registry/README.md) — two new rows.
  - [Idea](../features/idea/README.md) — optional back-link: an Idea's `Alternatives Considered` section MAY note "Decisions produced: D-0007, D-0008" but this is not required and does not introduce a new managed field in MVP.
- **Dependencies:** `feature`, `adherence-footer`, `document-types-registry`. Blocks nothing; purely additive.

## Open Questions

- **Filename form: numbered-and-slugged (`0007-postgres-over-mongo.md`) or slug-only with a `**Number:**` header field?** Numbered-filename is the ADR convention and gives short citations. Slug-only preserves the existing SpecScore invariant. Recommend numbered-filename; flag as the biggest convention break on this Idea.
- **Status lifecycle: `Proposed → Accepted → Superseded`, or also include `Deprecated` (no successor, just "don't follow this anymore")?** ADR community uses both. Recommend Deprecated be included — sometimes a decision expires without a replacement.
- **Should `Declined Alternatives` require a minimum count of entries (e.g. ≥1)?** Matches `idea`'s `not-doing-non-empty` discipline. Recommend yes; the whole point of the artifact is the trade-off.
- **Where does the `Observed Consequences` append-only rule get enforced?** Lint by body-hash-minus-that-section, or simply by convention? Recommend lint enforces immutability of all other sections once `Accepted`; `Observed Consequences` is the only mutable section.
- **Does `decisions-index` need `Affected Features` as a first-class index column, or is that over-scoping MVP?** Recommend ship without it; add if the volume justifies a faceted index.

---
*This document follows the https://specscore.md/idea-specification*
