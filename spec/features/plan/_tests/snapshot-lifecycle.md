# Scenario: Snapshot lifecycle

**Validates:** [plan#req:snapshot-table-format](../README.md#req-snapshot-table-format), [plan#req:snapshot-actions](../README.md#req-snapshot-actions), [plan#req:snapshot-git-hash](../README.md#req-snapshot-git-hash)

## Steps

GIVEN a plan with an `## Snapshots` section containing a table
WHEN the snapshot table is validated
THEN it contains columns: Date, Git Hash, Action, and Comment

GIVEN a plan that transitions to `approved` status
WHEN a snapshot is created with action `approved` and a valid git commit hash
THEN the snapshot row is added to the Snapshots table
AND the git hash references a commit where the plan content matches the approved state

GIVEN a plan with an existing `approved` snapshot
WHEN the plan content is edited and a new snapshot is created with action `checkpoint`
THEN the new snapshot row is appended to the Snapshots table
AND the previous `approved` snapshot remains unchanged

GIVEN a snapshot with a git hash that does not exist in the repository
WHEN the snapshot table is validated
THEN validation rejects with an error indicating the git hash is invalid

GIVEN a plan with a snapshot using a user-defined action `milestone-1`
WHEN the snapshot table is validated
THEN validation passes because user-defined actions are permitted
