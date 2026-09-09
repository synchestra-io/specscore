---
format: https://specscore.md/scenario-specification
---

# Scenario: External Plan checkout paths stay machine-local

**Validates:** [repo-config#req:repo-checkouts-machine-local](../README.md#req-repo-checkouts-machine-local)

## Steps

GIVEN an external destination `github.com/acme/plans`
AND local, organization, or user config maps it under `repo_checkouts` to `/work/checkouts/acme/plans`
AND that path is the git top-level whose origin is `github.com/acme/plans`
WHEN Plan storage is resolved
THEN the checkout is accepted

GIVEN `repo_checkouts` appears in committed `specscore.yaml`
WHEN config is loaded
THEN it fails before repository access and directs the value to machine-local config

GIVEN a checkout value is relative, names a nested directory, has a different origin, or resolves through a symlink to a different repository
WHEN external Plan storage is resolved
THEN it fails without reading or writing Plan artifacts

---
*This document follows the https://specscore.md/scenario-specification*
