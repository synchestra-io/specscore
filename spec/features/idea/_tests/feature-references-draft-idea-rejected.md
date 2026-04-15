# Scenario: Feature referencing a `Draft` idea is rejected

**Validates:** [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:feature-cross-reference](../README.md#req-feature-cross-reference)

## Steps

GIVEN an idea at `spec/ideas/offline-mode.md` with `**Status:** Draft`
AND a feature at `spec/features/offline-sync/README.md` with `**Source Ideas:** offline-mode`
WHEN the spec linter validates the spec tree
THEN the linter reports an error: feature "offline-sync" references idea "offline-mode" whose status is `Draft`
AND the error message notes that referenced ideas must have `Status ∈ {Approved, Specified}`
