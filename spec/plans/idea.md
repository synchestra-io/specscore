---
format: https://specscore.md/plan-specification
status: Approved
---

# Plan: Idea Types and Lifecycle

**Status:** Approved
**Source Feature:** idea
**Date:** 2026-05-25
**Owner:** alex
**Supersedes:** —

## Summary

Implements the idea Feature's type system (feature-request / change-request), extended lifecycle (Specifying, Implemented), change-request location at `{feature}/proposals/`, and CLI tooling — across specscore-cli (lint rules, sync logic, scaffold commands) and the specscore spec repo (forward reference updates).

## Approach

Grouped by subsystem: lint validation first (the foundation everything else depends on), then lifecycle sync logic (depends on lint recognizing the new statuses and types), then archival/cross-reference behavior (depends on both), then scaffold/CLI surface (depends on lint accepting the output), then regression verification. Two ACs (`idea-structure`, `authoring-independence`) are unchanged from the current implementation but covered by a regression-verification task to confirm the new code paths don't break existing behavior. Forward reference updates in `feature/README.md` and `plan/README.md` are tracked separately as change requests against those Features.

## Tasks

### Task 1: Lint rule updates for idea types and headers

**Verifies:** idea#ac:idea-header, idea#ac:idea-types

Implement lint validation for the new idea type system in specscore-cli:

- Accept `# Proposal: <Title>` as a valid title prefix that dispatches to idea validation rules, requiring `Type: change-request` and a valid `Targets:` field. Reject `# Proposal:` without these fields, and `# Idea:` with `Type: change-request`.
- Validate `**Type:**` field values (`feature-request` | `change-request`); treat absent Type as `feature-request`.
- Validate `**Targets:**` field: required and non-empty when Type is `change-request`; must reference an existing Feature directory. Prohibited (or `—`) when Type is `feature-request`.
- Validate `**Phase:**` field: optional free-form non-empty string when present.
- Validate change-request idea location: must reside at `spec/features/<feature-slug>/proposals/<slug>.md` where `<feature-slug>` matches `Targets:` value.
- Validate header field ordering: `Type:` and `Targets:` after `Status:`, before `Date:`; `Phase:` after `Targets:`.

### Task 2: Lifecycle derivation and sync updates

**Verifies:** idea#ac:promotion-lifecycle, idea#ac:change-request-lifecycle, idea#ac:sync-strictness

Extend the idea status system and sync logic in specscore-cli:

- Add `Specifying`, `Specified` (redefined), `Implementing` (adjusted), and `Implemented` to the valid status set.
- Implement derivation rules for feature-request ideas: `Specifying` when referenced Feature at Draft/Under Review; `Specified` when at Approved; `Implementing` when at Implementing; `Implemented` when at Stable.
- Bypass derivation for change-request ideas: accept author-managed `Specifying`/`Specified`/`Implementing`/`Implemented` without checking Feature references.
- Enforce `derived-status-not-author-set` for feature-request ideas only.
- Update `specscore idea sync` / `specscore lint --fix` to recompute statuses using the new derivation rules.
- Update `sync-lint-strict` to cover the new statuses and to skip change-request ideas.
- Update `implementing-requires-promotion` to apply only to feature-request ideas.
- Update `feature-cross-reference` to accept Ideas at `Specifying`, `Specified`, `Implementing`, and `Implemented` in addition to `Approved`.

### Task 3: Archival and cross-reference extensions

**Verifies:** idea#ac:archival, idea#ac:related-ideas

Extend archival behavior and cross-reference resolution in specscore-cli:

- Implement `--include-archived` flag on `specscore idea list`; exclude Archived ideas from default output; all other statuses (including `Implemented`) remain visible by default.
- For change-request ideas, archival leaves the file at `spec/features/<feature>/proposals/<slug>.md` (no move to `spec/ideas/archived/`). Validate that archived change-request ideas are NOT listed in the archived index.
- Extend `related-ideas-target-exists` resolution to scan `spec/features/*/proposals/` in addition to `spec/ideas/` and `spec/ideas/archived/`.
- Extend `index-completeness` to scan `spec/features/*/proposals/` for change-request ideas and require them in the active index.

Note: Forward reference updates in `feature/README.md` and `plan/README.md` (changing the proposal link targets and adding the `## Proposals` index section convention) are tracked separately as change requests against the `feature` and `plan` Features — they are outside the `idea` Feature's AC scope.

### Task 4: Scaffold command and CLI alias

**Verifies:** idea#ac:scaffold-behavior

Extend the scaffold command and add the proposal alias in specscore-cli:

- Add `--type` and `--targets` flags to `specscore idea new`. When `--type change-request --targets <feature-slug>` is supplied, scaffold at `spec/features/<feature-slug>/proposals/<slug>.md` with `# Proposal:` prefix and `Type:`/`Targets:` fields pre-populated. Auto-create the `proposals/` directory if absent.
- Implement `specscore proposal new <feature-slug> <slug>` as a convenience alias that delegates to `specscore idea new <slug> --type change-request --targets <feature-slug>`.
- Ensure scaffolded output passes `specscore lint` (existing `always-lint-clean-on-exit` contract).
- Support `--phase` flag for optional Phase field pre-population.

### Task 5: Regression verification for unchanged ACs

**Verifies:** idea#ac:idea-structure, idea#ac:authoring-independence

Verify that existing validation for required-sections, not-doing-non-empty, must-be-true-present, and authoring-agnostic rules still passes after the type/lifecycle changes. Run the existing test suite for these rules against both feature-request and change-request idea fixtures to confirm no regressions. No new logic expected — this task validates that the new code paths do not break existing behavior.

### Task 6: Update CLI specifications in specscore-cli

**Verifies:** idea#ac:idea-header, idea#ac:idea-types, idea#ac:promotion-lifecycle, idea#ac:change-request-lifecycle, idea#ac:archival, idea#ac:scaffold-behavior

Update the CLI feature specs in `specscore-cli/spec/features/cli/` to reflect all new lint rules, lifecycle changes, and CLI commands added in Tasks 1–4. Per AGENTS.md convention ("When changing rule behavior, update the spec in the same commit"), this task ensures the CLI spec stays in sync with the implementation:

- Update `spec/features/cli/idea/README.md` — document `idea new --type/--targets` flags, `proposal new` alias, and the extended lifecycle.
- Update `spec/features/cli/spec/lint/README.md` (or create an `idea-rules/` child) — document the 6 new lint rules added in Task 1 (type-values, type-title-consistency, targets-required, targets-exists, change-request-location, phase-non-empty) and the modified rules from Task 2 (status-values, sync-lint-strict, specified-requires-promotion, feature-cross-reference).
- Update `spec/features/cli/lifecycle-transitions/README.md` — document the new Specifying and Implemented statuses and the extended transition matrix.

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/plan-specification*
