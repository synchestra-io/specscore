# Scenario: Recursive nesting

**Validates:** [plan#req:recursive-nesting](../README.md#req-recursive-nesting), [plan#req:child-plan-format](../README.md#req-child-plan-format)

## Steps

GIVEN a plan `chat-feature` with a sub-plan `chat-infrastructure`
AND `chat-infrastructure` contains a sub-plan `database-setup`
AND `database-setup` contains a leaf task `create-schema`
WHEN the plan hierarchy is validated
THEN validation passes with no nesting depth errors

GIVEN a plan nested five levels deep
WHEN the plan hierarchy is validated
THEN validation passes because there is no maximum nesting limit

GIVEN a sub-plan `chat-infrastructure` nested under `chat-feature`
WHEN the sub-plan document is validated
THEN it follows the same format as a top-level plan, including title, header fields, Context, Acceptance criteria, and Tasks sections
