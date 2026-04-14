# The Ecosystem

SpecScore is the specification layer in a broader ecosystem of tools for AI-driven development. Each tool works independently but is designed to complement the others.

## SpecScore — The Score

Defines *what* gets built. An open specification format with validation tooling.

- **Format:** Markdown and YAML — version-controlled, portable, no vendor lock-in
- **Tooling:** CLI linter, LSP for editor integration
- **Standalone value:** Write better specs, validate them automatically, link them to code

## Rehearse — The Rehearsal

Tests specifications automatically *before* implementation begins.

- Validates that specs are complete, consistent, and testable
- Catches ambiguity and gaps before agents start building
- Works with any project that uses the SpecScore format

[rehearse.ink](https://rehearse.ink)

## Synchestra — The Performance

Orchestrates multi-agent work *across* specifications.

- Coordinates AI agents claiming and executing tasks from SpecScore specs
- Schema-validated state at every commit — git is the protocol
- Scales from a single developer to distributed teams

[synchestra.io](https://synchestra.io)

## How They Fit Together

```
Write specs with SpecScore
         ↓
Validate with Rehearse
         ↓
Orchestrate with Synchestra
         ↓
Agents coordinate using SpecScore as the shared protocol
```

## Each Tool Standalone, Better Together

You don't need the full ecosystem. SpecScore works with any orchestration tool — Jira, Linear, your own scripts. Rehearse works in any project that uses SpecScore-formatted specs. Synchestra is optimized for SpecScore but is not required.

The recommended path: start with SpecScore. Add Rehearse when you want automated validation. Add Synchestra when you need multi-agent coordination.
