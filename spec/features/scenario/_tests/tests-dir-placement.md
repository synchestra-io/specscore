# Scenario: Scenarios outside _tests are rejected

**Validates:** [scenario#req:tests-dir-only](../README.md#req-tests-dir-only)

## Steps

GIVEN a feature directory `spec/features/todo/`
AND a file `spec/features/todo/add-item.md` exists with a valid `# Scenario:` heading (placed directly in the feature directory, not in `_tests/`)
WHEN the spec validator scans the feature tree
THEN a validation error is reported indicating scenarios must live in `_tests/`
AND the file is not recognized as a valid scenario
