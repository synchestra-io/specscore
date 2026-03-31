# Scenario: Feature README without Acceptance Criteria section is rejected

**Validates:** [acceptance-criteria#ac:empty-section-handling](../README.md#ac-empty-section-handling), [acceptance-criteria#req:section-required](../README.md#req-section-required)

## Steps

GIVEN a feature README that contains all required sections except `## Acceptance Criteria`
WHEN the spec linter validates the feature
THEN the linter reports an error: Acceptance Criteria section is required
AND the validation fails
