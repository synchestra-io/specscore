# Scenario: Hand-authored idea lints identically to a skill-authored idea

**Validates:** [idea#ac:authoring-independence](../README.md#ac-authoring-independence), [idea#req:authoring-agnostic](../README.md#req-authoring-agnostic)

## Steps

GIVEN an idea at `spec/ideas/offline-mode.md` written by hand, with all required sections, a non-empty `Not Doing` list, and at least one Must-be-true assumption
AND an identical file would have been produced by the `spec-studio:ideate` skill
WHEN the spec linter validates the idea
THEN the linter accepts the file without error or warning
AND no lint rule inspects or references authoring provenance

---
*This document follows the https://specscore.md/scenario-specification*
