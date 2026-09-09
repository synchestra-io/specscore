---
format: https://specscore.md/scenarios-index-specification
---

# Scenarios: Plan

Test scenarios validating the plan feature requirements.

| Scenario | Validates |
|---|---|
| [plan-file-structure](plan-file-structure.md) | [plan#req:plan-file](../README.md#req-plan-file), [plan#req:plan-logical-id](../README.md#req-plan-logical-id), [plan#req:plan-slug-format](../README.md#req-plan-slug-format) |
| [external-plan-store](external-plan-store.md) | [plan#req:external-source-namespace](../README.md#req-external-source-namespace), [plan#req:source-project-context](../README.md#req-source-project-context), [plan#req:read-only-plan-operations](../README.md#req-read-only-plan-operations), [plan#req:plan-path-conflicts](../README.md#req-plan-path-conflicts), [plan#req:external-store-lifecycle-lock](../README.md#req-external-store-lifecycle-lock) |
| [plan-document-format](plan-document-format.md) | [plan#req:plan-title-format](../README.md#req-plan-title-format), [plan#req:plan-required-sections](../README.md#req-plan-required-sections) |
| [header-field-validation](header-field-validation.md) | [plan#req:required-header-fields](../README.md#req-required-header-fields), [plan#req:source-binding](../README.md#req-source-binding) |
| [valid-status-values](valid-status-values.md) | [plan#req:valid-statuses](../README.md#req-valid-statuses), [plan#req:execution-status-derived](../README.md#req-execution-status-derived) |
| [status-transition-rules](status-transition-rules.md) | [plan#req:status-transitions](../README.md#req-status-transitions) |
| [parallel-eligibility](parallel-eligibility.md) | [plan#req:parallel-eligibility](../README.md#req-parallel-eligibility) |
| [task-verifies-feature-ac](task-verifies-feature-ac.md) | [plan#req:task-verifies-feature-ac](../README.md#req-task-verifies-feature-ac) |
| [index-structure](index-structure.md) | [plans-index#req:namespace-root-only](../../plans-index/README.md#req-namespace-root-only), [plans-index#req:required-sections](../../plans-index/README.md#req-required-sections), [plans-index#req:sub-plan-indentation](../../plans-index/README.md#req-sub-plan-indentation), [plans-index#req:contents-columns](../../plans-index/README.md#req-contents-columns), [plans-index#req:recently-closed-present](../../plans-index/README.md#req-recently-closed-present) |
| [feature-back-reference](feature-back-reference.md) | [plan#req:feature-back-reference](../README.md#req-feature-back-reference) |
| [proposal-forward-reference](proposal-forward-reference.md) | [plan#req:proposal-forward-reference](../README.md#req-proposal-forward-reference) |
| [source-binding](source-binding.md) | [plan#req:source-binding](../README.md#req-source-binding) |
| [snapshot-lifecycle](snapshot-lifecycle.md) | [plan#req:snapshot-table-format](../README.md#req-snapshot-table-format), [plan#req:snapshot-actions](../README.md#req-snapshot-actions), [plan#req:snapshot-git-hash](../README.md#req-snapshot-git-hash) |
| [recursive-nesting](recursive-nesting.md) | [plan#req:recursive-nesting](../README.md#req-recursive-nesting), [plan#req:child-plan-format](../README.md#req-child-plan-format), [plan#req:prerequisite-plan-logical-ids](../README.md#req-prerequisite-plan-logical-ids) |
| [mixed-children](mixed-children.md) | [plan#req:mixed-children](../README.md#req-mixed-children) |
| [status-rollup](status-rollup.md) | [plan#req:status-rollup](../README.md#req-status-rollup) |

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/scenarios-index-specification*
