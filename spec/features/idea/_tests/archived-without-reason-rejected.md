# Scenario: Archived idea without an `Archive Reason` is rejected

**Validates:** [idea#ac:archival](../README.md#ac-archival), [idea#req:archive-reason](../README.md#req-archive-reason)

## Steps

GIVEN an idea file at `spec/ideas/archived/offline-mode.md` with `**Status:** Archived`
AND the `**Archive Reason:**` field is absent or set to `—`
WHEN the spec linter validates the idea
THEN the linter reports an error: archived idea "offline-mode" must include a non-empty `Archive Reason`
AND the validation fails
