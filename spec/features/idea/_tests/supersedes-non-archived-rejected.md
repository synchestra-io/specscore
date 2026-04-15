# Scenario: `Supersedes` target that is not archived is rejected

**Validates:** [idea#ac:archival](../README.md#ac-archival), [idea#req:supersedes-target-archived](../README.md#req-supersedes-target-archived)

## Steps

GIVEN an idea file at `spec/ideas/offline-sync.md` with `**Supersedes:** offline-mode`
AND `spec/ideas/offline-mode.md` exists with `**Status:** Approved` (not Archived)
WHEN the spec linter validates the idea tree
THEN the linter reports an error: `Supersedes` target "offline-mode" must exist under `spec/ideas/archived/` with `Status: Archived`
AND the validation fails
