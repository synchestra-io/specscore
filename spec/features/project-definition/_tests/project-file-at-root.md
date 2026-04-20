# Scenario: Project file must exist at repository root

**Validates:** [project-definition#req:project-file-location](../README.md#req-project-file-location)

## Steps

GIVEN a repository with no `specscore-project.yaml` at the root
WHEN SpecScore attempts to discover the project
THEN it MUST report an error indicating the project file is missing

---

GIVEN a repository with `specscore-project.yaml` at the root
WHEN SpecScore attempts to discover the project
THEN it MUST successfully load the project configuration

---
*This document follows the https://specscore.md/scenario-specification*
