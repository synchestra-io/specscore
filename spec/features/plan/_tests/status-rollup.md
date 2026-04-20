# Scenario: Status rollup

**Validates:** [plan#req:status-rollup](../README.md#req-status-rollup)

## Steps

GIVEN a plan with three child tasks all in `draft` status
WHEN the plan's derived status is computed
THEN the derived status is `draft`

GIVEN a plan with two children: one `in_review` and one `draft`
WHEN the plan's derived status is computed
THEN the derived status is `in_review`

GIVEN a plan with all children in `approved` status
WHEN the plan's derived status is computed
THEN the derived status is `approved`

GIVEN a plan with a derived status of `draft` but an explicitly set status of `in_review`
WHEN the plan status is evaluated
THEN the explicit status `in_review` takes precedence over the derived value
AND tooling surfaces a discrepancy between the explicit and derived statuses

---
*This document follows the https://specscore.md/scenario-specification*
