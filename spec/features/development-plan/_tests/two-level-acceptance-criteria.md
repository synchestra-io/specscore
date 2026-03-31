# Scenario: Two-level acceptance criteria

**Validates:** [development-plan#req:two-level-acceptance-criteria](../README.md#req-two-level-acceptance-criteria)

## Steps

GIVEN a plan with a `## Acceptance criteria` section containing cross-cutting criteria
AND Step 1 has its own `**Acceptance criteria:**` bullet list
AND Step 2 has its own `**Acceptance criteria:**` bullet list
WHEN the plan is parsed by an execution tool
THEN the plan-level criteria are identified as integration/end-to-end targets
AND the step-level criteria are identified as unit/component targets
AND both levels are available for populating work item descriptions
