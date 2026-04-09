# Skills Ecosystem Analysis: SpecScore + Superpowers + Agent-Skills

**Status:** Draft
**Date:** 2026-04-09
**Author:** Alexander Trakhimenok (with Claude)

## Problem Statement

Three open-source projects occupy adjacent spaces in AI-driven development:

- **SpecScore** (synchestra-io) — specification format and tooling
- **Superpowers** (obra) — agentic workflow skills
- **Agent-Skills** (addyosmani) — engineering discipline enforcement

Each solves a different slice of the same problem: helping AI agents produce production-quality software. This document analyzes where they overlap, where they diverge, and how SpecScore can absorb the best of both while maintaining its unique identity as a specification-first framework.

---

## 1. What Each Project Does

### SpecScore — The Artifact Layer

SpecScore defines *what* gets built. It provides:

- **Formal specification schema** — features, requirements, acceptance criteria, scenarios, plans, tasks, project definitions
- **Machine-readable format** — YAML/Markdown with validation rules
- **Source-to-spec traceability** — `specscore:` annotations linking code back to specs
- **Linting engine** — pluggable rules for spec validation (`specscore lint`)
- **CLI tooling** — real Go binary, not just prompt files
- **Spec hierarchy** — Feature > Requirement > Acceptance Criteria > Scenario, with formal relationships

SpecScore is tool-agnostic. It works standalone or with Rehearse (testing) and Synchestra (orchestration).

### Superpowers — The Workflow Layer

Superpowers defines *how* agents execute. It provides:

- **14 skills** covering the development workflow:
  - `brainstorming` — Socratic questioning to refine requirements before implementation
  - `writing-plans` / `executing-plans` — plan authoring and execution with review checkpoints
  - `test-driven-development` — RED-GREEN-REFACTOR enforcement
  - `systematic-debugging` — structured root cause analysis
  - `using-git-worktrees` — isolated development branches
  - `subagent-driven-development` — parallel agent dispatch with two-stage review
  - `dispatching-parallel-agents` — independent task parallelization
  - `requesting-code-review` / `receiving-code-review` — structured review workflow
  - `finishing-a-development-branch` — merge/PR/cleanup decision framework
  - `verification-before-completion` — evidence-based completion claims
  - `writing-skills` — meta-skill for creating new skills
  - `using-superpowers` — entry point and skill discovery
- **Auto-triggering** — skills activate based on context, not manual invocation
- **Subagent orchestration** — complex tasks decompose into parallel workstreams
- **Multi-platform** — Claude Code, Cursor, Gemini CLI, Codex, OpenCode

### Agent-Skills — The Discipline Layer

Agent-Skills defines *what standards* agents must meet. It provides:

- **21 skills** organized across 6 SDLC phases:
  - **Define:** `idea-refine`, `spec-driven-development`
  - **Plan:** `planning-and-task-breakdown`
  - **Build:** `incremental-implementation`, `test-driven-development`, `context-engineering`, `source-driven-development`, `frontend-ui-engineering`, `api-and-interface-design`
  - **Verify:** `browser-testing-with-devtools`, `debugging-and-error-recovery`
  - **Review:** `code-review-and-quality`, `code-simplification`, `security-and-hardening`, `performance-optimization`
  - **Ship:** `git-workflow-and-versioning`, `ci-cd-and-automation`, `deprecation-and-migration`, `documentation-and-adrs`, `shipping-and-launch`
- **Anti-rationalization tables** — each skill documents common agent excuses with documented rebuttals
- **Reference checklists** — security (OWASP Top 10), performance (Core Web Vitals), accessibility (WCAG 2.1 AA), testing patterns
- **Agent personas** — `code-reviewer` (Staff Engineer lens), `security-auditor`, `test-engineer`
- **Slash commands** — `/spec`, `/plan`, `/build`, `/test`, `/review`, `/code-simplify`, `/ship`
- **Ideas output** — `idea-refine` produces structured one-pagers saved to `docs/ideas/[name].md`

---

## 2. Where SpecScore Is Stronger

| Capability | SpecScore | Superpowers | Agent-Skills |
|---|---|---|---|
| Formal spec schema with validation | Yes (features, requirements, scenarios, acceptance criteria) | No | No (ad-hoc markdown) |
| Source-to-spec traceability | Yes (`specscore:` annotations) | No | No |
| Linting engine | Yes (pluggable rules) | No | No |
| CLI tooling | Yes (Go binary) | No | No |
| Spec hierarchy (Feature > Requirement > AC > Scenario) | Yes | No | No |
| Integration with test framework (Rehearse) | Yes | No | No |
| Integration with orchestrator (Synchestra) | Yes | No | No |
| Machine-readable project definition | Yes (`specscore-project.yaml`) | No | No |

**Key insight:** Neither Superpowers nor Agent-Skills has a structured specification format. Their "specs" are freeform markdown. SpecScore's formal schema is a unique differentiator.

---

## 3. What We're Missing

### 3.1 Pre-Spec Ideation

**Gap:** SpecScore has no place for ideas that haven't matured into formal specs.

Both competitors address this:
- **Agent-Skills** `idea-refine` — divergent/convergent ideation producing a structured one-pager with Problem Statement, MVP Scope, Not-Doing List, Key Assumptions, Open Questions. Output: `docs/ideas/[name].md`.
- **Superpowers** `brainstorming` — Socratic questioning to refine requirements before planning. More interactive and conversation-driven.

**Proposal:** Merge both approaches into a unified SpecScore ideation workflow (see Section 5).

### 3.2 Reference Checklists

**Gap:** No reusable quality checklists.

Agent-Skills ships four:
- `security-checklist.md` — OWASP Top 10 assessment
- `performance-checklist.md` — Core Web Vitals targets
- `accessibility-checklist.md` — WCAG 2.1 AA compliance
- `testing-patterns.md` — test structure and anti-patterns

**Opportunity:** These could become `references/` in SpecScore and eventually inform lint rules (e.g., a security lint rule that checks whether acceptance criteria cover OWASP items).

### 3.3 Agent Personas

**Gap:** No reusable reviewer/auditor personas.

Agent-Skills defines three (`code-reviewer`, `security-auditor`, `test-engineer`) as markdown files in `agents/`. These give the AI a specific lens to review through.

**Opportunity:** An `agents/` directory in SpecScore, with personas that reference SpecScore specs during review (e.g., "verify each requirement in spec X has a passing scenario").

### 3.4 Anti-Rationalization Tables

**Gap:** No defense against agents skipping steps.

Agent-Skills includes tables like:

| Rationalization | Rebuttal |
|---|---|
| "This is too simple for a spec" | Simple features grow. A 15-minute spec prevents hours of rework. |
| "I'll add tests later" | Later never comes. Write the test first. |
| "The security risk is low" | Attackers decide what's low risk, not developers. |

**Opportunity:** Could be added to SpecScore's spec format itself — a `rationalizations` section in feature specs that documents why shortcuts are unacceptable for that feature.

### 3.5 Deployment/Launch Coverage

**Gap:** SpecScore covers idea through task but has nothing for deployment.

Agent-Skills has `shipping-and-launch` (pre-launch checklists, rollout procedures) and `ci-cd-and-automation` (shift-left principles).

**Opportunity:** A `deployment-gate` or `launch-checklist` SpecScore feature, or simply a reference checklist.

### 3.6 ADR (Architecture Decision Records)

**Gap:** No formal decision-record format.

Agent-Skills' `documentation-and-adrs` skill guides creation of Architecture Decision Records. SpecScore has `architecture/` specs but no lightweight ADR format.

**Opportunity:** ADRs could live alongside plans as a SpecScore feature or as a convention within `docs/`.

---

## 4. What Neither Has That We Do

These are SpecScore's durable advantages — things that would be difficult for either project to replicate without becoming a different kind of project:

1. **Formal specification format with validation rules** — they're process docs; we're a schema
2. **Linting engine with pluggable rules** — programmatic spec validation
3. **Source-to-spec traceability** — bidirectional code-to-spec linking
4. **Acceptance criteria as first-class citizens** — testable, referenceable, machine-readable
5. **Scenario format** (Given/When/Then) — structured behavior examples
6. **Integration with Rehearse** — automated spec testing
7. **Integration with Synchestra** — multi-agent orchestration of specs
8. **CLI binary** — real tooling, not just prompts

---

## 5. Proposed: Unified Ideation Workflow

### The Problem with Two Approaches

Agent-Skills' `idea-refine` and Superpowers' `brainstorming` solve the same problem differently:

| | Agent-Skills `idea-refine` | Superpowers `brainstorming` |
|---|---|---|
| **Approach** | Structured 3-phase process (Expand > Evaluate > Sharpen) | Socratic questioning to surface intent |
| **Output** | Written one-pager (`docs/ideas/[name].md`) | Refined understanding leading to a plan |
| **Strength** | Produces a tangible artifact with MVP scope and exclusions | Interactive, uncovers unstated assumptions |
| **Weakness** | Can feel mechanical; may skip user's real intent | No durable artifact; insights lost after conversation |

### Our Melt: SpecScore Ideation

We can merge both into something neither offers alone — ideation that is both interactive *and* produces a formal artifact that feeds directly into SpecScore's spec pipeline.

**Proposed workflow:**

```
Phase 1: Discover (from Superpowers brainstorming)
├── Socratic questioning to surface real intent
├── "What problem does this solve?" / "Who feels this pain?"
├── Challenge unstated assumptions
└── Output: shared understanding of the problem space

Phase 2: Diverge & Converge (from Agent-Skills idea-refine)
├── "How Might We" reframing
├── Generate 5-8 variations across creative lenses
├── Group into 2-3 distinct directions
├── Stress-test against user value, feasibility, differentiation
└── Output: ranked directions with trade-off analysis

Phase 3: Crystallize (unique to SpecScore)
├── Produce structured idea document (see format below)
├── Map to existing SpecScore features if applicable
├── Identify which SpecScore specs this idea would need
├── Flag dependencies on other ideas or in-flight work
└── Output: docs/ideas/[name].md — ready to promote to spec
```

**Idea document format:**

```markdown
# [Idea Name]

**Status:** Draft | Under Review | Approved | Promoted | Archived
**Date:** YYYY-MM-DD
**Author:** [name]

## Problem Statement
[How Might We framing]

## Context
[What triggered this idea; relevant background]

## Recommended Direction
[Why this path over alternatives considered]

## Alternatives Considered
[2-3 other directions explored and why they were deprioritized]

## MVP Scope
[Minimum viable version that tests core assumptions]

## Not-Doing List
[Deliberate exclusions with reasoning]

## Key Assumptions to Validate
[Specific bets that need verification before full investment]

## SpecScore Integration
- **New features needed:** [list of spec/features/ entries this would require]
- **Existing features affected:** [specs that would need updates]
- **Dependencies:** [other ideas or in-flight work this depends on]

## Open Questions
[Unanswered items before this can be promoted to spec]
```

**What makes this unique:** The "SpecScore Integration" section is something neither competitor has. It connects ideation directly to the formal spec pipeline, making ideas *actionable* rather than just documented.

**Promotion path:** When an idea is approved, it gets promoted:
- Status changes to `Promoted`
- Corresponding feature specs are created in `spec/features/`
- The idea document links to the created specs
- The specs link back to the idea as provenance

---

## 6. Integration Roadmap: Using Both Skill Sets

We already use Superpowers (see `docs/superpowers/`). The question is whether to also integrate Agent-Skills and how.

### Complementary Roles

```
Agent-Skills                    SpecScore                     Superpowers
(discipline)                    (artifacts)                   (workflow)
─────────────                   ─────────────                 ─────────────
idea-refine          ──►        docs/ideas/[name].md
                                      │
                                      ▼
spec-driven-dev      ──►        spec/features/[name]/    ◄── brainstorming
                                      │
planning-and-task-   ──►        spec/features/[name]/    ◄── writing-plans
  breakdown                       plans/
                                      │
                                      ▼
                                                          ◄── executing-plans
incremental-impl     ──►        [implementation]          ◄── subagent-driven-dev
test-driven-dev                                           ◄── using-git-worktrees
                                      │
                                      ▼
code-review          ──►        [review]                  ◄── requesting-code-review
security-and-hardening                                    ◄── receiving-code-review
performance-opt
                                      │
                                      ▼
shipping-and-launch  ──►        [deployment]              ◄── finishing-a-dev-branch
ci-cd-and-automation                                      ◄── verification-before-completion
```

**Superpowers** handles *how agents work* (parallel dispatch, worktrees, review loops).
**Agent-Skills** handles *what agents check* (security, performance, accessibility, launch readiness).
**SpecScore** handles *what gets written* (specs, plans, tasks) and *validates it* (lint, source refs).

### Proposed Integration Model

1. **Keep Superpowers as the workflow engine** — it already orchestrates our development sessions
2. **Adopt Agent-Skills selectively** — not as a plugin, but by incorporating its best ideas:
   - `references/` directory with quality checklists (adapted, not copied)
   - `agents/` directory with review personas that reference SpecScore specs
   - Anti-rationalization patterns embedded in our skill documentation
   - `docs/ideas/` directory (already created) with our unified ideation format
3. **SpecScore remains the source of truth** — specs, validation, and traceability are ours; skills from either ecosystem feed into and consume SpecScore artifacts

---

## 7. Action Items

| Priority | Action | Rationale |
|---|---|---|
| P0 | Create `docs/ideas/` directory with ideation template | Gives pre-spec ideas a home (this document is the first) |
| P1 | Design unified ideation skill merging brainstorming + idea-refine | Our unique melt; neither competitor offers this |
| P1 | Add `references/` directory with security, performance, accessibility checklists | Low effort, high value; can later inform lint rules |
| P2 | Add `agents/` directory with SpecScore-aware review personas | Review quality improvement |
| P2 | Evaluate adding anti-rationalization sections to SpecScore feature specs | Defense against agents cutting corners |
| P3 | Design `deployment-gate` feature or checklist | Closes the idea-to-ship gap |
| P3 | Evaluate lightweight ADR format compatible with SpecScore | Decision provenance |

---

## 8. Open Questions

1. Should the unified ideation skill live in SpecScore or as a Superpowers skill that references SpecScore?
2. Should reference checklists be static markdown or eventually become machine-readable lint rule inputs?
3. How do we version/deprecate ideas? Is `Archived` status sufficient or do we need explicit expiry?
4. Should agent personas be part of SpecScore format (formalized) or stay as convention (`agents/` directory with freeform markdown)?
