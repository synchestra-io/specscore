# Scenario: Requirement without RFC 2119 keywords triggers a warning

**Validates:** [requirement#req:rfc2119-keywords](../README.md#req-rfc2119-keywords)

## Steps

GIVEN a feature README with a requirement:
```markdown
#### REQ: title-required
A todo item needs to have a title that is not empty.
```
AND the requirement body contains no RFC 2119 keywords (MUST, SHOULD, MAY, etc.)
WHEN the spec linter validates the feature
THEN the linter reports a warning: requirement "title-required" does not use RFC 2119 keywords
AND the warning suggests rephrasing with MUST, SHOULD, or MAY

---
*This document follows the https://specscore.md/scenario-specification*
