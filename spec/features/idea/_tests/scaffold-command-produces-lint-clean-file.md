# Scenario: `specscore idea new` scaffold is lint-clean on exit

**Validates:** [idea#ac:scaffold-behavior](../README.md#ac-scaffold-behavior), [idea#req:scaffold-command](../README.md#req-scaffold-command)

## Steps

GIVEN a spec repository with no file at `spec/ideas/offline-mode.md`
WHEN the user runs `specscore idea new offline-mode`
THEN a file is created at `spec/ideas/offline-mode.md` with all required sections, each prefilled with an HTML-comment prompt (e.g. `<!-- One "How Might We…" sentence. -->`)
AND `**Promotes To:**`, `**Supersedes:**`, and `**Related Ideas:**` are set to `—`
AND running `specscore lint` against the new file succeeds with no errors (placeholder rule U-005 does not trip on the HTML-comment prompts)

---
*This document follows the https://specscore.md/scenario-specification*
