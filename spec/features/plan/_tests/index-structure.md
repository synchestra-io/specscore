# Scenario: Index structure

**Validates:** [plan#req:index-completeness](../README.md#req-index-completeness), [plan#req:index-child-indentation](../README.md#req-index-child-indentation), [plan#req:index-roi-columns](../README.md#req-index-roi-columns), [plan#req:index-recently-closed](../README.md#req-index-recently-closed)

## Steps

GIVEN three plans exist: `user-auth`, `add-batch-mode`, and `refactor-output`
AND the plans index lists only `user-auth` and `add-batch-mode`
WHEN the index is validated
THEN validation rejects with an error indicating `refactor-output` is not listed

GIVEN a plan `chat-feature` with sub-plan `chat-infrastructure`
WHEN the plans index is rendered
THEN the sub-plan row is indented with `&ensp;`
AND its link path includes the parent directory: `chat-feature/chat-infrastructure/`

GIVEN a plan with Effort `M` and Impact `high`
AND another plan with no Effort or Impact metadata
WHEN the plans index is rendered
THEN the first plan shows `M` and `high` in the respective columns
AND the second plan shows `-` in both columns

GIVEN a plan that was completed within the last 5 entries
WHEN the plans index is rendered
THEN the plan appears in the Recently Closed section

---
*This document follows the https://specscore.md/scenario-specification*
