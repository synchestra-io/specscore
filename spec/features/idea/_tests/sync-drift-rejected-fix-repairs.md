# Scenario: Strict lint fails on sync drift; `--fix` repairs it

**Validates:** [idea#ac:sync-strictness](../README.md#ac-sync-strictness), [idea#req:sync-lint-strict](../README.md#req-sync-lint-strict)

## Steps

GIVEN a feature at `spec/features/offline-sync/README.md` with `**Source Ideas:** offline-mode`
AND an idea at `spec/ideas/offline-mode.md` with `**Status:** Approved` and `**Promotes To:** —`
WHEN `specscore lint` runs without `--fix`
THEN the linter reports an error: idea "offline-mode" drift — `Promotes To` and `Status` disagree with referencing features
AND re-running `specscore lint --fix` rewrites the idea header to `**Status:** Specified` and `**Promotes To:** offline-sync`
AND a subsequent `specscore lint` run succeeds

---
*This document follows the https://specscore.md/scenario-specification*
