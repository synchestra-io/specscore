# Scenario: AC without Requirements metadata is rejected

**Validates:** [acceptance-criteria#ac:well-formed-ac](../README.md#ac-well-formed-ac), [acceptance-criteria#req:requirements-metadata](../README.md#req-requirements-metadata)

## Steps

GIVEN a feature README with an Acceptance Criteria section
AND the section contains `### AC: my-check` followed by a prose statement but no `**Requirements:**` field
WHEN the spec linter validates the feature
THEN the linter reports an error: AC must include a `**Requirements:**` field
AND the validation fails
