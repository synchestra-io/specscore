# Scenario: Unknown prefix falls back to repo-root-relative path

**Validates:** [source-references#req:resolution-order](../README.md#req-resolution-order)

## Steps

GIVEN a source file contains the comment `// specscore:spec/features/cli/task/claim`
AND the first segment `spec` is not a known type prefix
AND the path `spec/features/cli/task/claim` exists in the repository
WHEN the resolver processes the reference
THEN the reference resolves to the repo path `spec/features/cli/task/claim`
