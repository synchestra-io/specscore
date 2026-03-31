# Scenario: README missing a required section is rejected

**Validates:** [feature#ac:readme-structure](../README.md#ac-readme-structure), [feature#req:required-sections](../README.md#req-required-sections)

## Steps

GIVEN a feature README with title, status, Summary, and Behavior sections
AND the README does not contain a Problem section
WHEN the spec linter validates the feature
THEN the linter reports an error: required section "Problem" is missing
AND the validation fails
