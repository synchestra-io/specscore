# Scenario: Idea as a directory is rejected

**Validates:** [idea#req:single-file](../README.md#req-single-file)

## Steps

GIVEN a directory at `spec/ideas/offline-mode/` containing a `README.md`
AND no file exists at `spec/ideas/offline-mode.md`
WHEN the spec linter validates the idea tree
THEN the linter reports an error: idea "offline-mode" must be a single markdown file, not a directory
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
