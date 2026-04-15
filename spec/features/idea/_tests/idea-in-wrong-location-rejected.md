# Scenario: Idea outside `spec/ideas/` is rejected

**Validates:** [idea#req:idea-location](../README.md#req-idea-location)

## Steps

GIVEN an idea file at `docs/ideas/payment-fraud-signals.md` with otherwise valid content
AND the file is not present under `spec/ideas/` or `spec/ideas/archived/`
WHEN the spec linter validates the idea tree
THEN the linter reports an error: idea "payment-fraud-signals" is not in `spec/ideas/` or `spec/ideas/archived/`
AND the validation fails
