# Scenario: Key Assumptions table without a Must-be-true row is rejected

**Validates:** [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:must-be-true-present](../README.md#req-must-be-true-present)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with all required sections present
AND the `## Key Assumptions to Validate` table contains only Should-be-true and Might-be-true rows
WHEN the spec linter validates the idea
THEN the linter reports an error: `Key Assumptions to Validate` must list at least one Must-be-true assumption
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
