# Scenario: Flow files live in flows subdirectory

**Validates:** [scenario#req:flow-location](../README.md#req-flow-location), [scenario#req:flow-format](../README.md#req-flow-format)

## Steps

GIVEN a feature with a `_tests/flows/` subdirectory
AND a flow file `_tests/flows/populated-list.md` exists with Given/When/Then steps:
  ```
  GIVEN a new todo list
  WHEN the user runs `todo add "Item 1"` and `todo add "Item 2"`
  THEN the list contains two items
  ```
AND a scenario file `_tests/complete-item.md` references the flow:
  ```
  GIVEN the state from [populated-list](flows/populated-list.md)
  WHEN the user runs `todo complete 1`
  THEN item 1 is marked complete
  ```
WHEN the spec validator checks both files
THEN the flow file is accepted as valid (correct location and format)
AND the scenario's flow reference is resolved successfully
