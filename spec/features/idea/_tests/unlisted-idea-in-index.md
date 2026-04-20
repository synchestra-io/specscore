# Scenario: Idea missing from index is a validation error

**Validates:** [idea#req:index-completeness](../README.md#req-index-completeness)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Status:** Approved`
AND the idea index at `spec/ideas/README.md` does not list `offline-mode`
WHEN the spec linter validates the idea tree
THEN the linter reports an error: idea "offline-mode" exists but is not listed in the active idea index
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
