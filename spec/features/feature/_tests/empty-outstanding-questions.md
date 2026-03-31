# Scenario: Outstanding Questions with no items uses placeholder text

**Validates:** feature#ac:empty-state-text, feature#req:outstanding-questions

## Steps

GIVEN a feature README with all required sections
AND the feature has no open questions
WHEN the author writes the Outstanding Questions section
THEN the section contains the text "None at this time."
AND the spec linter accepts the feature as valid
