# Scenario: Valid status values

**Validates:** [development-plan#req:valid-statuses](../README.md#req-valid-statuses), [development-plan#req:no-execution-status](../README.md#req-no-execution-status)

## Steps

GIVEN a plan document with status `draft`
WHEN the status is validated
THEN validation passes

GIVEN a plan document with status `approved`
WHEN the status is validated
THEN validation passes

GIVEN a plan document with status `completed`
WHEN the status is validated
THEN validation rejects the status with an error indicating `completed` is an execution concern and not a valid plan status

GIVEN a plan document with status `failed`
WHEN the status is validated
THEN validation rejects the status with an error indicating `failed` is an execution concern and not a valid plan status
