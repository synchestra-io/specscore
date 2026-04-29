# Ideas Workflow

An **Idea** is a pre-spec, lintable one-pager: a typed artifact for the "should we?" conversation that happens before a Feature's "what?". It captures a problem, a recommended direction, the MVP scope, and the dealbreaker assumptions that must hold for the direction to be worth pursuing.

For the authoritative schema and validation rules, see the [Idea feature spec](../spec/features/idea/README.md). This page is the narrative guide.

## When to Use an Idea

Reach for an Idea when:

- The concept is vague and you need to sharpen it.
- Multiple directions are plausible and none is obviously best.
- You want to stress-test assumptions before committing to a design.
- You need a reviewable artifact to get stakeholder buy-in.
- A retrospective may later need to audit the thinking behind a shipped Feature.

Skip the Idea step when the Feature is obvious, the direction is high-conviction, and the work can go straight to `specscore:design` or a hand-authored Feature. SpecScore does not force an Idea to exist for every Feature.

## Anatomy of an Idea

An Idea is a single markdown file at `spec/ideas/<slug>.md` with a fixed section schema. Every section is required; lint enforces the structure.

| Section | Purpose |
|---|---|
| Title `# Idea: <Name>` | Dispatch key for the Idea lint rule set. |
| Header fields | `Status`, `Date`, `Owner`, `Promotes To`, `Supersedes`, optional `Related Ideas`, `Archive Reason`. |
| Problem Statement | One "How Might We…" sentence. |
| Context | Triggering observation, related specs, prior art. |
| Recommended Direction | The direction you are recommending, and why over the alternatives. |
| Alternatives Considered | At least two directions that lost, with reasons. |
| MVP Scope | The single job the MVP nails — timeboxed, not feature-listed. |
| Not Doing (and Why) | Explicit exclusions with reasons. Must be non-empty. |
| Key Assumptions to Validate | Must-be-true, Should-be-true, Might-be-true. At least one Must-be-true. |
| SpecScore Integration | Features this would create or affect, and dependencies. |
| Open Questions | Questions blocking promotion. Empty state: "None at this time." |

A short example:

```markdown
# Idea: Offline Mode

**Status:** Draft
**Date:** 2026-04-15
**Owner:** alex@synchestra.io
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** depends_on:local-storage-schema

## Problem Statement
How might we let users keep working when their network drops without losing data or breaking sync?

## Context
Support tickets show a 12% weekly rate of "lost my work" complaints traced to transient connectivity…
```

## Creating an Idea — Three Paths

All three paths produce the same artifact. Validation does not care how it was authored — only that `specscore lint` passes.

### 1. Manual

Create `spec/ideas/<slug>.md` in any editor, follow the schema above, and run `specscore spec lint` until it is clean. Appropriate when you already know the shape of the Idea and just need to write it down.

### 2. `specscore new idea <slug>` (CLI scaffolder)

```
specscore new idea offline-mode
specscore new idea offline-mode --title "Offline Mode" --owner alex@synchestra.io \
    --hmw "How might we let users keep working offline?" \
    --not-doing "conflict resolution UI — v2"
specscore new idea offline-mode -i          # interactive TUI
specscore new idea offline-mode --force     # overwrite existing
```

The scaffold is lint-clean on exit regardless of which flags you supply. Missing content is filled with inline HTML-comment prompts that describe what belongs there; these are explicitly allowed by the placeholder-detection rule. Flag values replace the corresponding prompt with real text.

### 3. `spec-studio:ideate` skill

For Claude Code and compatible AI tools, the [`spec-studio:ideate`](https://github.com/synchestra-io/spec-studio/tree/main/skills/ideate) skill runs a three-phase divergent/convergent refinement process (Understand & Expand → Evaluate & Converge → Crystallize). It asks sharpening questions, generates variations, stress-tests directions on user value / feasibility / differentiation, and writes the final artifact to `spec/ideas/<slug>.md`.

The skill delegates file creation to `specscore new idea` when the CLI is installed and falls back to writing the file directly otherwise.

## The Lifecycle

```
Draft → Under Review → Approved → Specified → Archived
```

- **Draft** — first lint-clean write. Author is still iterating.
- **Under Review** — author has requested feedback from stakeholders.
- **Approved** — Recommended Direction has been approved. Ready for promotion.
- **Specified** — at least one Feature lists the Idea in its `**Source Ideas:**` field. **This transition is derived**: tooling sets it, not the author. Writing `Status: Specified` by hand without a matching Feature reference is a lint error.
- **Archived** — abandoned or superseded. The file is **moved** to `spec/ideas/archived/<slug>.md` and `**Archive Reason:**` becomes required.

When every Feature that references an Idea is removed or unlinks it, the Idea reverts from `Specified` back to `Approved`. An Idea is never stuck in `Specified` without a live reference.

## Promotion to Features

Ideas and Features cross-reference in a **many-to-many** relationship.

- A Feature carries the authoritative link via its `**Source Ideas:**` header field listing one or more Idea slugs.
- Each referenced Idea's `**Promotes To:**` is the derived reverse index — managed by tooling.

When a Feature is created with `**Source Ideas:** offline-mode, local-storage-schema`, tooling appends the Feature slug to each referenced Idea's `**Promotes To:**` and transitions each from `Approved` to `Specified`.

`specscore spec lint` (strict) fails on any drift between a Feature's `**Source Ideas:**` and the corresponding Idea's `**Promotes To:**` / `**Status:**`. `specscore spec lint --fix` rewrites the Idea headers in place to reconcile. Run `--fix` locally before committing; CI rejects drift.

## Related Ideas

The optional `**Related Ideas:**` header carries typed links to other Ideas. The vocabulary is fixed:

- `depends_on` — this Idea only works if the referenced one is also pursued.
- `alternative_to` — same problem, different approach; normally at most one is promoted.
- `extends` — this Idea builds on the scope of another; both can coexist.
- `conflicts_with` — incompatible directions; promoting both creates contradictory Features.

Format: `**Related Ideas:** depends_on:payment-rails-audit, alternative_to:single-click-checkout`. Cycles in `depends_on` are permitted — pre-spec thinking is legitimately interdependent.

## Day-to-Day Commands

```
specscore new idea <slug>              # scaffold a new Idea
specscore new idea <slug> -i           # interactive scaffold
specscore spec lint                    # validate the spec tree
specscore spec lint --fix              # auto-repair drift (idea ↔ feature)
```

## Further Reading

- [Idea feature spec](../spec/features/idea/README.md) — authoritative schema, lint rules, acceptance criteria.
- [Feature feature spec](../spec/features/feature/README.md) — how `**Source Ideas:**` drives the Idea's derived state.
- [`spec-studio:ideate` skill](https://github.com/synchestra-io/spec-studio/tree/main/skills/ideate) — recommended AI-assisted authoring path.
