# Scenario: Reference inside string literal is not detected

**Validates:** [source-references#req:comment-prefix-required](../README.md#req-comment-prefix-required)

## Steps

GIVEN a Go source file contains the line `fmt.Println("specscore:feature/x")`
AND the `specscore:` text is inside a string literal, not preceded by a comment prefix
WHEN the detector scans the file
THEN the line is not detected as a source reference
