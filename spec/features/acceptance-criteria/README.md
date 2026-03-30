# Feature: Acceptance Criteria

**Status:** Conceptual

## Summary

Acceptance criteria are the contract between what a feature promises and what the system actually delivers. Each AC is a standalone markdown file — readable by product owners, auditable by reviewers, and executable by the [test runner](https://github.com/synchestra-io/rehearse/blob/main/spec/features/testing-framework/test-runner/). ACs live alongside the features they verify, carry their own lifecycle, and compose into [test scenarios](https://github.com/synchestra-io/rehearse/blob/main/spec/features/testing-framework/test-scenario/) for end-to-end validation.

The full specification for this feature — file format, supported languages, identification scheme, statuses, and validation rules — lives in the [synchestra-io/rehearse](https://github.com/synchestra-io/rehearse/blob/main/spec/features/acceptance-criteria/) repository.

## SpecScore Conventions

SpecScore extends the base AC specification with project-specific conventions:

### Mandatory AC section in feature READMEs

Every SpecScore feature README must include an **Acceptance Criteria** section. The "Not defined yet." state triggers a mandatory Outstanding Question: "Acceptance criteria not yet defined for this feature." This ensures missing ACs are visible — not forgotten.

### Relationship to development plan ACs

Feature ACs and plan ACs serve different audiences and have different lifecycles:

| AC type | Lives in | Answers | Lifecycle |
|---|---|---|---|
| **Feature AC** | `spec/features/{feature}/_acs/` | "Does this feature work correctly?" | Evolves with the feature; long-lived |
| **Plan-level AC** | `spec/plans/{plan}/README.md` (inline or `_acs/` subdir) | "Were this plan's goals achieved?" | Frozen with the plan; immutable |
| **Plan step-level AC** | Within each plan step | "Was this step's deliverable produced?" | Frozen with the plan; immutable |

Plan step ACs may *reference* feature ACs — for example, "the feature AC `cli/project/remove/not-in-list` must pass after this step." But they are not the same artifact. Feature ACs are the long-lived, canonical verification units. Plan ACs are scoped to a single implementation effort and frozen on approval.

When execution units are derived from a plan, both plan step ACs and any referenced feature ACs are included in the execution context. Implementers know exactly what "done" looks like before they write a line of code.

### Outstanding Questions linkage

If the AC section says "Not defined yet.", the Outstanding Questions section must include the corresponding question. This keeps missing ACs visible until addressed.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Features gain a mandatory Acceptance Criteria section and `_acs/` directory convention. The feature spec defines the structural rules; this feature defines what goes inside. |
| [Development Plan](../development-plan/README.md) | Plan step ACs may reference feature ACs. Plan-level ACs follow the same format but are frozen with the plan. |
| [Testing Framework](../testing-framework/README.md) | Test scenarios reference ACs via table syntax. The test runner resolves and executes verification scripts. |
| [Outstanding Questions](../outstanding-questions/README.md) | Missing ACs surface as outstanding questions, keeping them visible until addressed. |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
