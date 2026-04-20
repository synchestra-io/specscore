# Scenario: Default values applied for omitted project_dirs

**Validates:** [project-definition#req:optional-field-defaults](../README.md#req-optional-field-defaults)

## Steps

GIVEN a `specscore-project.yaml` with only `title: My Project` and no `project_dirs`
WHEN SpecScore resolves the specifications directory
THEN it MUST use `spec` as the specifications path

---

GIVEN a `specscore-project.yaml` with only `title: My Project` and no `project_dirs`
WHEN SpecScore resolves the documents directory
THEN it MUST use `docs` as the documents path

---
*This document follows the https://specscore.md/scenario-specification*
