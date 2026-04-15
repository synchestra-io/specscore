# Scenario: Author-edited `Promotes To` is rejected

**Validates:** [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:promotes-to-managed](../README.md#req-promotes-to-managed)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Status:** Approved`
AND an author has manually written `**Promotes To:** offline-sync`
AND no feature in `spec/features/` lists `offline-mode` in its `**Source Ideas:**` field
WHEN the spec linter validates the idea tree
THEN the linter reports an error: `Promotes To` disagrees with referencing features
AND the error message notes that `Promotes To` is managed state and must not be edited manually
