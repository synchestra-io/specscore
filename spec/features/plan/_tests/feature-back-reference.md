# Scenario: Feature back-reference

**Validates:** [plan#req:feature-back-reference](../README.md#req-feature-back-reference)

## Steps

GIVEN a plan `user-auth` that lists features `api` and `ui/hub`
WHEN the feature README for `api` is checked
THEN it contains a `## Plans` section with a table row linking to `user-auth`
AND the table includes columns for Plan, Status, Author, and Approved

GIVEN a feature that appears in both a plan and one of its sub-plans
WHEN the feature README is checked
THEN both the plan and the sub-plan appear as separate rows in the Plans table
