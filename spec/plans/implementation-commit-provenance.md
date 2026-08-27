---
format: https://specscore.md/plan-specification
status: Implemented
---
# Plan: Implementation Commit Provenance

**Status:** Implemented
**Source Feature:** implementation-commit-provenance
**Date:** 2026-06-26
**Owner:** alex
**Supersedes:** —

## Summary

Master plan for the three-repo `implementation-commit-provenance` effort. This plan owns the **data-model and enforcement** layer — the optional `implementation_commit` Task-entity property, the syntactic provenance-ref-format lint rule, the plan-level evidence rollup, and the Snapshots boundary. The actual lint/rollup code lands in `specscore-cli`; this plan tracks those cross-cutting data-model ACs. Two child sub-plans carry the rest of the work and declare this plan as their `**Parent:**`.

## Approach

Linear, bottom-up. **Task 1** pins the data shape (the entity property — already added to `task.entity.md`). **Task 2** implements the syntactic format validation, and **Task 3** the rollup + Snapshots-distinct surfacing; both are realized by `specscore-cli` lint code (delegated to the CLI sub-plan, which also carries the verb itself). Implementation order across the composition: this plan's data model → the CLI verb+lint (sub-plan `cli-task-change-status` in `specscore-cli`) → the agent skill (sub-plan `lifecycle-status-skill-task-kind` in `ai-plugin-specscore`). The two sub-plans are wired to this master via their cross-repo `**Parent:** specscore:implementation-commit-provenance` reference, which is the canonical, lint-validated composition link (`P-005`).

## Tasks

### Task 1: Define the optional `implementation_commit` Task-entity property

**Verifies:** implementation-commit-provenance#ac:provenance-is-optional, implementation-commit-provenance#ac:provenance-never-inferred
**Depends-On:** —
**Status:** complete

Define the optional `implementation_commit` property on the Task entity (`spec/features/task/task.entity.md`) and its document surface `**Implemented-by:**` — the typed data shape, optional and actor-supplied (never inferred from ambient `HEAD`). The property is already added; this task confirms the shape and the never-inferred guarantee that downstream code relies on.

### Task 2: Provenance-ref-format syntactic lint rule (in specscore-cli)

**Verifies:** implementation-commit-provenance#ac:well-formed-ref-accepted, implementation-commit-provenance#ac:malformed-ref-rejected, implementation-commit-provenance#ac:same-repo-bare-sha
**Depends-On:** 1
**Status:** complete

Implement the syntactic validator for the `<repo>@<sha> (<branch>)` reference (bare `<sha>` for same-repo), accepting on shape alone without resolving the repo. The code lands in `specscore-cli` (sub-plan `cli-task-change-status`); this task tracks the cross-cutting data-model ACs it satisfies.

### Task 3: Plan-evidence rollup and Snapshots boundary (in specscore-cli)

**Verifies:** implementation-commit-provenance#ac:plan-evidence-rolls-up, implementation-commit-provenance#ac:snapshots-stay-distinct
**Depends-On:** 2
**Status:** complete

Implement the derived plan-level rollup of child tasks' `implementation_commit` values and surface it **query-only via `specscore plan info`** (no plan-body/frontmatter write), as evidence distinct from the plan `## Snapshots` Git Hash. Code lands in `specscore-cli` (`cli/plan/info` output); this task tracks the data-model ACs.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/plan-specification*
