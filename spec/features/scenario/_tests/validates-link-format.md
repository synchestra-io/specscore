# Scenario: Validates field uses correct link format

**Validates:** [scenario#req:validates-format](../README.md#req-validates-format), [scenario#req:validates-optional](../README.md#req-validates-optional)

## Steps

GIVEN a scenario file with the validates field:
  `**Validates:** [todo#req:title-required](../README.md#req-title-required)`
WHEN the spec validator checks the validates metadata
THEN the reference is parsed as feature-path `todo`, type `req`, slug `title-required`
AND the markdown link target `../README.md#req-title-required` is validated
AND the scenario is accepted

GIVEN a scenario file with no `**Validates:**` field
WHEN the spec validator checks the scenario
THEN no validation error is reported for the missing field
AND the scenario is accepted as an exploratory or cross-feature example
