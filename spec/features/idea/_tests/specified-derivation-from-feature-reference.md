# Scenario: Creating a feature that references an idea transitions it to `Specified`

**Validates:** [idea#ac:promotion-lifecycle](../README.md#ac-promotion-lifecycle), [idea#req:specified-derivation](../README.md#req-specified-derivation)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with `**Status:** Approved` and `**Promotes To:** —`
AND a feature at `spec/features/offline-sync/README.md` is created with `**Source Ideas:** offline-mode`
WHEN `specscore lint --fix` runs over the spec tree
THEN the idea's `**Status:**` is updated to `Specified`
AND the idea's `**Promotes To:**` lists `offline-sync`
AND subsequent `specscore lint` invocations succeed without drift

---
*This document follows the https://specscore.md/scenario-specification*
