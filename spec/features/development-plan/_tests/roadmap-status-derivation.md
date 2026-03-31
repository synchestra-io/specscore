# Scenario: Roadmap status derivation

**Validates:** [development-plan#req:roadmap-status-derivation](../README.md#req-roadmap-status-derivation)

## Steps

GIVEN a roadmap with two child plans, one `draft` and one `approved`
WHEN the roadmap status is derived
THEN the roadmap status is `draft`

GIVEN a roadmap with two child plans, both `in_review`
WHEN the roadmap status is derived
THEN the roadmap status is `in_review`

GIVEN a roadmap with two child plans, both `approved`
WHEN the roadmap status is derived
THEN the roadmap status is `approved`

GIVEN a roadmap with two child plans, one `approved` and one with linked work items in progress
WHEN the roadmap status is derived
THEN the roadmap status is `in_progress`

GIVEN a roadmap that is explicitly set to `superseded`
WHEN the roadmap status is checked
THEN the roadmap status is `superseded` regardless of child plan statuses
