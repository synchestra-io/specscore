# Scenario: Parallel eligibility

**Validates:** [development-plan#req:parallel-eligibility](../README.md#req-parallel-eligibility)

## Steps

GIVEN a plan with three steps
AND Step 1 has no `Depends on` field
AND Step 2 has no `Depends on` field
AND Step 3 depends on Step 1
WHEN the dependency graph is computed
THEN Step 1 and Step 2 are marked as parallel-eligible
AND Step 3 is sequenced after Step 1
AND the critical path includes Step 1 followed by Step 3
