# Examples: Reading an Idea Well

This page calibrates taste. Good Ideas vary — they disagree about directions, weight alternatives differently, and frame problems in their own voice. Bad Ideas fail in predictable ways: the same handful of shapes, again and again.

Rather than enumerate rules in the abstract, this page dissects one deliberately bad Idea inline, shows what a repair looks like, and points you at the real-world good Ideas in the repository and the companion skill.

## How to read this page

1. Read the [bad Idea](#a-bad-idea-ai-dashboard) in the first code block below, cold. Notice what makes you uncomfortable.
2. Work through the annotated walk-through. Each subsection quotes the offending fragment and names the failure.
3. Read the [repaired Idea](#the-repaired-idea). It is the same concept, same author, same afternoon — just thought through.
4. When you write your own Idea, run `specscore spec lint`. Bad Ideas fail structurally before they fail on judgment; fix the structure first and the judgment pressure becomes visible.

Do not skim the bad Idea for laughs. It is written to be plausible. If you have never written an Idea this sloppy, you have been lucky, or you have not reviewed enough of other people's.

## A bad Idea: `ai-dashboard`

Saved as `spec/ideas/ai-dashboard.md`. Ignore, for a moment, that `specscore spec lint` rejects it with nine errors. Read it as prose.

```markdown
# AI Dashboard Idea

**Status:** Specified
**Date:** 2026-04-15
**Owner:** product
**Supersedes:** offline-mode
**Related Ideas:** blocks:payment-rails-audit
**Promotes To:** —

## Problem Statement

Users want better visibility into their Ideas and Features.

## Context

We've been hearing from users.

## Recommended Direction

We should leverage AI and best-in-class dashboards to deliver a world-class
experience that empowers users across the product lifecycle.

## Alternatives Considered

None really.

## MVP Scope

- Dashboard with charts
- AI summaries
- Slack integration
- Export to PDF
- Mobile app
- SSO
- Role-based access control

## Not Doing (and Why)

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Should-be-true | Users will like the dashboard | Ask them |
| Might-be-true | AI summaries will be accurate | See what happens |

## SpecScore Integration

TBD.

## Open Questions

- What should it look like?
```

### What's wrong with the title and header

> `# AI Dashboard Idea`

The required format is `# Idea: <Name>` — the `Idea:` prefix is load-bearing because `specscore lint` uses it as the dispatch key to select the Idea rule set. Without it the linter cannot tell whether to apply Idea rules or Feature rules. This is rule `idea-title-format`.

> `**Status:** Specified`
>
> `**Promotes To:** —`

`Specified` is a *derived* status. A human never writes it; it appears when a Feature lands with `**Source Ideas:** ai-dashboard`. Writing it by hand with an empty `Promotes To` is a two-rule failure (`idea-specified-requires-promotion` and `idea-sync-lint-strict`) and, more importantly, it lies about the artifact's state. The author is trying to skip the promotion gate.

> `**Supersedes:** offline-mode`

`Supersedes:` must point to an Idea that actually exists under `spec/ideas/archived/` and has `Status: Archived`. Pointing it at a non-existent or non-archived slug (rule `idea-supersedes-target-archived`) is the pre-spec equivalent of a dangling pointer: any tool that traces lineage will trip over it.

> `**Related Ideas:** blocks:payment-rails-audit`

The `Related Ideas` vocabulary is fixed at four relationships: `depends_on`, `alternative_to`, `extends`, `conflicts_with`. `blocks` is a reasonable English word and not one of them (rule `idea-related-ideas-format`). Do not invent relationships; the vocabulary is deliberately small so consumers of the Idea graph can enumerate edges.

The header field order is also wrong — the canonical order is `Status, Date, Owner, Promotes To, Supersedes, Related Ideas` — which `idea-header-fields` catches.

### What's wrong with the Problem Statement

> Users want better visibility into their Ideas and Features.

Three problems stacked in one sentence:

- **No How-Might-We framing.** The HMW form is idiomatic because it forces the author to name a user, a desired outcome, and an implicit constraint. The declarative "users want X" form smuggles in both the diagnosis and the solution without examining either (rule `idea-hmw-framing`, warning).
- **The user is unspecified.** "Users" is everyone; therefore no one. Is this the solo user running `specscore` locally? The stakeholder reviewing Ideas at a weekly meeting? Synchestra dashboards consumed by a PM? Each implies a different artifact.
- **"Visibility" is a vitamin, not a painkiller.** No one is bleeding. The Problem Statement does not name a moment where the lack of visibility costs someone something concrete.

A repaired Problem Statement would sound like: *"How might we let a solo SpecScore user see, at a glance, which Ideas have gone stale — so that weekly Idea review takes minutes instead of a manual sweep of every file?"* It names the user (solo SpecScore user), the outcome (see stale Ideas at a glance), and the cost being paid today (manual sweep).

### What's wrong with the Context

> We've been hearing from users.

This is filler. Context exists so the reader can audit the triggering observation without trusting the author's memory. Good Context cites: a support ticket rate, a specific conversation, a prior design doc, an adjacent Feature. "We've been hearing from users" is a claim that cannot be falsified and therefore cannot be learned from.

### What's wrong with the Recommended Direction

> We should leverage AI and best-in-class dashboards to deliver a world-class experience that empowers users across the product lifecycle.

Take a position, or do not write a Recommended Direction. This sentence takes none. It is a compressed list of buzzwords ("AI", "dashboards", "world-class", "empowers", "lifecycle") stitched together with connective tissue. A reviewer cannot disagree with it, because there is nothing specific to disagree with — and that is the tell. A direction that no one could plausibly oppose is not a direction, it is a mood.

The Recommended Direction should choose one approach, name its tradeoff against the alternatives, and say why the author is recommending it *over* the rejected paths. "Ship `specscore ideas stale` as a CLI command, not a web UI, because the bottleneck is surfacing the data, not rendering it" is a direction. You can argue against it. That is what makes it useful.

### What's wrong with the Alternatives Considered

> None really.

Spec requires at least two alternatives with reasons they lost (section schema plus rule `idea-required-sections` covers its presence). Writing "None" is the author confessing that they stopped thinking after the first idea survived first contact with their own head. Every direction has at least two plausible alternatives; if you cannot name them, you have not yet understood the problem.

### What's wrong with the MVP Scope

> - Dashboard with charts
> - AI summaries
> - Slack integration
> - Export to PDF
> - Mobile app
> - SSO
> - Role-based access control

MVP Scope is *the single job the MVP nails*, timeboxed. A seven-line feature list is the opposite: it is a wishlist that has been mislabelled. If everything is MVP, nothing is. The shape the section wants is a sentence: *"One command — `specscore ideas stale` — that identifies Ideas exceeding per-status age thresholds and prints them as a table. Target: usable end-to-end within one week."* That is timeboxed and has a sharp edge. A reviewer can hold you to it.

Feature lists also pre-empt the downstream Feature decomposition. That is the job of whoever writes the Features from the Idea; the Idea's job is to commit to the *first useful slice* and defer the rest.

### What's wrong with the Not Doing section

> ## Not Doing (and Why)
>
> *(section empty)*

Saying no to 1,000 things is a tenet, not a decoration. An empty Not Doing list means the author has not yet confronted the real tradeoff: what they are *giving up* in order to ship the MVP. Rule `idea-not-doing-non-empty` rejects the file. Good: the rejection is the point. If you cannot name three things you are deliberately skipping, you have not sharpened the Idea enough to promote it.

### What's wrong with the Assumptions table

> | Tier | Assumption | How to validate |
> |------|------------|-----------------|
> | Should-be-true | Users will like the dashboard | Ask them |
> | Might-be-true | AI summaries will be accurate | See what happens |

Two failures:

- **No Must-be-true row.** `idea-must-be-true-present` rejects this. Must-be-true assumptions are *dealbreakers* — the assumption whose falsification kills the whole direction. Without one, the Idea has not been stress-tested. The author has not asked "what single fact, if I learned it tomorrow, would make this Idea not worth pursuing?"
- **Validation plans that are not plans.** "Ask them" and "See what happens" are verbal shrugs. A real validation plan names an instrument (a metric, a prototype, an experiment, a pair of design-partner repos) and a threshold ("if fewer than half the Ideas it surfaces are ones the partner would have flagged manually, the direction is wrong").

### What's wrong with the SpecScore Integration

> TBD.

The section exists so the Idea can be traced into and out of the Feature graph. "TBD" removes the only handle downstream tooling has. At minimum, name the Features this would create (or explicitly `TBD at design time` with a note on why), the existing Features affected, and the dependencies. The bad Idea punts on the one piece of machine-readable metadata the section was invented to carry.

### What's wrong with the Open Questions

> - What should it look like?

Open Questions are questions that *block promotion to a Feature*. "What should it look like?" is a design question the Feature is meant to answer. If every question is at that altitude, the Idea is not ready for review; it is still a concept. The questions that belong here are the ones a reviewer needs to answer before approving the Recommended Direction — typically about assumptions, scope cuts, or constraints the author could not resolve alone.

## The repaired Idea

Same author, same afternoon, same underlying impulse — thought through. This file lints clean.

```markdown
# Idea: AI Dashboard

**Status:** Draft
**Date:** 2026-04-15
**Owner:** alex@synchestra.io
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we let a solo SpecScore user see, at a glance, which Ideas in their
repo have gone stale — so that weekly Idea review takes minutes instead of a
manual sweep of every file?

## Context

Two design partners running SpecScore on 40+ active Ideas each report spending
~30 minutes every Monday re-reading Idea files to find ones that have drifted:
`Under Review` for more than a week, `Approved` with no `Promotes To` after a
fortnight, or `Draft` untouched for a month. The data is already in the header
fields; no one has a view over it.

## Recommended Direction

Ship a single `specscore ideas stale` CLI command that prints a table of Ideas
exceeding configurable age thresholds per status. No web UI, no daemon, no
"dashboard". The command reads existing header fields, applies thresholds from
`specscore-project.yaml`, and writes plain text plus optional JSON. Chosen over
a web dashboard because the bottleneck is surfacing information, not rendering
it — a `grep`-able table composes with the user's existing tooling.

## Alternatives Considered

- **Web dashboard served by `specscore serve`.** Rejected: builds, hosts, and
  maintains a UI for a problem that a 20-line table solves. Defers shipping by
  weeks.
- **GitHub Action that posts a weekly comment on a tracking issue.** Rejected:
  ties the feature to one forge and requires setting up CI for users who
  haven't. Can be layered on top of the CLI later.

## MVP Scope

One command — `specscore ideas stale` — that identifies Ideas exceeding
per-status age thresholds and prints them as a table. Target: usable end-to-end
by the two design partners within one week of work.

## Not Doing (and Why)

- Web UI — a CLI table solves the surfacing problem at 1% of the cost; a UI
  can come after we know what columns matter.
- Configurable thresholds per-Idea — project-level thresholds only. Per-Idea
  overrides are premature until we see whether anyone needs them.
- Notifications, email, Slack integration — out of scope until the read-only
  view is validated.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Age-since-last-edit is a good enough staleness signal; users will not demand semantic staleness on day one. | Run the command against the two partners' repos for two weeks; compare against Ideas they flag manually. |
| Should-be-true | Default thresholds (Draft 30d, Under Review 7d, Approved 14d) fit most teams without tuning. | Ship the defaults; measure how many teams override them after one month. |
| Might-be-true | Users will also want a `--json` mode for piping into other tools. | Offer it; count invocations with `--json` versus default text. |

## SpecScore Integration

- **New Features this would create:** `ideas-stale-command`.
- **Existing Features affected:** `project-definition` gains an optional
  `ideas.staleness` block for thresholds.
- **Dependencies:** none — the command reads artifacts that already exist.

## Open Questions

- Should the command exit non-zero when stale Ideas are found, for CI use, or
  always exit zero and leave CI integration to a later Idea?
```

### What changed, and why

- **Title now uses `# Idea: …`.** Dispatch key for the lint rule set.
- **Status is `Draft`**, not `Specified`. The author writes `Draft`; tooling writes `Specified` when a Feature is created that references this slug.
- **Header fields in canonical order.** `Status, Date, Owner, Promotes To, Supersedes, Related Ideas`. The fabricated `Supersedes` and invented `blocks:` relationship are gone.
- **Problem Statement is a single HMW sentence** that names a user, an outcome, and the cost being paid today.
- **Recommended Direction takes a position** (CLI, not web UI) and names the tradeoff out loud.
- **Alternatives Considered lists two directions with reasons they lost.** A reviewer now has something to argue with.
- **MVP Scope is one sentence, timeboxed.** The seven-item wishlist is gone. The deferred items live in Not Doing.
- **Not Doing has three explicit exclusions with reasons.** Each is a commitment; later Ideas can lift them with deliberate intent.
- **Assumptions table has a Must-be-true row** naming a dealbreaker, plus real validation plans with thresholds.
- **SpecScore Integration names the Feature** that would be created and the existing Feature it touches.
- **Open Questions carries a genuine blocker** — a scope decision the author is deferring to review, not a design question.

## Finding good Ideas

This page shows what bad looks like. For what good looks like, read real ones.

- **`spec/ideas/`** in this repository is the canonical reference once real Ideas exist. At the time of writing, the directory contains no active Ideas — the first real ones, when written, become the ongoing good-examples library. Prefer a real in-repo Idea over any synthetic example.
- **The `specscore:ideate` skill's `references/examples.md`** at `/home/ai/projects/synchestra-io/ai-plugin-sdd/skills/specscore-ideate/references/examples.md` contains external ideation-session examples — full three-phase sessions (Understand & Expand → Evaluate & Converge → Crystallize) that terminate in a complete Idea artifact. Useful when you want to see the *process* that produced a good Idea, not just the artifact.

When a real Idea in `spec/ideas/` has aged well through a Feature promotion and a retrospective, cite it from this page.

## Common failure modes — self-audit checklist

Scan this before sending an Idea for review. If you can answer "yes" to any of these, fix it first.

- **Title does not start with `# Idea: `.** Fix the dispatch key.
- **Status is `Specified` without a Feature that references this slug.** `Specified` is derived; never author it.
- **Status is `Archived` but the file is still at `spec/ideas/<slug>.md`.** Move it to `spec/ideas/archived/`.
- **`Supersedes:` points to an Idea that does not exist in `archived/` with `Status: Archived`.** Dangling pointer.
- **A `Related Ideas:` entry uses a relationship other than `depends_on`, `alternative_to`, `extends`, `conflicts_with`.** The vocabulary is fixed; invent at your peril.
- **Header field order is not canonical.** `Status, Date, Owner, Promotes To, Supersedes, Related Ideas`.
- **Problem Statement is not an HMW sentence, or does not name a specific user.** "Users want X" is almost always a smell.
- **Problem is a vitamin, not a painkiller.** No one is bleeding; no one will miss the Idea if it ships six months late.
- **Context cannot be falsified.** "We've been hearing from users" is not Context.
- **Recommended Direction takes no position.** If a reviewer cannot plausibly disagree, rewrite it.
- **Alternatives Considered has fewer than two alternatives, or says "none really".** Think harder.
- **MVP Scope is a bulleted feature list instead of a single timeboxed job.** Cut until one job remains.
- **Not Doing is empty or contains only one throwaway item.** Name three things you are giving up to ship.
- **Assumptions table has no Must-be-true row.** You have not stress-tested the direction.
- **Validation plans are verbal shrugs** ("ask them", "see what happens"). Name an instrument and a threshold.
- **SpecScore Integration is `TBD` with no detail.** At minimum, name the Features this would create or mark them `TBD at design time` with a reason.
- **Open Questions contain design questions instead of promotion blockers.** Move design questions to the downstream Feature.

If every check passes and `specscore spec lint` is clean, send for review.
