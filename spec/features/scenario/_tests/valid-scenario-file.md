# Scenario: Valid scenario file accepted

**Validates:** [scenario#req:title-prefix](../README.md#req-title-prefix), [scenario#req:steps-required](../README.md#req-steps-required), [scenario#req:filename-slug](../README.md#req-filename-slug)

## Steps

GIVEN a feature with a `_tests/` directory
AND a file `_tests/add-item.md` exists with the content:
  ```
  # Scenario: Add item to list
  **Validates:** [todo#req:title-required](../README.md#req-title-required)
  ## Steps
  GIVEN an empty todo list
  WHEN the user runs `todo add "Buy milk"`
  THEN the list contains one item titled "Buy milk"
  ```
WHEN the spec validator checks the scenario file
THEN the file is accepted as a valid scenario
AND no validation errors are reported
