# Scenario: Cross-reference uses correct identifier syntax

**Validates:** [requirement#req:ref-syntax](../README.md#req-ref-syntax)

## Steps

GIVEN a scenario file with `**Validates:** [todo-item/manage#req:title-required](../README.md#req-title-required)`
AND an AC section with `**Requirements:** todo-item/manage#req:title-required`
WHEN the spec linter validates the references
THEN both references resolve to the correct requirement heading
AND no validation errors are reported

---
*This document follows the https://specscore.md/scenario-specification*
