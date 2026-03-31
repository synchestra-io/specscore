# Scenario: AC with incorrect heading format is rejected

**Validates:** [acceptance-criteria#ac:well-formed-ac](../README.md#ac-well-formed-ac), [acceptance-criteria#req:inline-heading](../README.md#req-inline-heading)

## Steps

GIVEN a feature README with an Acceptance Criteria section
AND the section contains a heading `### Acceptance Criteria: my-ac` instead of `### AC: my-ac`
WHEN the spec linter validates the feature
THEN the linter reports an error: AC heading must use `### AC: {slug}` format
AND the validation fails
