# Scenario: Nesting limit enforcement

**Validates:** [development-plan#req:nesting-limit](../README.md#req-nesting-limit)

## Steps

GIVEN a plan with steps at level 1 (e.g., `### 1. Define schema`)
AND sub-steps at level 2 (e.g., `#### 2.1. Add streaming support`)
WHEN the plan is validated
THEN validation passes

GIVEN a plan with a sub-sub-step at level 3 (e.g., `##### 2.1.1. Handle edge cases`)
WHEN the plan is validated
THEN validation rejects the plan with an error indicating steps must not exceed two levels of nesting
