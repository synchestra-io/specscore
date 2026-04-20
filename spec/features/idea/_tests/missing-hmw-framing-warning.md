# Scenario: Problem Statement without `How Might We` framing emits a warning

**Validates:** [idea#req:hmw-framing](../README.md#req-hmw-framing)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with all required sections present
AND the `## Problem Statement` section contains a descriptive sentence that does not begin with "How might we"
WHEN the spec linter validates the idea
THEN the linter emits a warning (not an error): Problem Statement should contain exactly one "How Might We…" sentence
AND the validation succeeds at error severity

---
*This document follows the https://specscore.md/scenario-specification*
