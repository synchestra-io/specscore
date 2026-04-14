# SpecScore Marketing Plan — Positioning, Website, Content & Launch

**Date:** 2026-04-14
**Status:** Approved
**Approach:** Standard-first (B) with ecosystem hint (C)

## Context

SpecScore is being extracted from Synchestra as an independent open-source specification framework. It needs its own marketing identity, website (specscore.md), content strategy, and launch plan. The ecosystem brand system, voice, and visual identity are defined in `synchestra-marketing/` and apply to SpecScore.

### Existing Assets

- Brand voice, color system, typography, illustration style: `synchestra-marketing/branding/`
- Ecosystem positioning and audience tiers: `synchestra-marketing/POSITIONING.md`
- Messaging hierarchy: `synchestra-marketing/ECOSYSTEM.md`
- Extraction plan: `synchestra-marketing/EXTRACTION_PLAN.md`

### Key Decisions Made During Brainstorm

- **Domain:** specscore.md (inspired by obsidian.md — Markdown as brand signal, audience is increasingly Markdown-literate including BAs via AI skills adoption)
- **Go-to-market:** Developer-first + AI-skills-first (complementary — devs adopt the tool, AI skills produce/consume the format)
- **Competition:** Inertia (no formal specs) and AI-native alternatives (SpecFlow, OpenSpec)
- **Core differentiator:** Structure without bureaucracy, enforced by tooling (linter + LSP) — SpecScore validates specs, not just defines them
- **Brand voice:** Inherited from Synchestra — technical but warm, confident and declarative, orchestra metaphor used selectively

## Section 1: Positioning & Core Message

**Domain:** specscore.md

**Tagline:** "The open specification standard for AI-driven development"

**Elevator pitch:**

> SpecScore is an open specification format that makes requirements machine-readable without making them human-unreadable. It's Markdown and YAML — version-controlled, portable, no vendor lock-in. A linter and LSP catch ambiguity before your agents do. Use it standalone or with any orchestration tool.

**Primary positioning:** SpecScore is a format and standard first, tooling second. The CLI validates specs the way ESLint validates JavaScript — it enforces the standard, it's not the standard itself.

**Key differentiators:**

| vs. No formal specs | vs. SpecFlow / OpenSpec |
|---|---|
| "Your AI agents are guessing what you meant" | "SpecScore validates specs, not just defines them" |
| "Specs as plain text = version-controlled, diffable, reviewable" | "Linter + LSP = catch ambiguity before execution" |
| "Adopt incrementally — one spec at a time" | "Plain Markdown/YAML, not a proprietary format" |

**Ecosystem hint (visible but not primary):**

> SpecScore defines *what* gets built. When you're ready to test specs automatically (Rehearse) or orchestrate agents across them (Synchestra), the ecosystem is there. But SpecScore works beautifully on its own.

**Owned term:** "spec-driven development" — SpecScore should own this phrase the way Obsidian owns "second brain" and Tailwind owns "utility-first CSS."

## Section 2: Website Structure

**Philosophy:** The site is a specification reference first, a product page second. The format docs *are* the product.

### Page Hierarchy

```
specscore.md/
├── Home (hero + value prop + quick start + ecosystem hint)
├── Spec/
│   ├── Overview (what is the SpecScore format)
│   ├── Feature
│   ├── Requirement
│   ├── Acceptance Criteria
│   ├── Scenario
│   ├── Plan
│   ├── Task
│   ├── Project Definition
│   └── Source References
├── Guides/
│   ├── Getting Started (write your first spec)
│   ├── For Developers
│   ├── For Product Owners
│   ├── For QAs
│   ├── For Business Analysts
│   ├── For Project Managers
│   └── For Architects
├── Tools/
│   ├── CLI (lint, validate, traverse)
│   ├── LSP (editor integration)
│   └── Integrations (Claude, Cursor, Copilot, etc.)
├── Blog/
├── GitHub (→ repo)
└── Ecosystem (→ Synchestra, Rehearse, Impresario — light touch)
```

### Home Page Structure

1. **Hero** — Tagline + one-sentence pitch + "Read the spec" primary CTA + "Install the CLI" secondary CTA
2. **The problem** — Technical-but-warm voice. "AI agents work from specifications. Most specifications are ambiguous, scattered, and untestable."
3. **The format** — Show a real SpecScore spec (feature with requirements, acceptance criteria). Let the format sell itself.
4. **The tooling** — CLI linting in action + LSP catching issues in-editor. "Write specs. Validate them. Ship them."
5. **Role-based entry points** — Six cards linking to persona guides.
6. **Ecosystem hint** — One row. "SpecScore defines what. Rehearse validates. Synchestra orchestrates."
7. **Footer** — GitHub, contributing, license (open source emphasis)

### Key Design Decisions

- "Read the spec" is the primary CTA, not "Install the CLI" — standard-first positioning
- Role-based guides are prominent, not buried — multi-stakeholder tool
- Blog lives on specscore.md (not synchestra.io) — SpecScore owns "spec-driven development" content
- Ecosystem page is a single page, not a section — keeps SpecScore's identity clean

## Section 3: Content Strategy

### Content Pillars

| Pillar | Purpose | Examples |
|---|---|---|
| **The Methodology** | Establish spec-driven development as a practice | "What is spec-driven development?", "Why AI agents need formal specs", "Specs vs. tickets" |
| **The Format** | Teach the SpecScore format | "Writing your first SpecScore feature", "Acceptance criteria that machines can test", "Source references: linking code back to specs" |
| **The Tooling** | Show the CLI + LSP in action | "Linting specs like you lint code", "SpecScore in VS Code", "CI/CD for specifications" |
| **The Comparison** | Position against alternatives | "SpecScore vs. no specs", "SpecScore vs. Gherkin", "SpecScore vs. SpecFlow" |
| **The Community** | Build credibility and adoption | Case studies, contributor spotlights, "how we use SpecScore at [team]" |

### Content Cadence

- **Launch:** 3-4 pillar articles ready on day one (methodology intro, getting started, one comparison, one tooling walkthrough)
- **Post-launch:** 1-2 posts per month, alternating pillars
- **Evergreen:** Spec reference docs and role-based guides double as SEO content

### SEO Plays

- "spec-driven development" — own this term before anyone else does
- "AI agent specifications" / "specifications for AI agents" — high intent, low competition
- "how to write specs for [Claude/Cursor/Copilot]" — integration-specific long-tail
- "SpecScore vs [alternative]" — comparison pages

### Distribution Channels

- **GitHub README** — most OSS adoption starts here
- **Hacker News** — launch post + follow-up "Show HN" when CLI is polished
- **Reddit** — r/programming, r/artificial, r/ChatGPTCoding, r/LocalLLaMA
- **Dev.to / Hashnode** — cross-post pillar articles
- **Twitter/X** — short-form spec tips, CLI demos, threading the methodology
- **Obsidian community** — natural overlap if the Obsidian companion angle materializes

### What We Deliberately Don't Do Yet

- No Product Hunt launch until CLI + site are polished
- No paid ads — too early, audience too niche
- No newsletter — wait until there's enough content to sustain it

## Section 4: Launch Strategy

**Launch philosophy:** Quiet, credible, developer-first. The standard should feel like it already exists and you're just now finding out about it.

### Phase 0: Foundation (before any public announcement)

- specscore.md live with spec docs, getting started guide, and role-based guides
- GitHub repo clean, README polished, CONTRIBUTING.md in place
- CLI installable via `go install` or Homebrew
- 3-4 pillar blog posts published
- At least one real-world project using SpecScore publicly (own repos count)

### Phase 1: Soft Launch (week 1-2)

- **"Show HN: SpecScore — an open spec format for AI-driven development"**
  - Lead with the problem, show the format, link to spec docs (not the CLI)
  - Be present in comments for 24h
- **Reddit posts** — different angle per subreddit:
  - r/programming: "We built a linter for specifications"
  - r/artificial: "A spec format AI agents can actually parse and validate"
- **Twitter/X thread:** "Why we created a specification standard for AI agents" — the story, not the pitch
- **Goal:** 50-100 GitHub stars, initial feedback, early adopters who file issues

### Phase 2: Content Flywheel (month 1-3)

- 1-2 posts per month from the content pillars
- "SpecScore vs. X" comparison pages go live
- Engage in relevant threads — anywhere someone asks "how do I write specs for AI agents," SpecScore shows up naturally
- Invite early adopters to contribute
- **Goal:** 500+ stars, community contributions, "spec-driven development" starts appearing in search

### Phase 3: Expansion (month 3-6)

- Product Hunt launch
- Conference talks / meetup presentations
- Obsidian plugin or companion (if the angle proves out)
- Rehearse and Synchestra start appearing in ecosystem content — the funnel opens
- **Goal:** 1000+ stars, spec-driven development recognized as a methodology, early enterprise interest

### Success Criteria at 6 Months

- "spec-driven development" returns SpecScore as the top result
- Developers are writing SpecScore specs in projects outside the ecosystem
- PMs/BAs are finding SpecScore through role-based guides
- The ecosystem hint is converting visitors into Synchestra/Rehearse awareness

## Updates Needed in Existing Docs

These changes live in the `synchestra-marketing` repo (not this repo) and should be made when specscore.md goes live:

- `synchestra-marketing/STRATEGY.md` — update all `specscore.org` references to `specscore.md`
- `synchestra-marketing/EXTRACTION_PLAN.md` — update all `specscore.org` references to `specscore.md`
- `synchestra-marketing/ECOSYSTEM.md` — no domain references, but verify messaging alignment with this plan

## Outstanding Questions

None at this time.
