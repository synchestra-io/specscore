# Scenario: Requirement with multiple independent obligations is flagged

**Validates:** [requirement#req:single-obligation](../README.md#req-single-obligation)

## Steps

GIVEN a feature README with a requirement:
```markdown
#### REQ: title-rules
A todo item MUST have a non-empty title, MUST NOT exceed 200 characters,
and SHOULD be unique within the list.
```
WHEN the spec linter validates the feature
THEN the linter reports a warning: requirement "title-rules" contains multiple independent obligations
AND the linter suggests splitting into separate requirements

---
*This document follows the https://specscore.md/scenario-specification*
