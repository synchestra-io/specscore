---
title: What Is Spec-Driven Development?
description: Why AI agents need formal specifications — and how SpecScore makes it practical.
date: 2026-04-14
---

# What Is Spec-Driven Development?

Running one AI agent on one task works fine. Tell it what to build, it builds something. Maybe it's right. Maybe you spend an hour fixing what it misunderstood.

Running five agents across three platforms on a tree of interdependent tasks? That's where things collapse. Not because the agents are bad — because the specifications are.

## The Problem With Informal Specs

Most teams specify work in one of three ways:

1. **Jira tickets** — a title, a description, maybe some acceptance criteria written by someone in a hurry
2. **Notion docs** — rich, detailed, and never read by the agent that actually does the work
3. **Conversation** — "just make it work like the mockup"

None of these are machine-readable. None of them are validatable. None of them link back to the code that implements them.

When an AI agent works from an ambiguous spec, it guesses. Sometimes it guesses right. When it doesn't, you get rework — and rework with AI agents is expensive because it compounds across every downstream task that depended on the wrong assumption.

## Spec-Driven Development

Spec-driven development is a practice where specifications are:

- **Structured** — features, requirements, acceptance criteria, plans, and tasks follow a consistent format
- **Machine-readable** — AI agents can parse the spec and understand what to build without human translation
- **Validatable** — a linter catches ambiguity, missing fields, and structural problems before anyone starts building
- **Traceable** — source code links back to the spec it implements via inline annotations

The specification is the source of truth. Not the ticket. Not the conversation. Not the PR description.

## What This Looks Like in Practice

A SpecScore feature specification is Markdown with structure:

```yaml
title: User Authentication
status: draft
priority: high
```

Inside the feature directory, you define requirements with acceptance criteria:

```markdown
## Requirements

### R1: Email/password login

**Acceptance Criteria:**
- AC1: Given valid credentials, the system returns a session token
- AC2: Given invalid credentials, the system returns a 401 with a descriptive error
- AC3: Passwords are never logged or returned in API responses
```

Then you run `specscore lint` and it tells you what's missing, what's ambiguous, and what doesn't conform to the schema.

The developer (or agent) implementing this feature adds a source reference:

```go
// specscore:user-authentication/R1
func Login(email, password string) (*Session, error) {
```

Now you can trace from spec to code and from code to spec.

## Why Now?

Three things changed:

1. **AI agents are doing real work** — not just autocomplete, but multi-step implementation across codebases
2. **Specifications became the bottleneck** — agent capability outpaced spec quality
3. **Markdown won** — it's the lingua franca of developer tools, and now BAs and PMs write it too (thanks to AI assistants)

Spec-driven development is what test-driven development was for code quality — but for the layer above code. The spec is the contract. Everything flows from it.

## Getting Started

SpecScore is the open specification framework for spec-driven development:

- [Read the spec](/specifications) — understand the format
- [Install the CLI](/install) — start validating your specifications
- [Find your role](/for/developers) — guides for developers, PMs, QAs, BAs, architects, and project managers

The format is Markdown and YAML. The schema is published and versioned. The CLI is open source. No vendor lock-in.
