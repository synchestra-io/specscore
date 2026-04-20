# Scenario: Reference to non-existent resource produces an error

**Validates:** [source-references#req:nonexistent-is-error](../README.md#req-nonexistent-is-error)

## Steps

GIVEN a source file contains `// specscore:feature/does-not-exist`
AND the path `spec/features/does-not-exist` does not exist in the repository
AND `does-not-exist` does not exist as a repo-root-relative path either
WHEN the linter validates the file
THEN the linter reports an error with the file path and line number
AND the error indicates that the referenced resource does not exist
AND the validation fails (not a warning)

---
*This document follows the https://specscore.md/scenario-specification*
