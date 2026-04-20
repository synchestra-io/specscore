# Scenario: Step keywords must be uppercase

**Validates:** [scenario#req:keywords-uppercase](../README.md#req-keywords-uppercase)

## Steps

GIVEN a scenario file with steps using lowercase keywords:
  ```
  given an empty list
  when the user adds an item
  then the list has one item
  ```
WHEN the spec validator checks the steps section
THEN a validation error is reported indicating keywords must be uppercase
AND the scenario is rejected

---
*This document follows the https://specscore.md/scenario-specification*
