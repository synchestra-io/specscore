# Scenario: Sub-feature must explicitly reference parent requirement

**Validates:** [requirement#req:no-inheritance](../README.md#req-no-inheritance)

## Steps

GIVEN a parent feature `billing` with `#### REQ: currency-required` in its Behavior section
AND a sub-feature `billing/payments` with a scenario that validates currency handling
AND the scenario's `**Validates:**` field references `billing/payments#req:currency-required`
BUT no such requirement exists in `billing/payments` — it only exists in `billing`
WHEN the spec linter validates the reference
THEN the linter reports an error: requirement `billing/payments#req:currency-required` not found
AND suggests referencing the parent requirement as `billing#req:currency-required`
