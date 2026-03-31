# Scenario: Requirement outside Behavior section is rejected

**Validates:** [requirement#req:topic-scoped](../README.md#req-topic-scoped)

## Steps

GIVEN a feature README with a `#### REQ: some-rule` heading under `## Acceptance Criteria`
AND the requirement does not appear under `## Behavior`
WHEN the spec linter validates the feature
THEN the linter reports an error: requirement "some-rule" is outside `## Behavior`
AND the validation fails
