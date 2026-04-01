# Scenario: Parallel eligibility

**Validates:** [plan#req:parallel-eligibility](../README.md#req-parallel-eligibility)

## Steps

GIVEN a plan with three tasks
AND Task 1 has no `Depends on` field
AND Task 2 has no `Depends on` field
AND Task 3 depends on Task 1
WHEN the dependency graph is computed
THEN Task 1 and Task 2 are marked as parallel-eligible
AND Task 3 is sequenced after Task 1
AND the critical path includes Task 1 followed by Task 3
