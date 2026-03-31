# Scenario: Immutability after approval

**Validates:** [development-plan#req:immutability-after-approval](../README.md#req-immutability-after-approval)

## Steps

GIVEN a plan with status `approved` and content hash `abc123`
WHEN an edit is attempted to the plan's Steps section
THEN the edit is rejected with an error indicating approved plans are immutable

GIVEN a plan with status `approved` that needs a different approach
WHEN a new plan is created with a `supersedes` reference to the original
AND the original plan's status is changed to `superseded`
THEN both plans exist: the original as a historical record and the new plan as the active approach
