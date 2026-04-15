# Scenario: `depends_on` cycle between ideas is accepted

**Validates:** [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:cycles-allowed](../README.md#req-cycles-allowed)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Related Ideas:** depends_on:sync-engine`
AND an idea file at `spec/ideas/sync-engine.md` with `**Related Ideas:** depends_on:offline-mode`
WHEN the spec linter validates the idea tree
THEN the linter accepts both files without error or warning about the cycle
AND any traversal tooling terminates safely rather than looping forever
