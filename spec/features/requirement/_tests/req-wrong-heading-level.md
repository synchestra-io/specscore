# Scenario: Requirement at wrong heading level is rejected

**Validates:** [requirement#req:heading-level](../README.md#req-heading-level)

## Steps

GIVEN a feature README with a `### topic` heading under `## Behavior`
AND a `##### REQ: deep-rule` heading under that topic (two levels below instead of one)
WHEN the spec linter validates the feature
THEN the linter reports an error: requirement "deep-rule" is not one heading level below its topic
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
