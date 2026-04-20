# Scenario: AC without a prose condition statement is rejected

**Validates:** [acceptance-criteria#ac:well-formed-ac](../README.md#ac-well-formed-ac), [acceptance-criteria#req:condition-statement](../README.md#req-condition-statement)

## Steps

GIVEN a feature README with an Acceptance Criteria section
AND the section contains `### AC: my-check` with a `**Requirements:**` field but no prose statement after it
WHEN the spec linter validates the feature
THEN the linter reports an error: AC must include a condition statement
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
