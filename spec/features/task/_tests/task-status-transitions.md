# Scenario: Task status transitions

**Validates:** [task#req:valid-task-statuses](../README.md#req-valid-task-statuses), [task#req:status-transitions](../README.md#req-status-transitions), [task#req:terminal-states](../README.md#req-terminal-states)

## Steps

GIVEN a task with status `planning`
WHEN the status is changed to `queued`
THEN the transition is accepted

GIVEN a task with status `queued`
WHEN the status is changed to `in_progress`
THEN the transition is accepted

GIVEN a task with status `in_progress`
WHEN the status is changed to `blocked`
THEN the transition is accepted

GIVEN a task with status `blocked`
WHEN the status is changed to `in_progress`
THEN the transition is accepted

GIVEN a task with status `in_progress`
WHEN the status is changed to `complete`
THEN the transition is accepted

GIVEN a task with status `in_progress`
WHEN the status is changed to `failed`
THEN the transition is accepted

GIVEN a task with status `failed`
WHEN the status is changed to `queued` (retry)
THEN the transition is accepted

GIVEN a task with status `planning`
WHEN the status is changed to `aborted`
THEN the transition is accepted

GIVEN a task with status `queued`
WHEN the status is changed to `aborted`
THEN the transition is accepted

GIVEN a task with status `in_progress`
WHEN the status is changed to `aborted`
THEN the transition is accepted

GIVEN a task with status `blocked`
WHEN the status is changed to `aborted`
THEN the transition is accepted

GIVEN a task with status `complete`
WHEN the status is changed to any other status
THEN the transition is rejected with an error indicating `complete` is a terminal state

GIVEN a task with status `aborted`
WHEN the status is changed to any other status
THEN the transition is rejected with an error indicating `aborted` is a terminal state

GIVEN a task with status `planning`
WHEN the status is changed directly to `in_progress` (skipping `queued`)
THEN the transition is rejected with an error indicating the allowed transitions from `planning`

GIVEN a task with status `queued`
WHEN the status is changed to `complete` (skipping `in_progress`)
THEN the transition is rejected with an error indicating the allowed transitions from `queued`

GIVEN a task with status `planning`
WHEN the status is changed to `invalid_status`
THEN the status is rejected with an error indicating the value is not a valid task status
