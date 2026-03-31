# Scenario: README without Feature prefix in title is rejected

**Validates:** feature#ac:readme-structure, feature#req:title-format

## Steps

GIVEN a feature directory `spec/features/authentication/` with a `README.md`
AND the README title is `# Authentication` (missing `Feature:` prefix)
WHEN the spec linter validates the feature
THEN the linter reports an error: title must use `# Feature: {Title}` format
AND the validation fails
