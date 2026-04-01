# Scenario: Dependency resolution

**Validates:** [task#req:dependency-sibling](../README.md#req-dependency-sibling), [task#req:dependency-cousin](../README.md#req-dependency-cousin), [task#req:dependency-cross-project](../README.md#req-dependency-cross-project), [task#req:dependency-blocks-execution](../README.md#req-dependency-blocks-execution)

## Steps

GIVEN a plan with tasks `set-up-infrastructure` and `implement-auth`
AND `implement-auth` has `depends_on: set-up-infrastructure`
WHEN the dependency is resolved
THEN `set-up-infrastructure` is found as a sibling task in the same plan directory

GIVEN a plan with task `implement-auth`
AND `implement-auth` has `depends_on: setup-db`
AND no sibling task `setup-db` exists in the same plan
WHEN the dependency is resolved
THEN resolution fails with an error indicating the sibling task was not found

GIVEN a project with plans `auth-plan` and `infra-plan`
AND task `implement-auth` in `auth-plan` has `depends_on: ../infra-plan/configure-db`
WHEN the dependency is resolved
THEN `configure-db` is found in the `infra-plan` directory via relative path resolution

GIVEN a task with `depends_on: ../nonexistent-plan/some-task`
WHEN the dependency is resolved
THEN resolution fails with an error indicating the target path does not exist

GIVEN a task with `depends_on: https://github.com/org/other-project/spec/plans/setup/configure-db`
WHEN the dependency is resolved
THEN the URL is accepted as a cross-project reference to the remote task

GIVEN a task with status `queued`
AND its dependency `set-up-infrastructure` has status `in_progress`
WHEN the task attempts to transition to `in_progress`
THEN the transition is blocked because the dependency is not `complete`

GIVEN a task with status `queued`
AND its dependency `set-up-infrastructure` has status `complete`
WHEN the task attempts to transition to `in_progress`
THEN the transition is accepted because all dependencies are met

GIVEN a task with status `in_progress`
AND a previously `complete` dependency becomes invalidated (e.g., marked `failed` after retry)
WHEN the dependency state is re-evaluated
THEN the task transitions to `blocked` because a dependency is no longer `complete`

GIVEN a task with multiple dependencies
AND all dependencies have status `complete`
WHEN the task attempts to transition to `in_progress`
THEN the transition is accepted

GIVEN a task with multiple dependencies
AND one dependency has status `complete` and another has status `in_progress`
WHEN the task attempts to transition to `in_progress`
THEN the transition is blocked because not all dependencies are `complete`
