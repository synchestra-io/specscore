# Scenario: AC wrapping a single requirement without added value is flagged

**Validates:** [acceptance-criteria#req:no-duplicate-conditions](../README.md#req-no-duplicate-conditions)

## Steps

GIVEN a feature README with an Acceptance Criteria section
AND an AC lists exactly one requirement in its `**Requirements:**` field
AND the AC's condition statement restates the requirement verbatim without additional context
WHEN the spec linter validates the feature
THEN the linter reports a warning: AC should be omitted when it wraps a single requirement without added value
