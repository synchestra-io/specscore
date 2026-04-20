# Scenario: Empty `Not Doing` section is rejected

**Validates:** [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:not-doing-non-empty](../README.md#req-not-doing-non-empty)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with all required sections present
AND the `## Not Doing (and Why)` section contains no bullet entries
WHEN the spec linter validates the idea
THEN the linter reports an error: `Not Doing (and Why)` must contain at least one explicit exclusion with a reason
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
