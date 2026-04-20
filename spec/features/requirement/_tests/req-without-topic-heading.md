# Scenario: Requirement directly under Behavior without topic heading is rejected

**Validates:** [requirement#req:topic-scoped](../README.md#req-topic-scoped)

## Steps

GIVEN a feature README with a `#### REQ: orphan-rule` heading directly under `## Behavior`
AND there is no intervening `###` topic heading between `## Behavior` and the requirement
WHEN the spec linter validates the feature
THEN the linter reports an error: requirement "orphan-rule" has no parent topic heading
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
