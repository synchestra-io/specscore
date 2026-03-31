# Scenario: Undefined ACs trigger a mandatory Outstanding Question

**Validates:** feature#ac:empty-state-text, feature#req:ac-section

## Steps

GIVEN a feature README where the Acceptance Criteria section reads "Not defined yet."
AND the Outstanding Questions section does not include "Acceptance criteria not yet defined for this feature."
WHEN the spec linter validates the feature
THEN the linter reports an error: missing required Outstanding Question for undefined ACs
AND the validation fails
