# Scenario: Undefined ACs without corresponding Outstanding Question fails validation

**Validates:** [acceptance-criteria#ac:empty-section-handling](../README.md#ac-empty-section-handling), [acceptance-criteria#req:outstanding-question-linkage](../README.md#req-outstanding-question-linkage)

## Steps

GIVEN a feature README where the Acceptance Criteria section reads "Not defined yet."
AND the Outstanding Questions section does not include "Acceptance criteria not yet defined for this feature."
WHEN the spec linter validates the feature
THEN the linter reports an error: missing required Outstanding Question for undefined ACs
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
