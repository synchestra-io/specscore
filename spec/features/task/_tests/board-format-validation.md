# Scenario: Board format validation

**Validates:** [task#req:board-columns](../README.md#req-board-columns), [task#req:board-done-strikethrough](../README.md#req-board-done-strikethrough), [task#req:board-recently-finished](../README.md#req-board-recently-finished), [task#req:board-reflects-current-state](../README.md#req-board-reflects-current-state)

## Steps

GIVEN a plan with a task status board
WHEN the board is validated
THEN the board contains columns: Queued, In Progress, Blocked, and Done

GIVEN a plan with a task status board
AND the board is missing the Blocked column
WHEN the board is validated
THEN validation fails with an error indicating the required columns

GIVEN a plan with a task in status `planning`
WHEN the task status board is rendered
THEN the task appears in the Queued column

GIVEN a plan with a task in status `complete`
WHEN the task status board is rendered
THEN the task appears in the Done column with strikethrough formatting (`~~task-slug~~`)

GIVEN a plan with a task in status `failed`
WHEN the task status board is rendered
THEN the task appears in the Done column with strikethrough formatting

GIVEN a plan with a task in status `aborted`
WHEN the task status board is rendered
THEN the task appears in the Done column with strikethrough formatting

GIVEN a plan with tasks that have reached terminal status
WHEN the task status board is rendered
THEN the board includes a Recently Finished section listing those tasks with their status and completion date

GIVEN a plan where a task has transitioned from `queued` to `in_progress`
AND the board still shows the task in the Queued column
WHEN the board is validated against current task state
THEN validation fails with an error indicating the board is stale

GIVEN a plan where all task statuses match their board column placements
WHEN the board is validated against current task state
THEN validation passes with no errors

GIVEN a plan with a Recently Finished section
WHEN the section is validated
THEN each entry includes a Task column, a Status column, and a Completed date column
