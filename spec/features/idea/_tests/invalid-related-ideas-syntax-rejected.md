# Scenario: Unknown `Related Ideas` relationship is rejected

**Validates:** [idea#ac:related-ideas](../README.md#ac-related-ideas), [idea#req:related-ideas-format](../README.md#req-related-ideas-format)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with valid content
AND the file contains `**Related Ideas:** blocks:payment-rails-audit`
WHEN the spec linter validates the idea
THEN the linter reports an error: relationship "blocks" is not a recognized `Related Ideas` relationship
AND the error message lists the valid relationships: depends_on, alternative_to, extends, conflicts_with

---
*This document follows the https://specscore.md/scenario-specification*
