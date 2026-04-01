# Scenario: Plan document format

**Validates:** [plan#req:plan-title-format](../README.md#req-plan-title-format), [plan#req:plan-required-sections](../README.md#req-plan-required-sections)

## Steps

GIVEN a plan document with title `# Plan: Add batch mode to CLI`
WHEN the document is validated
THEN validation passes for the title format

GIVEN a plan document with title `# Add batch mode to CLI` (missing `Plan:` prefix)
WHEN the document is validated
THEN validation rejects the document with an error indicating the `Plan:` prefix is required

GIVEN a plan document that is missing the `## Context` section
WHEN the document is validated
THEN validation rejects the document with an error listing the missing required section

GIVEN a plan document that is missing the `## Tasks` section
WHEN the document is validated
THEN validation rejects the document with an error listing the missing required section
