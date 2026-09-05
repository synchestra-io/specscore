---
format: https://specscore.md/idea-specification
status: Specified
---

# Idea: Spec-code traceability via CodeGrapher

**Status:** Specified
**Date:** 2026-06-10
**Owner:** trakhimenok
**Promotes To:** source-references
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we link SpecScore requirements and acceptance criteria to the code symbols that implement them and the executable tests that verify them, then keep those links current across edits, refactors, worktrees, and exact Git revisions?

## Context

SpecScore already owns stable Feature, REQ, AC, scenario, and source-reference semantics. Its CLI scans source comments itself. CodeGrapher already parses SpecScore documents and source code, but it cannot currently address normalized REQ and AC records or express typed implementation and verification edges. This duplicates scanning while leaving the most valuable traceability relationship unavailable.

CodeGrapher imports SpecScore parser packages. Importing CodeGrapher back into SpecScore would create a Go module dependency cycle and tightly couple SpecScore to one index implementation. The tools need a clear ownership boundary and a versioned process contract.

## Recommended Direction

1. **SpecScore owns meaning.** Export normalized artifact, requirement, acceptance-criterion, and scenario records with stable IDs and source ranges. Define typed source relations such as `implements`, `verifies`, and `references` while preserving canonical clickable SpecScore URLs.
2. **CodeGrapher owns discovery and graph storage.** Scan source and specification files, attach annotations to the nearest code or test symbol, and store typed edges between symbols, tests, REQs, ACs, scenarios, and features.
3. **SpecScore queries CodeGrapher out of process.** A versioned JSON CLI or service contract supplies trace, suggestion, validation, and coverage results. SpecScore does not import the CodeGrapher module or read its SQLite files.
4. **WB coordinates incremental work.** It supplies changed paths, schedules low-priority graph refreshes, binds evidence to the exact target and worktree revisions, and reports affected specifications and dependent repositories.
5. **Agents propose; committed links decide.** An agent may identify candidate code/REQ and test/AC pairs from existing comments, test names, scenario references, call graphs, coverage, and package proximity. Each suggestion carries evidence and confidence. It does not count as accepted implementation or verification evidence until a human or implementation agent commits an explicit annotation.
6. **SpecScore Studio renders backlinks from the graph.** Feature, requirement, and AC pages show current implementing symbols, verifying tests, coverage evidence, and drift at an exact revision. Specifications do not store generated line-number backlink tables.

### Coverage attribution

WB SHOULD collect the language-native coverage artifact once at the selected validation gate. CodeGrapher maps covered source ranges to code symbols, then follows accepted typed links to attribute coverage to REQs, ACs, and features. SpecScore presents at least two separate signals:

- **implementation exercised** — linked implementation code was executed by the measured suite;
- **acceptance criterion verified** — an executable test explicitly linked with `verifies` passed.

Coverage of implementation code MUST NOT by itself prove an AC. Every coverage result records the exact commit, command, package or project scope, coverage mode, and collection time. Partial or stale evidence is labelled rather than silently combined. Expensive per-test coverage is optional; aggregate suite coverage is the default.

## Alternatives Considered

- **Keep scanning in SpecScore.** This preserves the current dependency shape but duplicates parsers, cannot use the symbol graph efficiently, and makes worktree overlays and cross-repository impact harder.
- **Import CodeGrapher into SpecScore.** This offers an in-process API but creates a module dependency cycle and prevents CodeGrapher from remaining a replaceable provider.
- **Infer all links without annotations.** This is convenient initially but turns probabilistic guesses into apparent evidence. Suggestions are valuable; accepted traceability must remain explicit and reviewable.

## MVP Scope

1. Export normalized SpecScore Feature, REQ, AC, and scenario records with source ranges.
2. Define and parse typed code annotations for implementation, verification, and general reference relationships.
3. Extend CodeGrapher extraction and JSON queries for typed links, missing links, dangling links, and candidate pairs.
4. Update the SpecScore implementation skill so implementation symbols link to REQs and executable tests link to ACs or directly verified REQs.
5. Let WB refresh changed paths, validate accepted links, and attach exact-revision coverage evidence without blocking edit hooks.
6. Measure suggestion precision, incremental refresh latency, coverage attribution completeness, and agent tool-call savings on SpecScore CLI and WB.

## Not Doing (and Why)

- Treating an agent suggestion as accepted proof — inferred links remain proposals.
- Treating line or branch coverage as proof that an AC passed — explicit verifying-test evidence remains distinct.
- Enforcing complete link coverage immediately — completeness begins as a report until adoption data supports a gate.
- Reading CodeGrapher database files from SpecScore — integration uses a versioned provider contract.
- Initializing, uploading, or fully rebuilding a graph from an edit hook — edit refresh is incremental, debounced, and non-blocking.
- Requiring per-test coverage — it can multiply test cost and is only collected when its additional attribution value is requested.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Explicit typed annotations remain readable and cheap to maintain | Apply them to one Go CLI and review refactor behavior |
| Must-be-true | CodeGrapher can attach comments and coverage ranges to stable symbols accurately | Compare extracted edges with a hand-reviewed fixture corpus |
| Should-be-true | Changed-path refresh plus a worktree overlay is substantially cheaper than full indexing | Benchmark edits in WB and SpecScore CLI worktrees |
| Should-be-true | Aggregate coverage attribution is useful without per-test coverage | Compare feature and AC reports with linked-test evidence |
| Might-be-true | Agent suggestions remove most manual pairing work | Measure accepted, rejected, and corrected suggestions |

## SpecScore Integration

- **New Features this would create:** Code graph provider and traceability queries may become separate Features after the spike
- **Existing Features affected:** source-references, requirement, acceptance-criteria, scenario
- **Dependencies:** CodeGrapher JSON/service contract; WB validation and evidence receipts

## Open Questions

- Should canonical typed annotations use a keyword before the URL, a URL query parameter, or another grammar that remains clickable and language-agnostic?
- Should fragments standardize on lowercase `#req:` and `#ac:` while accepting existing uppercase `#REQ:` references during migration?
- How should mappings be recorded for generated source that cannot carry comments?
- Which coverage formats must the first provider contract support beyond Go cover profiles?
- When is per-test coverage worth its additional execution cost?
