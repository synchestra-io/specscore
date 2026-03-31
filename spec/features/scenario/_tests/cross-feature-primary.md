# Scenario: Cross-feature scenario placed in primary feature

**Validates:** [scenario#req:cross-feature-placement](../README.md#req-cross-feature-placement), [scenario#req:cross-feature-validates](../README.md#req-cross-feature-validates)

## Steps

GIVEN two features: `due-dates` (primary) and `todo-list` (secondary)
AND a scenario file in `spec/features/due-dates/_tests/overdue-with-filter.md` with:
  `**Validates:** [due-dates#req:overdue-detection](../README.md#req-overdue-detection), [todo-list#req:filter-by-status](../../todo-list/README.md#req-filter-by-status)`
WHEN the spec validator checks the scenario
THEN the scenario is accepted in the primary feature's `_tests/` directory
AND both cross-feature validates references are resolved
AND no duplicate scenario exists in `spec/features/todo-list/_tests/`
