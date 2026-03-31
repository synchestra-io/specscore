# Scenario: Reference without comment prefix is not detected

**Validates:** [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required)

## Steps

GIVEN a source file contains the line `specscore:feature/cli/task/claim` without any comment prefix
WHEN the detector scans the file
THEN the line is not detected as a source reference
