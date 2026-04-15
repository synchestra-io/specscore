# Scenario: Idea id is derived from the filename

**Validates:** [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:id-is-slug](../README.md#req-id-is-slug)

## Steps

GIVEN an idea file at `spec/ideas/team-billing.md` with valid content
AND the file does not contain a separate `**Id:**` or `id:` field
WHEN the spec linter validates the idea
THEN the linter accepts the file
AND the idea's canonical id is reported as `team-billing`, matching the filename without `.md`
