# Scenario: Scenario correctly references an AC in its Validates field

**Validates:** [acceptance-criteria#req:scenario-validates-ac](../README.md#req-scenario-validates-ac)

## Steps

GIVEN a feature with an AC `### AC: input-validation` in its README
AND a scenario file in `_tests/` with `**Validates:** feature-id#ac:input-validation`
WHEN the spec linter resolves the Validates reference
THEN the reference resolves to the AC heading in the feature README
AND validation passes

---
*This document follows the https://specscore.md/scenario-specification*
