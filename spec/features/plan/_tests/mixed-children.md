# Scenario: Mixed children

**Validates:** [plan#req:mixed-children](../README.md#req-mixed-children)

## Steps

GIVEN a plan `chat-feature` with the following children:
  - `chat-infrastructure/` (a sub-plan with its own tasks)
  - `chat-workflow-engine/` (a sub-plan with its own tasks)
  - `send-notifications/` (a leaf task with no children)
WHEN the plan structure is validated
THEN validation passes because leaf tasks and sub-plans may coexist at the same level

GIVEN a plan where all children are leaf tasks (no sub-plans)
WHEN the plan structure is validated
THEN validation passes

GIVEN a plan where all children are sub-plans (no leaf tasks)
WHEN the plan structure is validated
THEN validation passes
