# Impresario — Project Lifecycle Skills Architecture

**Status:** Approved
**Date:** 2026-04-09
**Author:** Alexander Trakhimenok (with Claude)

## Decision

**Impresario** is the lifecycle producer that drives projects from idea to ship. It provides both **skills** (structured workflows) and **personas** (specialist roles like Marketing VP, Troubleshooter, Security Auditor) that can be called upon at any stage. It consumes domain-specific skills from SpecScore and Synchestra to orchestrate the full development lifecycle.

## Name Origin

- **Impresario** (Italian) = the organizer/producer of a musical or theatrical production
- From *impresa* (enterprise, undertaking)
- The impresario doesn't play instruments — they assemble the talent, manage the show, and make it happen

Single S. Correct Italian/English spelling. No footnotes needed.

## The Three-Layer Architecture

```
┌─────────────────────────────────────────────────────┐
│                    IMPRESARIO                        │
│            (lifecycle producer)                      │
│                                                     │
│  impresario:idea → impresario:spec → impresario:plan│
│  → impresario:build → impresario:test               │
│  → impresario:review → impresario:ship              │
└────────────┬──────────────────────┬─────────────────┘
             │                      │
             ▼                      ▼
┌────────────────────┐  ┌────────────────────────────┐
│   SPECSCORE SKILLS │  │     SYNCHESTRA SKILLS      │
│   (Maestria)       │  │                            │
│                    │  │                            │
│ Spec authoring     │  │ Agent orchestration        │
│ Spec linting       │  │ Worktree isolation         │
│ Ideation           │  │ Parallel agent dispatch    │
│ Plan structure     │  │ Subagent coordination      │
│ Requirement writing│  │ Code review workflows      │
│ Scenario design    │  │ Branch management          │
│ Source ref linking  │  │ CI/CD integration          │
└────────────────────┘  └────────────────────────────┘
```

**Impresario** owns the "what happens next" — the lifecycle sequence.
**SpecScore Skills (Maestria)** own the "how to write/validate specs."
**Synchestra Skills** own the "how to orchestrate agents and workflows."

### Personas

Beyond skills, Impresario provides **personas** — specialist roles the agent can adopt for specific tasks. An impresario assembles the right talent for each scene:

| Persona | When summoned |
|---|---|
| **Marketing VP** | Naming, positioning, brand strategy, go-to-market |
| **Troubleshooter** | Debugging, root cause analysis, incident response |
| **Security Auditor** | Threat modeling, OWASP review, vulnerability assessment |
| **Code Reviewer** | Code quality, architecture review, Staff Engineer lens |
| **Test Engineer** | Test strategy, coverage analysis, QA perspective |
| **Product Owner** | Requirements refinement, acceptance criteria, prioritization |

Personas differ from skills: a **skill** is a workflow (steps to follow), a **persona** is a perspective (lens to think through). Impresario composes both — summon the right persona, run the right skill.

## How It Works

The impresario calls on domain skills at each lifecycle stage:

| Lifecycle Stage | Impresario Command | SpecScore Skills Used | Synchestra Skills Used |
|---|---|---|---|
| **Ideation** | `impresario:idea` | Ideation, problem framing | — |
| **Specification** | `impresario:spec` | Spec authoring, linting, validation | — |
| **Planning** | `impresario:plan` | Plan structure, task breakdown | Agent dispatch strategy |
| **Implementation** | `impresario:build` | Source ref linking | Worktrees, parallel agents, TDD |
| **Testing** | `impresario:test` | Scenario verification, AC validation | Test orchestration |
| **Review** | `impresario:review` | Spec compliance checking | Code review workflows |
| **Ship** | `impresario:ship` | — | CI/CD, branch management |

## Command Namespace

```
impresario:idea         # ideation — refine a raw concept into a structured idea
impresario:spec         # specification — author and validate SpecScore specs
impresario:plan         # planning — break specs into executable plans and tasks
impresario:build        # implementation — execute tasks with agent coordination
impresario:test         # verification — validate against acceptance criteria
impresario:review       # code review — spec compliance + code quality
impresario:ship         # launch — pre-launch checks, deployment, release
```

This parallels existing namespace patterns in the ecosystem:

```
specscore:lint          # SpecScore CLI — spec validation
specscore:<ref>         # SpecScore source annotations
impresario:<stage>      # Impresario — lifecycle commands
```

## The Orchestra Metaphor — Complete Picture

| Product | Role in the Orchestra | What It Does |
|---|---|---|
| **Impresario** | The producer | Runs the show end-to-end; decides what happens when |
| **Synchestra** | The orchestra | Coordinates the performers (agents) in concert |
| **SpecScore** | The musical score | Defines what gets performed (specifications) |
| **Maestria** | The performers' skill | SpecScore domain expertise (spec quality, authoring) |
| **Rehearse** | The rehearsal | Tests the performance before opening night |

The impresario reads the score (SpecScore), calls on the orchestra (Synchestra) and its performers' mastery (Maestria), runs rehearsals (Rehearse), and produces the show.

## Comparison to Competitors

| | **Impresario** | **Superpowers** (obra) | **Agent-Skills** (addyosmani) |
|---|---|---|---|
| **Scope** | Full lifecycle producer consuming domain skills | Development workflow skills | Engineering discipline skills |
| **Architecture** | Three-layer (producer → domain skills → tools) | Flat (14 peer skills) | Flat (21 peer skills, 6 phases) |
| **Metaphor** | "The producer runs the show" | "You get powers" | "Skills for agents" |
| **Composability** | Impresario composes SpecScore + Synchestra skills | Skills are standalone | Skills are standalone |
| **Spec integration** | Native — built on SpecScore | None | Ad-hoc markdown specs |

The key differentiator: Impresario is a **composer of skills**, not just a collection. It orchestrates domain-specific skills from multiple sources into a coherent lifecycle.

## Borrowing the Best

From this architecture, we can integrate the best ideas from both competitors:

**From Agent-Skills:**
- Reference checklists (security, performance, accessibility) → become SpecScore Skills (Maestria)
- Anti-rationalization tables → embedded in Impresario lifecycle commands
- Agent personas (reviewer, auditor, tester) → available at `impresario:review`
- `docs/ideas/` convention → already adopted

**From Superpowers:**
- Brainstorming via Socratic questioning → merged into `impresario:idea`
- Git worktree isolation → available at `impresario:build`
- Subagent-driven development → available at `impresario:build`
- Verification before completion → built into every `impresario:` transition

## Alternatives Considered

- **Maestria for everything** — strong standalone name but doesn't convey "lifecycle producer"
- **Skillestra** — clear family tie to Synchestra but collision risk in conversation
- **Impressario** (double-s) — "impress" pun, but: typo friction, spelling corrections, broken imports. The real word is already better.
- **GrandMaestro** — overkill for a skills framework name

## Domain

No dedicated domain needed at launch. Competitors (Superpowers, Agent-Skills) succeed with GitHub URLs alone. Developers find skills through package managers, plugin marketplaces, and GitHub search.

- Primary home: `github.com/synchestra-io/impresario`
- Marketing page if needed: `synchestra.io/impresario`
- Candidates evaluated and parked: `impresario.show` (available, strong metaphor but reads as demo site), `impresario.run` (available, developer-friendly but generic)

## Open Questions

1. Should Impresario be a separate repo or live within SpecScore?
2. Package name: `impresario`, `@synchestra/impresario`?
3. Can third parties create Impresario-compatible skill sets? Plugin architecture?
4. How does Impresario interact with existing Superpowers installation — replacement, wrapper, or coexistence?
