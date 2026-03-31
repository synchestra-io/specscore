# Scenario: Unrecognized status value is rejected

**Validates:** [feature#ac:readme-structure](../README.md#ac-readme-structure), [feature#req:status-field](../README.md#req-status-field)

## Steps

GIVEN a feature README with title `# Feature: Authentication`
AND the status field reads `**Status:** Draft`
WHEN the spec linter validates the feature
THEN the linter reports an error: "Draft" is not a valid status
AND the error message lists the valid values: Conceptual, In Progress, Stable, Deprecated
