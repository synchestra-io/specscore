# Scenario: Header field validation

**Validates:** [plan#req:required-header-fields](../README.md#req-required-header-fields), [plan#req:source-type-values](../README.md#req-source-type-values)

## Steps

GIVEN a plan document with Status, Features, Source type, Source, Author, and Created fields all present
AND Source type is `feature`
WHEN the document is validated
THEN validation passes for the header fields

GIVEN a plan document missing the `**Author:**` field
WHEN the document is validated
THEN validation rejects the document with an error listing `Author` as a missing required field

GIVEN a plan with status `approved` that is missing the `**Approver:**` field
WHEN the document is validated
THEN validation rejects the document with an error indicating `Approver` is required when status is `approved`

GIVEN a plan with Source type set to `task`
WHEN the document is validated
THEN validation rejects the document with an error indicating Source type must be `feature` or `change-request`
