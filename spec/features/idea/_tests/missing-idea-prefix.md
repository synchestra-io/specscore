# Scenario: Idea title without `Idea:` prefix is rejected

**Validates:** [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:title-format](../README.md#req-title-format)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md`
AND its first heading reads `# Offline Mode`
WHEN the spec linter validates the idea
THEN the linter reports an error: idea title must use the `# Idea: <Title>` format
AND the error message notes that the `Idea:` prefix is required to dispatch the idea rule set

---
*This document follows the https://specscore.md/scenario-specification*
