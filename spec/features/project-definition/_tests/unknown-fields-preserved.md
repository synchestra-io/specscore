# Scenario: Unknown fields are ignored without error

**Validates:** [project-definition#req:unknown-fields-ignored](../README.md#req-unknown-fields-ignored)

## Steps

GIVEN a `specscore-project.yaml` containing `title: My Project` and an unknown field `state_repo: https://github.com/org/state`
WHEN SpecScore validates the project file
THEN validation MUST succeed without errors
AND the unknown field MUST NOT cause a warning or validation failure

---
*This document follows the https://specscore.md/scenario-specification*
