# Scenario: Status transition rules

**Validates:** [development-plan#req:status-transitions](../README.md#req-status-transitions)

## Steps

GIVEN a plan with status `draft`
WHEN the status is changed to `in_review`
THEN the transition is accepted

GIVEN a plan with status `in_review`
WHEN the status is changed to `draft` (revisions requested)
THEN the transition is accepted

GIVEN a plan with status `in_review`
WHEN the status is changed to `approved`
THEN the transition is accepted

GIVEN a plan with status `approved`
WHEN the status is changed to `superseded`
THEN the transition is accepted

GIVEN a plan with status `draft`
WHEN the status is changed directly to `approved` (skipping `in_review`)
THEN the transition is rejected with an error indicating the allowed transitions from `draft`

GIVEN a plan with status `approved`
WHEN the status is changed to `draft`
THEN the transition is rejected with an error indicating approved plans can only transition to `superseded`
