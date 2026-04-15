# Scenario: Archived idea index entries out of chronological order are rejected and auto-fixed

**Validates:** [idea#ac:archival](../README.md#ac-archival), [idea#req:archived-index-chronological](../README.md#req-archived-index-chronological)

## Steps

GIVEN `spec/ideas/archived/README.md` lists archived ideas with entries of the form `- YYYY-MM-DD — [slug](<slug>.md) — <archive reason>`
AND the entries are not ordered by `**Date:**` oldest first (e.g. a `2025-03-10` entry appears before a `2024-11-02` entry)
WHEN the spec linter validates the archived index
THEN the linter reports an error: archived index entries must appear in chronological order by Date (oldest first)
AND the validation fails
WHEN `specscore lint --fix` is run against the same tree
THEN the linter rewrites `spec/ideas/archived/README.md` so entries appear oldest-first by `**Date:**`
AND a subsequent `specscore lint` run passes
