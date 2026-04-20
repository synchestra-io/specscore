# Scenario: Type prefix expansion falls back when expanded path does not exist

**Validates:** [source-references#req:resolution-order](../README.md#req-resolution-order)

## Steps

GIVEN a source file contains the comment `// specscore:feature/nonexistent`
AND the directory `spec/features/nonexistent` does not exist
AND a file or directory `feature/nonexistent` exists at the repo root
WHEN the resolver processes the reference
THEN the resolver falls back to treating `feature/nonexistent` as a repo-root-relative path
AND the reference resolves to `feature/nonexistent`

---
*This document follows the https://specscore.md/scenario-specification*
