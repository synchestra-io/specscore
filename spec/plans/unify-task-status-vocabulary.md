---
format: https://specscore.md/plan-specification
status: Implemented
---
# Plan: Unify Task Status Vocabulary

**Status:** Implemented
**Source Feature:** unify-task-status-vocabulary
**Date:** 2026-06-26
**Owner:** alex
**Supersedes:** —

## Summary

Makes the canonical 7-value Task status enum the single task-status vocabulary in practice. Work lands in `specscore-cli` (the lint vocabulary rule + `--fix` migration, the `plan new` scaffold, the rollup confirmation) and the specstudio `implement` skill (the canonical progress lifecycle). This plan's own task blocks carry `**Status:** done` — the only value the *current* lint rule `P-004` accepts; `P-004` rejects the canonical `planning` until Task 1 lands. Once Task 1 ships, `lint --fix` migrates these to `planning`. (That bootstrap irony is the bug in miniature.)

## Approach

Linear, enforcement-first. **Task 1** builds the lint vocabulary rule (flag any non-enum/legacy token) plus the `--fix` value-for-value migration — the foundation everything else relies on. **Task 2** fixes the `plan new` scaffold to emit `planning`. **Task 3** confirms the execution-band rollup is canonical and that a migrated all-`complete` plan reaches `Implemented`; it depends on the migration (Task 1). **Task 4** aligns the specstudio `implement` skill's status writes to the `planning→queued→in_progress→complete` lifecycle; it depends on the vocabulary being enforced (Task 1) so the skill writes values lint accepts. Cross-repo: Tasks 1–3 are `specscore-cli`; Task 4 is the specstudio `implement` skill.

## Tasks

### Task 1: Lint vocabulary enforcement + legacy-token migration

**Verifies:** unify-task-status-vocabulary#ac:enum-is-sole-vocabulary, unify-task-status-vocabulary#ac:lint-flags-legacy, unify-task-status-vocabulary#ac:lint-fix-migrates-legacy
**Depends-On:** —
**Status:** complete

Change lint rule **P-004** — whose current accepted set is the hybrid `{pending, in-progress, done, blocked, failed, aborted}` — to validate plan-inline task `**Status:**` against the canonical 7-value enum (`planning, queued, in_progress, blocked, complete, failed, aborted`). Plain `lint`: legacy tokens (`pending`/`done`/`in-progress`) are flagged naming their canonical replacement; any other non-enum token (e.g. `shipped`) is rejected with no mapping. Under `--fix`: rewrite the legacy tokens value-for-value (`pending`→`planning`, `done`→`complete`, `in-progress`→`in_progress`), changing nothing else.

### Task 2: Scaffold emits `planning`

**Verifies:** unify-task-status-vocabulary#ac:scaffold-emits-planning
**Depends-On:** 1
**Status:** complete

Change the `specscore plan new` scaffold to emit `**Status:** planning` on each generated task block instead of `pending`, so newly created plans are canonical by construction (requires Task 1's `P-004` change to accept `planning`).

### Task 3: Confirm rollup canonical, migrated plan reaches Implemented

**Verifies:** unify-task-status-vocabulary#ac:rollup-reaches-implemented
**Depends-On:** 1
**Status:** complete

Confirm (and test) that the execution-band rollup reads canonical values and that a plan whose tasks were legacy `done`, after `lint --fix` migration to `complete`, rolls up to `Implemented` — the outcome the legacy token silently never produced. No rollup logic change is expected; this verifies the migration closes the loop end-to-end.

### Task 4: Align the specstudio implement skill lifecycle

**Verifies:** unify-task-status-vocabulary#ac:implement-skill-lifecycle
**Depends-On:** 1
**Status:** complete

In the specstudio `implement` skill, replace the `{pending, in-progress, done, blocked}` status writes with the canonical lifecycle `planning → queued → in_progress → complete` (plus `blocked`/`failed`/`aborted`), updating the skill's status-protocol table, per-task status-write transitions, and self-review token set so it never writes a legacy token. This is the cross-repo consumer change that realizes the lifecycle AC rather than leaving it a dangling contract.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/plan-specification*
