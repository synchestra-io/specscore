# Scenario: Feature slug with invalid characters is rejected

**Validates:** feature#req:slug-format

## Steps

GIVEN a feature directory named `spec/features/My_Feature/`
AND it contains a valid `README.md`
WHEN the spec linter validates the feature tree
THEN the linter reports an error: slug "My_Feature" contains invalid characters
AND the error message indicates that slugs must be lowercase, hyphen-separated, and URL-safe
