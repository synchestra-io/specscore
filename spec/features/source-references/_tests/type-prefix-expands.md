# Scenario: Type prefix expands to correct repo path

**Validates:** [source-references#req:type-prefix-expansion](../README.md#req-type-prefix-expansion)

## Steps

GIVEN a source file contains the comment `// specscore:feature/cli/task/claim`
AND the directory `spec/features/cli/task/claim` exists in the repository
WHEN the resolver processes the reference
THEN the reference resolves to the repo path `spec/features/cli/task/claim`
AND the type prefix `feature/` does not appear in the resolved path
