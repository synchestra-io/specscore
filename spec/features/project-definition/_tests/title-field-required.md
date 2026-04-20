# Scenario: Title field is mandatory

**Validates:** [project-definition#req:title-required](../README.md#req-title-required)

## Steps

GIVEN a `specscore-project.yaml` that contains no `title` field
WHEN SpecScore validates the project file
THEN it MUST report a validation error indicating `title` is required

---

GIVEN a `specscore-project.yaml` with `title: My Project`
WHEN SpecScore validates the project file
THEN validation MUST succeed

---
*This document follows the https://specscore.md/scenario-specification*
