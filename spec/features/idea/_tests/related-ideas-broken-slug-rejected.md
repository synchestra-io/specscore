# Scenario: `Related Ideas` entry with a non-existent slug is rejected

**Validates:** [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:related-ideas-target-exists](../README.md#req-related-ideas-target-exists)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Related Ideas:** depends_on:ghost-idea`
AND no file exists at `spec/ideas/ghost-idea.md` or `spec/ideas/archived/ghost-idea.md`
WHEN the spec linter validates the idea tree
THEN the linter reports an error: `Related Ideas` target "ghost-idea" does not resolve to an existing idea
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
