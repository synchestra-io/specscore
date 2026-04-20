# Scenario: Feature identified by path relative to spec/features

**Validates:** [feature#req:path-identification](../README.md#req-path-identification)

## Steps

GIVEN a feature at `spec/features/billing/payments/README.md`
WHEN the spec tooling resolves the feature identifier
THEN the identifier is `billing/payments`
AND plans, source references, and other specs use `billing/payments` to reference this feature

---
*This document follows the https://specscore.md/scenario-specification*
