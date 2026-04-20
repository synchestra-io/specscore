# Scenario: Idea slug with invalid characters is rejected

**Validates:** [idea#req:slug-format](../README.md#req-slug-format)

## Steps

GIVEN an idea file at `spec/ideas/Payment_Fraud.md`
AND the file contains an otherwise valid idea body
WHEN the spec linter validates the idea tree
THEN the linter reports an error: slug "Payment_Fraud" contains invalid characters
AND the error message indicates that idea slugs must be lowercase, hyphen-separated, and URL-safe

---
*This document follows the https://specscore.md/scenario-specification*
