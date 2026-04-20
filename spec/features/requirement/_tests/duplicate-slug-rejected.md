# Scenario: Duplicate requirement slug within a feature is rejected

**Validates:** [requirement#req:slug-unique](../README.md#req-slug-unique)

## Steps

GIVEN a feature README with two requirements both named `#### REQ: title-required`
AND they appear under different topic headings within `## Behavior`
WHEN the spec linter validates the feature
THEN the linter reports an error: duplicate requirement slug "title-required" in the same feature
AND the validation fails

---
*This document follows the https://specscore.md/scenario-specification*
