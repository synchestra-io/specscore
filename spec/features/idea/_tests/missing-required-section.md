# Scenario: Idea missing a required section is rejected

**Validates:** [idea#ac:idea-structure](../README.md#ac-idea-structure), [idea#req:required-sections](../README.md#req-required-sections)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with valid title and header fields
AND the file contains Problem Statement, Context, Recommended Direction, MVP Scope, Not Doing, Key Assumptions, SpecScore Integration, and Open Questions sections
AND the file does not contain an `Alternatives Considered` section
WHEN the spec linter validates the idea
THEN the linter reports an error: required section "Alternatives Considered" is missing
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
