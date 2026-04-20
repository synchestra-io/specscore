# Scenario: Feature missing from index is a validation error

**Validates:** [feature#req:index-completeness](../README.md#req-index-completeness)

## Steps

GIVEN a feature directory `spec/features/notifications/` with a valid README
AND the feature index (`spec/features/README.md`) does not list `notifications`
WHEN the spec linter validates the feature tree
THEN the linter reports an error: feature "notifications" exists but is not listed in the index
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
