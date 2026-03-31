# Feature: Acceptance Criteria

**Status:** Conceptual

## Summary

Acceptance criteria are **abstract verification conditions** — success/failure statements that define what must be true for a requirement to be satisfied. Each AC is a standalone markdown file that states a condition without prescribing specific inputs, flows, or implementation details. ACs are readable by product owners, auditable by reviewers, and verifiable through [scenarios](../scenario/README.md) that provide concrete Given/When/Then proof.

ACs answer: **"What must be true?"** — not **"How do we test it?"** That is the [scenario](../scenario/README.md)'s job.

The full specification for this feature — file format, supported languages, identification scheme, statuses, and validation rules — lives in the [synchestra-io/rehearse](https://github.com/synchestra-io/rehearse/blob/main/spec/features/acceptance-criteria/) repository.

## SpecScore Conventions

SpecScore extends the base AC specification with project-specific conventions:

### Mandatory AC section in feature READMEs

Every SpecScore feature README must include an **Acceptance Criteria** section. The "Not defined yet." state triggers a mandatory Outstanding Question: "Acceptance criteria not yet defined for this feature." This ensures missing ACs are visible — not forgotten.

### Requirement traceability

Each AC SHOULD include a `**Requirement:**` metadata field linking it to the requirement(s) it verifies:

```markdown
# AC: title-required

**Requirement:** todo-item/manage#req:title-required

Creating a todo without a title is rejected. Creating a todo with a title succeeds.
```

The reference format is `{feature-path}#req:{slug}`. An AC may reference multiple requirements as a comma-separated list when it verifies a condition spanning them.

### Acceptance criteria vs scenarios

ACs and [scenarios](../scenario/README.md) are complementary but distinct:

| | Acceptance Criteria | Scenario |
|---|---|---|
| **Abstraction** | Abstract condition | Concrete flow |
| **Format** | Prose statement | Given/When/Then steps |
| **Location** | `_acs/` directory | `_tests/` directory |
| **Purpose** | Define what must be true | Prove it with specific inputs |
| **Executable** | No (verified through scenarios) | Yes (via Rehearse) |

**Example pair:**

AC (`_acs/title-required.md`):
> Creating a todo without a title is rejected. Creating a todo with a title succeeds.

Scenario (`_tests/create-without-title.md`):
> GIVEN an empty todo list
> WHEN the user runs `todo add ""`
> THEN the CLI prints "Error: title is required"
> AND the list remains empty

The AC states the rule; the scenario proves it with a concrete example.

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
| [Requirement](../requirement/README.md) | ACs reference requirements via the `**Requirement:**` metadata field, creating traceability from verification conditions back to specific behavioral rules. |
| [Scenario](../scenario/README.md) | Scenarios validate ACs with concrete Given/When/Then flows. An AC is abstract; a scenario is its executable proof. |
| [Development Plan](../development-plan/README.md) | Plan step ACs may reference feature ACs. Plan-level ACs follow the same format but are frozen with the plan. |
| [Testing Framework](../testing-framework/README.md) | Test scenarios reference ACs via table syntax. The test runner resolves and executes verification scripts. |
| [Outstanding Questions](../outstanding-questions/README.md) | Missing ACs surface as outstanding questions, keeping them visible until addressed. |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
