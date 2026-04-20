# Scenario: Unrecognized idea status value is rejected

**Validates:** [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:status-values](../README.md#req-status-values)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with title `# Idea: Offline Mode`
AND the status field reads `**Status:** Pending`
WHEN the spec linter validates the idea
THEN the linter reports an error: "Pending" is not a valid idea status
AND the error message lists the valid values: Draft, Under Review, Approved, Specified, Archived

---
*This document follows the https://specscore.md/scenario-specification*
