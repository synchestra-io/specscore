# Scenario: Features field uniform declaration

**Validates:** [plan#req:features-field-uniform](../README.md#req-features-field-uniform)

## Steps

GIVEN a plan that affects a single feature `cli`
WHEN the plan header is checked
THEN the `**Features:**` field is present with `cli` listed

GIVEN a plan that affects three features: `api`, `cli`, and `ui/hub`
WHEN the plan header is checked
THEN the `**Features:**` field lists all three features, each linking to its feature spec README

GIVEN a plan document without a `**Features:**` field
WHEN the document is validated
THEN validation rejects the document with an error indicating the Features field is required

---
*This document follows the https://specscore.md/scenario-specification*
