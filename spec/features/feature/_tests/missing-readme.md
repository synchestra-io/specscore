# Scenario: Feature directory without README is rejected

**Validates:** feature/directory-readme

## Steps

GIVEN a feature directory `spec/features/my-feature/` exists
AND the directory does not contain a `README.md` file
WHEN the spec linter validates the feature tree
THEN the linter reports an error for `my-feature`: "Feature directory missing README.md"
AND the validation fails
