---
format: https://specscore.md/scenario-specification
---

# Scenario: Plan file structure

**Validates:** [plan#req:plan-file](../README.md#req-plan-file), [plan#req:plan-slug-format](../README.md#req-plan-slug-format)

## Steps

GIVEN a spec repository with a `spec/plans/` directory
WHEN a new plan is created with slug `add-batch-mode`
THEN `spec/plans/add-batch-mode/README.md` is created
AND `spec/plans/add-batch-mode.md` is not created

GIVEN both `spec/plans/add-batch-mode.md` and `spec/plans/add-batch-mode/README.md` exist
WHEN the Plan namespace is validated
THEN validation rejects the duplicate logical ID before any mutation

GIVEN `spec/plans/platform/api/README.md` and `spec/plans/payments/api/README.md` exist
WHEN the Plan namespace is enumerated
THEN their logical IDs are `platform/api` and `payments/api`
AND neither is collapsed to the leaf name `api`

GIVEN a plan slug containing uppercase letters or underscores (e.g., `Add_Batch`)
WHEN the slug is validated
THEN validation rejects the slug with an error indicating it must be lowercase, hyphen-separated, and URL-safe

---
*This document follows the https://specscore.md/scenario-specification*
