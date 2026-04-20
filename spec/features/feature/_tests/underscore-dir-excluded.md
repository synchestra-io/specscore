# Scenario: Underscore-prefixed directory excluded from feature index

**Validates:** [feature#req:underscore-reserved](../README.md#req-underscore-reserved)

## Steps

GIVEN a feature directory `spec/features/todo-item/` with a valid README
AND it contains child directories `_acs/`, `_tests/`, and `completion/`
WHEN the spec tooling builds the feature tree
THEN `completion` appears as a sub-feature of `todo-item`
AND `_acs` and `_tests` do not appear in the feature tree
AND `_acs` and `_tests` do not appear in the Contents table

---
*This document follows the https://specscore.md/scenario-specification*
