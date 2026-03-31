# Feature: Acceptance Criteria

**Status:** Conceptual

## Summary

Acceptance criteria (ACs) are optional inline sections within a feature README that bundle related [requirements](../requirement/README.md) into composite verification conditions. An AC states what must be true for a group of requirements to be satisfied — without prescribing specific inputs, flows, or implementation details. Concrete proof is the [scenario](../scenario/README.md)'s job.

ACs answer: **"What must be true?"** Scenarios answer: **"How do we prove it?"**

## Contents

| Directory | Description |
|---|---|
| [_tests/](_tests/README.md) | Test scenarios for acceptance criteria behavior |

### _tests

Scenarios that validate the requirements defined in this feature, covering AC formatting, requirement linkage, placeholder behavior, and the relationship between ACs and scenarios.

## Problem

Requirements (`#### REQ:`) capture individual enforceable rules, but reviewers and product owners often need to verify conditions that span multiple requirements. Without a grouping mechanism:

- Verification checklists repeat requirement references redundantly across scenarios.
- Reviewers cannot quickly see which requirements are verified together.
- There is no single place that states the composite condition a group of requirements must satisfy.

Acceptance criteria solve this by bundling related requirements under a named, addressable heading that scenarios can reference.

## Behavior

### AC format and location

ACs are inline sections within the feature README's `## Acceptance Criteria` section. They do not live in a separate directory.

#### REQ: inline-heading

Each AC MUST use the heading format `### AC: {slug}` within the `## Acceptance Criteria` section. The slug MUST be lowercase, hyphen-separated, and URL-safe.

#### REQ: requirements-metadata

Each AC MUST include a `**Requirements:**` field listing the requirements it bundles. The field appears on the line immediately after the `### AC:` heading. Requirements use the format `{feature-id}#req:{slug}`, comma-separated when multiple.

#### REQ: condition-statement

Each AC MUST include a prose statement (one or more sentences) after the `**Requirements:**` field that describes the composite verification condition. The statement is abstract — it defines what must be true, not how to test it.

### AC section in feature READMEs

Every feature README carries an Acceptance Criteria section. ACs are optional — the section may contain ACs, or it may indicate that none are defined yet.

#### REQ: section-required

Every feature README MUST include an `## Acceptance Criteria` section. This section is never omitted.

#### REQ: placeholder-when-empty

When no ACs are defined, the Acceptance Criteria section MUST state "Not defined yet." as its only content.

#### REQ: outstanding-question-linkage

When the Acceptance Criteria section reads "Not defined yet.", the feature's Outstanding Questions section MUST include "Acceptance criteria not yet defined for this feature." This keeps missing ACs visible until addressed.

### When to use ACs

ACs are a grouping mechanism, not a mandatory layer. Not every requirement needs an AC wrapper.

#### REQ: optional-grouping

ACs are OPTIONAL. They SHOULD only be created when bundling two or more related requirements adds clarity for reviewers or scenario authors. A requirement that stands alone SHOULD be referenced directly by its scenario without an intermediary AC.

#### REQ: no-duplicate-conditions

An AC MUST NOT restate a single requirement verbatim. If an AC would contain exactly one requirement and add no additional condition, the AC SHOULD be omitted and the scenario SHOULD reference the requirement directly.

### ACs vs scenarios

ACs and scenarios are complementary but distinct artifacts.

#### REQ: abstract-not-concrete

ACs MUST remain abstract verification conditions. They MUST NOT contain Given/When/Then steps, specific test inputs, or implementation details. Concrete flows belong in [scenarios](../scenario/README.md).

#### REQ: scenario-validates-ac

Scenarios that verify a bundled condition SHOULD reference the AC in their `**Validates:**` field using the format `{feature-id}#ac:{slug}`. When a scenario also needs to reference the underlying requirements, it MAY list both the AC and the individual REQs.

### Relationship to plan ACs

Feature ACs and development plan ACs serve different audiences and lifecycles.

#### REQ: feature-ac-long-lived

Feature ACs are long-lived and evolve with the feature specification. They MUST NOT be frozen or versioned with a development plan.

#### REQ: plan-ac-references

Development plan step ACs MAY reference feature ACs to indicate that a plan step's completion depends on the feature AC passing. The plan step AC is frozen with the plan; the referenced feature AC continues to evolve independently.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Features carry a mandatory Acceptance Criteria section. This feature defines what goes inside that section. |
| [Requirement](../requirement/README.md) | ACs bundle requirements via the `**Requirements:**` metadata field, creating traceability from composite verification conditions back to individual behavioral rules. |
| [Scenario](../scenario/README.md) | Scenarios validate ACs (or REQs directly) with concrete Given/When/Then flows. An AC is abstract; a scenario is its executable proof. |
| [Development Plan](../development-plan/README.md) | Plan step ACs may reference feature ACs. Plan-level ACs follow the same format but are frozen with the plan. |
| [Outstanding Questions](../outstanding-questions/README.md) | Missing ACs surface as outstanding questions, keeping them visible until addressed. |

## Acceptance Criteria

### AC: well-formed-ac

**Requirements:** acceptance-criteria#req:inline-heading, acceptance-criteria#req:requirements-metadata, acceptance-criteria#req:condition-statement

An AC uses the correct heading format, includes a Requirements metadata field linking to the REQs it bundles, and provides an abstract prose condition. An AC missing any of these elements fails validation.

### AC: empty-section-handling

**Requirements:** acceptance-criteria#req:section-required, acceptance-criteria#req:placeholder-when-empty, acceptance-criteria#req:outstanding-question-linkage

The Acceptance Criteria section is always present. When it has no ACs, it reads "Not defined yet." and a corresponding outstanding question is raised. Once ACs are added, the placeholder and question are removed.

### AC: ac-scenario-separation

**Requirements:** acceptance-criteria#req:abstract-not-concrete, acceptance-criteria#req:scenario-validates-ac, acceptance-criteria#req:optional-grouping

ACs remain abstract — they never contain Given/When/Then steps or test data. Scenarios reference ACs (or REQs directly) in their Validates field. ACs are only created when bundling adds value; standalone requirements are referenced directly.

## Outstanding Questions

- Should tooling enforce a minimum of two requirements per AC, or is the "optional grouping" convention sufficient without automated enforcement?
- How should ACs handle cross-feature requirement bundles — e.g., an AC in feature A that references requirements from both feature A and feature B?
