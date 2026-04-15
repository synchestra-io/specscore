# Scenario: Author-written `Status: Specified` without matching feature is rejected

**Validates:** [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-not-author-set](../README.md#req-specified-not-author-set)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` manually edited to `**Status:** Specified`
AND no feature in `spec/features/` lists `offline-mode` in its `**Source Ideas:**` field
WHEN the spec linter validates the idea tree
THEN the linter reports an error: `Status: Specified` may only be set by tooling in response to a referencing feature
AND the validation fails
