# Scenario: Only known type prefixes are recognized

**Validates:** [source-references#req:type-prefix-set](../README.md#req-type-prefix-set)

## Steps

GIVEN a source file contains `// specscore:config/database`
AND `config` is not in the known type prefix set (`feature/`, `plan/`, `doc/`)
WHEN the resolver processes the reference
THEN the resolver does not attempt type prefix expansion
AND the reference is treated as the repo-root-relative path `config/database`

---
*This document follows the https://specscore.md/scenario-specification*
