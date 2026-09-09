---
kind: entity
id: plan
singular: Plan
plural: Plans
description: A composite task — a task with subtasks — that bridges Feature specifications to executable work.
properties:
  - name: id
    data_type: string
    description: The Plan's full slash-separated logical ID relative to its source namespace.
    checks:
      required: true
      min_length: 1
      max_length: 128
      pattern: "^[a-z0-9]+(?:-[a-z0-9]+)*(?:/[a-z0-9]+(?:-[a-z0-9]+)*)*$"
  - name: status
    data_type: string
    description: The plan's full-lifecycle status in one field, across three bands — prep (human-authored), execution (derived by `lint --fix` from task rollup), and disposition (human-authored).
    checks:
      required: true
      enum:
        - Draft
        - In Review
        - Approved
        - Executing
        - Blocked
        - Implemented
        - Failed
        - Rejected
        - Withdrawn
        - Superseded
        - Deprecated
  - name: features
    data_type: array
    description: Features this Plan affects (its source Feature for feature-sourced plans). Empty for idea-sourced plans, which bind to an Idea via the `Source` line instead.
    checks:
      required: false
      items:
        data_type: ref
        checks:
          entity_ref: ../feature/feature.entity.md
  - name: tasks
    data_type: array
    description: Subtasks that make up this Plan. Empty list means the Plan is a leaf Task.
    checks:
      required: false
      items:
        data_type: ref
        checks:
          entity_ref: ../task/task.entity.md
  - name: tasks_count
    data_type: integer
    description: Derived count of the Plan's direct child tasks. Maintained by `specscore spec lint --fix` and surfaced in frontmatter (per artifact-frontmatter-convention); never hand-authored.
    checks:
      required: false
      min: 0
  - name: parent
    data_type: string
    description: The master Plan this Plan is a sub-plan of — the master/sub-plan composition link. A same-source full logical ID, or a cross-repository `<repo-slug>:<plan-logical-id>` soft reference. Absent for root plans. Single-parent (a tree) in the MVP; multi-parent DAGs are out of scope. Same-repo parents are resolved and checked for acyclicity by lint rule `P-005`; cross-repo parents are validated syntactically only (no sibling-repo resolution).
    checks:
      required: false
      max_length: 256
format: https://specscore.md/entity-specification
---

# Entity: Plan

## Description

The Plan entity is the typed shape of a SpecScore
[Plan](README.md) — the **how** layer that decomposes one or more
[Feature entities](../feature/feature.entity.md) into ordered
[Task entities](../task/task.entity.md).

A Plan is a composite Task; structurally, the two share their core
shape. The distinction is whether `tasks` is empty (leaf Task) or
non-empty (Plan). See the Plan spec's
[unified plan/task model](README.md#recursive-task-and-plan-model) for the full
recursive definition.

## Properties

<!-- managed-by: specscore lint --fix -->
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | yes | The Plan's full slash-separated logical ID relative to its source namespace. |
| `status` | string | yes | The plan's full-lifecycle status in one field, across three bands — prep (human-authored), execution (derived by `lint --fix` from task rollup), and disposition (human-authored). |
| `features` | array | no | Features this Plan affects (its source Feature for feature-sourced plans). Empty for idea-sourced plans, which bind to an Idea via the `Source` line instead. |
| `tasks` | array | no | Subtasks that make up this Plan. Empty list means the Plan is a leaf Task. |
| `tasks_count` | integer | no | Derived count of the Plan's direct child tasks. Maintained by `specscore spec lint --fix` and surfaced in frontmatter (per artifact-frontmatter-convention); never hand-authored. |
| `parent` | string | no | The master Plan this Plan is a sub-plan of — the master/sub-plan composition link. A same-source full logical ID, or a cross-repository `<repo-slug>:<plan-logical-id>` soft reference. Absent for root plans. Single-parent (a tree) in the MVP; multi-parent DAGs are out of scope. Same-repo parents are resolved and checked for acyclicity by lint rule `P-005`; cross-repo parents are validated syntactically only (no sibling-repo resolution). |
<!-- end-managed -->

## Referenced by

<!-- managed-by: specscore lint --fix -->
- _No references yet._
<!-- end-managed -->

---
*This document follows the https://specscore.md/entity-specification*
