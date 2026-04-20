# Scenario: Empty AC section displays placeholder text

**Validates:** [acceptance-criteria#ac:empty-section-handling](../README.md#ac-empty-section-handling), [acceptance-criteria#req:placeholder-when-empty](../README.md#req-placeholder-when-empty)

## Steps

GIVEN a feature README where the Acceptance Criteria section has no `### AC:` headings
AND the section does not contain "Not defined yet."
WHEN the spec linter validates the feature
THEN the linter reports an error: Acceptance Criteria section with no ACs must state "Not defined yet."
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
