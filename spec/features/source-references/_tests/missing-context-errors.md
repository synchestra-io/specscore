# Scenario: Missing host/org/repo context produces an error

**Validates:** [source-references#req:unresolvable-context-error](../README.md#req-unresolvable-context-error)

## Steps

GIVEN a source file contains `// specscore:feature/cli/task/claim`
AND there is no git remote named `origin`
AND `specscore-project.yaml` does not declare `project.host`, `project.org`, or `project.repo`
WHEN the linter attempts to expand the reference
THEN the linter reports an error indicating that host/org/repo cannot be resolved
AND expansion does not proceed

---
*This document follows the https://specscore.md/scenario-specification*
