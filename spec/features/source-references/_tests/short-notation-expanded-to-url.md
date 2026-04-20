# Scenario: Short notation is auto-expanded to canonical URL before commit

**Validates:** [source-references#req:canonical-url-form](../README.md#req-canonical-url-form), [source-references#req:auto-expansion](../README.md#req-auto-expansion)

## Steps

GIVEN a source file contains the comment `// specscore:feature/cli/task/claim`
AND the git remote origin is `git@github.com:acme/myproject.git`
AND the path `spec/features/cli/task/claim` exists in the repository
WHEN the linter runs with `--fix` (or the pre-commit hook runs)
THEN the comment is replaced with `// https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim`
AND no `specscore:` prefixed reference remains in the file

---
*This document follows the https://specscore.md/scenario-specification*
