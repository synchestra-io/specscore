# Scenario: AC containing Given/When/Then steps is rejected

**Validates:** [acceptance-criteria#ac:ac-scenario-separation](../README.md#ac-ac-scenario-separation), [acceptance-criteria#req:abstract-not-concrete](../README.md#req-abstract-not-concrete)

## Steps

GIVEN a feature README with an Acceptance Criteria section
AND an AC's condition statement contains "GIVEN" or "WHEN" or "THEN" keywords as step instructions
WHEN the spec linter validates the feature
THEN the linter reports a warning: AC should remain abstract and not contain Given/When/Then steps

---
*This document follows the https://specscore.md/scenario-specification*
