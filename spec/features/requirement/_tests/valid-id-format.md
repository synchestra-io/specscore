# Scenario: Requirement is addressable by canonical ID

**Validates:** [requirement#req:id-format](../README.md#req-id-format)

## Steps

GIVEN a feature at `spec/features/todo-item/manage/`
AND its README contains `#### REQ: title-required` under a topic heading in `## Behavior`
WHEN spec tooling resolves the identifier `todo-item/manage#req:title-required`
THEN the tool locates the requirement heading in `spec/features/todo-item/manage/README.md`
AND returns the requirement body text
