# Scenario: Requirement slug with invalid characters is rejected

**Validates:** [requirement#req:slug-format](../README.md#req-slug-format)

## Steps

GIVEN a feature README with `#### REQ: Title_Required` under a topic heading in `## Behavior`
WHEN the spec linter validates the feature
THEN the linter reports an error: slug "Title_Required" contains invalid characters
AND the error message indicates that slugs must be lowercase, hyphen-separated, and URL-safe
