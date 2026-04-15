# Scenario: `Specified` status without a `Promotes To` target is rejected

**Validates:** [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-requires-promotion](../README.md#req-specified-requires-promotion)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Status:** Specified`
AND the `**Promotes To:**` field value is `—`
WHEN the spec linter validates the idea
THEN the linter reports an error: an idea with `Status: Specified` must have a non-empty `Promotes To` list
AND the validation fails
