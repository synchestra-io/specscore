---
format: https://specscore.md/scenario-specification
---

# Scenario: Plan repository routing is explicit and layered

**Validates:** [repo-config#req:plans-repo-project-selection](../README.md#req-plans-repo-project-selection), [repo-config#req:plan-repos-aggregate-routing](../README.md#req-plan-repos-aggregate-routing), [repo-config#req:plan-route-required](../README.md#req-plan-route-required), [repo-config#req:plan-config-precedence](../README.md#req-plan-config-precedence)

## Steps

GIVEN source project `github.com/acme/app`
AND user config maps `acme/plans-user: [acme/app]`
AND organization config maps `github.com/acme/plans-org: [github.com/acme/app]`
AND committed `specscore.yaml` sets `plans_repo: acme/plans-project`
AND `specscore.local.yaml` sets `plans_repo: github.com/acme/plans-local`
WHEN a Plan operation resolves its destination
THEN it selects `github.com/acme/plans-local`
AND removing each more-specific selection yields project, organization, then user in that order

GIVEN `plans_repo: acme/app` in the source repository
WHEN Plan storage is resolved
THEN same-repository storage is selected explicitly at `spec/plans`

GIVEN no route matches `github.com/acme/app`
WHEN a Plan operation starts
THEN it fails before Plan access with configuration guidance
AND a Feature lint operation for the same project proceeds without Plan routing

GIVEN one config maps `acme/app` to two different normalized destinations
WHEN Plan routing is resolved
THEN it fails as ambiguous

---
*This document follows the https://specscore.md/scenario-specification*
