# Scenario: Custom project_dirs override defaults

**Validates:** [project-definition#req:directory-layout](../README.md#req-directory-layout)

## Steps

GIVEN a `specscore-project.yaml` with `project_dirs.specifications: specifications` and `project_dirs.documents: documentation`
WHEN SpecScore resolves the specifications directory
THEN it MUST use `specifications` as the specifications path

---

GIVEN the same project file
WHEN SpecScore resolves the documents directory
THEN it MUST use `documentation` as the documents path

---
*This document follows the https://specscore.md/scenario-specification*
