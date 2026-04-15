# Scenario: Idea missing a required header field is rejected

**Validates:** [idea#ac:idea-header](../README.md#ac-idea-header), [idea#req:header-fields](../README.md#req-header-fields)

## Steps

GIVEN an idea file at `spec/ideas/offline-mode.md` with title `# Idea: Offline Mode`
AND the file contains `**Status:** Draft` and `**Date:** 2026-04-15`
AND the file does not contain an `**Owner:**` field
WHEN the spec linter validates the idea
THEN the linter reports an error: required header field "Owner" is missing
AND the validation fails
