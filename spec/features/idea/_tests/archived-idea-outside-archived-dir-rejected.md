# Scenario: Archived idea outside `archived/` subdirectory is rejected

**Validates:** [idea#ac:archival](../README.md#ac-archival), [idea#req:archived-location](../README.md#req-archived-location)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Status:** Archived`
AND the file has not been moved into `spec/ideas/archived/`
WHEN the spec linter validates the idea tree
THEN the linter reports an error: archived idea "offline-mode" must reside at `spec/ideas/archived/offline-mode.md`
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
