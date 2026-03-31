# Scenario: Multiple WHEN steps are rejected

**Validates:** [scenario#req:single-when](../README.md#req-single-when)

## Steps

GIVEN a scenario file with steps containing two WHEN keywords:
  ```
  GIVEN an empty list
  WHEN the user adds an item
  WHEN the user deletes the item
  THEN the list is empty
  ```
WHEN the spec validator checks the steps section
THEN a validation error is reported indicating exactly one WHEN is required
AND the scenario is rejected
