# Scenario: Roadmap structure

**Validates:** [development-plan#req:roadmap-no-steps](../README.md#req-roadmap-no-steps), [development-plan#req:child-plan-format](../README.md#req-child-plan-format), [development-plan#req:hierarchy-nesting-limit](../README.md#req-hierarchy-nesting-limit)

## Steps

GIVEN a roadmap plan with a `## Child Plans` table listing two child plans
AND no `## Steps` section with implementation steps
WHEN the roadmap is validated
THEN validation passes

GIVEN a roadmap plan that has both a `## Child Plans` table and a `## Steps` section with implementation steps
WHEN the roadmap is validated
THEN validation rejects the document with an error indicating roadmaps must not have implementation steps

GIVEN a child plan with standard steps and acceptance criteria
WHEN the child plan is validated
THEN validation passes using the same rules as a standalone plan

GIVEN a roadmap whose child plan itself contains child plans (three levels deep)
WHEN the hierarchy is validated
THEN validation rejects with an error indicating hierarchy must not exceed two levels
