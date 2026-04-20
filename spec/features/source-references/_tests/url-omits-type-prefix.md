# Scenario: Expanded URL contains resolved path without type prefix

**Validates:** [source-references#req:url-structure](../README.md#req-url-structure)

## Steps

GIVEN a source file contains `// specscore:plan/v2-migration`
AND the git remote origin resolves to `github.com/acme/myproject`
WHEN the linter expands the reference
THEN the expanded URL is `https://specscore.org/github.com/acme/myproject/spec/plans/v2-migration`
AND the URL does not contain the segment `plan/` — only the resolved path `spec/plans/v2-migration`

---
*This document follows the https://specscore.md/scenario-specification*
