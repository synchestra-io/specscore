# Scenario: Plan step AC references a feature AC

**Validates:** [acceptance-criteria#req:plan-ac-references](../README.md#req-plan-ac-references)

## Steps

GIVEN a development plan step with an AC that references `my-feature#ac:input-validation`
AND the feature `my-feature` has an AC `### AC: input-validation` in its README
WHEN the spec linter resolves the plan step AC reference
THEN the reference resolves to the feature AC
AND the plan step AC is marked as frozen with the plan
AND the feature AC continues to evolve independently
