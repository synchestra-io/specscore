# Scenario: Scenario without title prefix is rejected

**Validates:** [scenario#req:title-prefix](../README.md#req-title-prefix)

## Steps

GIVEN a feature with a `_tests/` directory
AND a file `_tests/add-item.md` exists with the heading `# Add item to list` (missing `Scenario:` prefix)
WHEN the spec validator checks the scenario file
THEN a validation error is reported indicating the `Scenario:` prefix is missing
AND the file is rejected
