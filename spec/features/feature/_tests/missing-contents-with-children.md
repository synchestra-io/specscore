# Scenario: Feature with children but no Contents section is rejected

**Validates:** feature/contents-when-children

## Steps

GIVEN a feature directory `spec/features/ui/` with sub-feature directories `hub/` and `tui/`
AND the `ui/README.md` does not contain a Contents section
WHEN the spec linter validates the feature
THEN the linter reports an error: feature has child directories but no Contents section
AND the validation fails
