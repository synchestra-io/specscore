# Scenario: Plan directory structure

**Validates:** [development-plan#req:plan-directory](../README.md#req-plan-directory), [development-plan#req:plan-slug-format](../README.md#req-plan-slug-format)

## Steps

GIVEN a spec repository with a `spec/plans/` directory
WHEN a new plan is created with slug `add-batch-mode`
THEN a directory `spec/plans/add-batch-mode/` is created
AND a `README.md` file exists inside that directory

GIVEN a plan slug containing uppercase letters or underscores (e.g., `Add_Batch`)
WHEN the slug is validated
THEN validation rejects the slug with an error indicating it must be lowercase, hyphen-separated, and URL-safe
