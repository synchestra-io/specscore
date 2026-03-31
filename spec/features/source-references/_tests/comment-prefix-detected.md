# Scenario: Reference preceded by comment prefix is detected

**Validates:** [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required)

## Steps

GIVEN a source file contains the line `// specscore:feature/model-selection`
AND another line contains `# specscore:feature/model-selection`
AND another line contains `-- specscore:feature/model-selection`
AND another line contains `; specscore:feature/model-selection`
WHEN the detector scans the file
THEN all four lines are detected as valid source references
