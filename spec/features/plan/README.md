---
format: https://specscore.md/feature-specification
status: Stable
---

# Feature: Plan

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/plan?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/plan?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/plan?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/plan?op=request-change) |

**Status:** Stable
**Source Ideas:** —

## Summary

A plan is a composite task -- a task that contains subtasks. It bridges feature specifications and change requests to executable work. Plans are mutable documents; snapshots provide immutable reference points for review, approval, and retrospective.

There is one structural concept: the **task**. A task with children is a plan. A task without children is a leaf task. Every Plan has a canonical directory whose `README.md` is the Plan document; child directories express recursive sub-plans with no artificial depth limit. The typed shape of a single Plan is captured in the co-located [plan entity](plan.entity.md).

## Contents

| Directory | Description |
|---|---|
| [_tests](_tests/README.md) | Test scenarios validating plan feature requirements |

The Plan document type has a companion Index-Kind feature, [plans-index](../plans-index/README.md), which specifies the shape of the `spec/plans/README.md` aggregation file every repo maintains. Plans-Index lives as a top-level sibling of Plan under `spec/features/` — not as a sub-feature — so its specification URL stays flat (`plans-index-specification`, not a nested compound).

## Problem

Teams have well-defined execution systems -- work items, tickets, tasks -- but there is no structured way to go from "we know what to build" to "here are the work items to execute."

Today that decomposition happens ad hoc -- a human or AI agent reads a feature spec, mentally breaks it into tasks, and manually creates work items one by one. This creates three problems:

- **No review gate.** Work begins without explicit approval of the approach. A bad decomposition wastes time and effort.
- **No stable reference.** Work items are designed to be fluid -- agents add subtasks, humans cancel items, parallel work gets restructured. This fluidity is a feature of execution, but it means there is no fixed record of what was originally planned.
- **No retrospective anchor.** Without a snapshot of intent, you cannot compare what was planned against what actually happened. Lessons learned require a before-and-after.

## Design Philosophy

SpecScore separates **intent** from **execution** by design, with distinct artifacts for each stage of the workflow.

| Artifact | Question it answers | Audience | Mutability | Lives in |
|---|---|---|---|---|
| Feature spec | What do we want? | Product, engineering | Versioned | Spec repo |
| Change request | What should change in an existing feature? | Product, engineering | Versioned until approved | Spec repo |
| Plan | How will we build it? | Reviewers, planners | Mutable; snapshots provide fixed references | Spec repo |
| Work items | Who's doing what right now? | Agents, operators | Highly fluid | Execution system |

A **feature spec** defines something new. A **change request** (implemented as a [proposal](../proposals/README.md)) mutates something that already exists. Both are *what* artifacts -- they describe desired outcomes. The distinction matters because:

- **New features** start from a blank slate. The plan is unconstrained.
- **Change requests** operate on existing behavior. The plan must account for what is already there -- migration paths, backward compatibility, affected dependents. The review process is different: reviewers need to understand the delta, not just the destination.

From the planning pipeline's perspective, both converge to the same output -- a plan that produces executable tasks:

```mermaid
graph LR
    A["Feature spec"]
    B["Plan"]
    C["Tasks / Work items"]
    D["Change request<br/>(proposal)"]

    A --> B
    D --> B
    B --> C
```

**Why not use the work item tree as the plan?** Work items are designed to be fluid. Agents add subtasks when they discover complexity. Humans cancel items when priorities shift. Parallel work gets restructured on the fly. This fluidity is a feature -- it is how real development works. But fluidity is the enemy of reviewability. A human reviewer needs a stable, scannable document to approve before work begins. Snapshots provide that stability without sacrificing the ability to evolve the plan.

**No duplicated status tracking.** The plan does not track completion -- execution tools do. A progress view can be derived by mapping plan tasks to their linked work items and looking up live status. One source of truth, two views: the plan view for humans, the deep work item tree for agents.

## Behavior

### Plan location

A Plan is stored in a slug-named directory whose `README.md` is the document:

```text
spec/plans/
  README.md
  {plan-slug}/
    README.md
    {child-plan-slug}/
      README.md
```

`{plan-slug}` is a URL/path-safe identifier (e.g., `add-batch-mode`, `user-auth`). A nested Plan's logical ID is its full slash-separated path relative to its source namespace, such as `chat-feature/chat-infrastructure/database-setup`.

#### REQ: plan-file

Every newly created Plan MUST be a `README.md` in its canonical directory. A root Plan is `spec/plans/{plan-slug}/README.md`; each child Plan adds one slug directory below its parent. The plans index alone owns `spec/plans/README.md`. Readers MAY accept a legacy flat `{plan-slug}.md` during migration, but it MUST NOT coexist with a directory Plan of the same logical ID.

#### REQ: plan-logical-id

A Plan's logical ID MUST preserve every path segment from the namespace root. Creation, lookup, dependency resolution, status changes, and index generation MUST use the full logical ID; two leaf names under different parents are distinct Plans.

### External Plan namespace

When a source project routes Plans to another repository, the destination stores them under the source project's complete repository identity:

```text
spec/plans/{source-host}/{source-owner}/{source-repo}/{nested-plan-id}/README.md
```

For example, `github.com/datatug/datatug` Plan `phase-1/core-loop` lives at `spec/plans/github.com/datatug/datatug/phase-1/core-loop/README.md` in the destination. The host segment prevents repositories with the same owner/name on different forges from sharing a namespace. When the explicitly selected destination equals the source repository, the existing `spec/plans/{nested-plan-id}/README.md` layout is used without the identity prefix.

#### REQ: external-source-namespace

An external destination MUST isolate every source at `spec/plans/{source-host}/{source-owner}/{source-repo}/`. Repository identities and Plan IDs MUST be joined as validated path segments. Neither lexical traversal nor symlinks in an existing namespace may escape the destination repository or cross into another source namespace.

#### REQ: source-project-context

`--project` always identifies the source project, including when invoked from a generic Plan-store checkout. The source project owns Features, acceptance criteria, Ideas, and other specification context. Plan operations MAY read those source artifacts for validation but MUST NOT modify them; all Plan and Plan-index writes go only to the resolved Plan namespace.

#### REQ: read-only-plan-operations

Read-only and dry-run Plan operations MUST leave both source and destination repositories byte-for-byte unchanged. Validation fixes required for source artifacts MUST be reported rather than applied as a side effect of a Plan operation.

#### REQ: plan-path-conflicts

A logical Plan ID MUST have exactly one on-disk representation. If a legacy flat file and canonical directory form resolve to the same ID, validation and mutation fail before durable writes. `--force` MAY replace an existing canonical Plan document but MUST NOT choose between ambiguous representations or overwrite a parent Plan.

#### REQ: external-store-lifecycle-lock

An external Plan-store repository MUST ignore the anchored `/.specscore-lifecycle.lock` file. The persistent lock inode is shared lifecycle coordination state and MUST NOT be deleted after each operation. A destination that has not otherwise been initialized by SpecScore still requires this exact gitignore rule before Plan mutation.

#### REQ: plan-slug-format

Plan slugs MUST be lowercase, hyphen-separated, and URL-safe. Underscores, spaces, and special characters MUST NOT be used.

#### REQ: source-binding

Every plan MUST declare exactly one source line in its header, in one of three forms: a `**Source Feature:**` line naming the Feature it decomposes (a Feature-sourced plan), a `**Source:** idea:{slug}` line naming the Idea it plans directly (an Idea-sourced plan), or a `**Source:** none` line (a source-less plan, not tied to any Feature or Idea). A plan that declares more than one source, or whose source line is absent or carries an unrecognized value, is invalid.

### Plan document structure

```markdown
# Plan: Add batch mode to CLI

**Status:** Approved
**Source Feature:** cli
**Date:** 2026-03-14
**Owner:** @alex
**Supersedes:** —

## Summary

1-3 sentences. What this plan covers and how it decomposes the source
Feature (or, for an idea-sourced plan, the source Idea).

## Approach

<=1 paragraph on the decomposition strategy: what was grouped, what was
deferred and why.

## Tasks

### Task 1: Define batch input schema

**Verifies:** cli#ac:batch-schema-validates

Establish the YAML/JSON schema for batch input files. This determines
the contract for all downstream tasks.

### Task 2: Implement batch parser

**Verifies:** cli#ac:batch-parser-rejects-invalid

Parse and validate batch input files against the schema from Task 1,
rejecting invalid files with per-field error messages.

### Task 3: Update CLI entry point

**Verifies:** cli#ac:batch-flag-in-help

Add the `--batch <file>` flag and wire it to the parser; surface it in
help output, mutually exclusive with positional arguments.

## Deferred AC Coverage

<!-- Omit this section entirely when no ACs are deferred. -->

## Open Questions

None at this time.

---
*This document follows the https://specscore.md/plan-specification*
```

#### REQ: plan-title-format

Every plan document MUST use the `# Plan: {Title}` format for its title. The `Plan:` prefix is required.

#### REQ: plan-required-sections

Every plan document MUST include the following sections: title (`# Plan: X`), header metadata fields, `## Summary`, `## Approach`, `## Tasks`, and `## Open Questions`. A `## Deferred AC Coverage` section, `## Snapshots` section, and `## Dependency graph` section are OPTIONAL.

### Header fields

| Field | Required | Description |
|---|---|---|
| **Status** | Yes | Current plan status (see [Plan statuses](#plan-statuses)) |
| **Source Feature** | One of | Slug of the Feature this plan decomposes (feature-sourced plans) |
| **Source** | One of | `idea:{slug}` naming the Idea this plan decomposes (idea-sourced plans), or `none` for a source-less plan |
| **Date** | Yes | Date the plan was created |
| **Owner** | Yes | Who wrote the plan |
| **Supersedes** | Yes | `—`, or the slug of an older plan this one wholesale-replaces |
| **Parent** | No | The master plan this plan is a sub-plan of (master/sub-plan composition). A same-source full Plan logical ID or a cross-repository `<repo-slug>:<plan-logical-id>` soft reference. Absent for root plans. See [Cross-repo plan composition](#cross-repo-plan-composition). Validated by lint rule `P-005`. |
| **Prerequisite Plans** | No | Comma-separated full logical IDs of same-source Plans that must be implemented before this Plan begins. Nested IDs retain every ancestor segment. |
| **Coordination** | No | The repo/branch where this plan document's own mutations are authoritative -- `<owner>/<repo>@<branch>`. Absent means unrestricted (any repo/branch may mutate the plan, as before this field existed). See [Coordination branch](#coordination-branch). Validated by lint rule `P-010`. |
| **Effort** | No | `S` \| `M` \| `L` \| `XL` -- see [Optional ROI metadata](#optional-roi-metadata) |
| **Impact** | No | `low` \| `medium` \| `high` \| `critical` -- see [Optional ROI metadata](#optional-roi-metadata) |

#### REQ: required-header-fields

Every plan MUST include these header fields: Status, the source line (exactly one of `Source Feature`, `Source: idea:{slug}`, or `Source: none`, per [source-binding](#req-source-binding)), Date, Owner, and Supersedes. Effort and Impact are OPTIONAL.

#### REQ: proposal-forward-reference

When a plan is triggered by a change request (proposal), the Source field MUST link directly to the proposal. The proposal in turn MUST include a forward reference to the plan.

When a plan is triggered by a change request (proposal), the **Source** field links directly to the proposal. The proposal in turn gets a forward reference to the plan:

```markdown
# Proposal: Deprecate v1 endpoints

| Field  | Value                                             |
|--------|---------------------------------------------------|
| Status | `approved`                                        |
| Plan   | [migrate-to-v2](../../../plans/migrate-to-v2.md)  |
```

### Plan statuses

A plan's status models its full lifecycle in one field, in three bands. The values are capitalized, matching the SpecScore-wide status vocabulary used by Features and Ideas; the frontmatter `status:` mirror (per [artifact-frontmatter-convention](../artifact-frontmatter-convention/README.md)) carries the same value verbatim.

| Band | Status | Description | Set by |
|---|---|---|---|
| Prep | `Draft` | Plan is being written, not ready for review | Human author |
| Prep | `In Review` | Submitted for review | Human author |
| Prep | `Approved` | Reviewed and approved — ready/pending execution | Human author |
| Execution | `Executing` | At least one task in progress | Derived by `lint --fix` |
| Execution | `Blocked` | Tasks blocked; none progressing and none failed | Derived by `lint --fix` |
| Execution | `Implemented` | All tasks complete | Derived by `lint --fix` |
| Execution | `Failed` | A task failed/aborted and the plan cannot complete | Derived by `lint --fix` |
| Disposition | `Rejected` | Approach rejected outright at review (not sent back for revisions) | Human author |
| Disposition | `Withdrawn` | Abandoned | Human author |
| Disposition | `Superseded` | Replaced by a named successor plan | Human author |
| Disposition | `Deprecated` | Approach no longer recommended, with no named successor | Human author |

The bands are sequential phases of one lifecycle, not concurrent: reaching `Executing` already implies the approval gate passed (or was deliberately bypassed), so `Approved` carries no additional current-state information during a run. The authority handoff sits at `Approved`: `lint --fix` only ever transitions from `Approved` onward (deriving the execution band from task-status rollup) and MUST NEVER overwrite a human-authored prep state. Plans are mutable; if the approach changes after approval, edit the plan and record a new snapshot rather than creating a separate document.

#### REQ: valid-statuses

A plan's Status field MUST be one of: `Draft`, `In Review`, `Approved`, `Executing`, `Blocked`, `Implemented`, `Failed`, `Rejected`, `Withdrawn`, `Superseded`, or `Deprecated`. No other values are permitted. A `Superseded` plan MUST carry a reference to its successor plan.

#### REQ: execution-status-derived

The execution-band statuses (`Executing`, `Blocked`, `Implemented`, `Failed`) MUST NOT be hand-authored. They are derived by `specscore spec lint --fix` from the rollup of the plan's task statuses (see [Status rollup](#status-rollup)), and `lint --fix` transitions only from `Approved` onward, never overwriting a `Draft`/`In Review`/`Approved` prep state.

### Status transitions

```mermaid
graph LR
    A["Draft"]
    B["In Review"]
    C["Approved"]
    E["Executing"]
    K["Blocked"]
    I["Implemented"]
    F["Failed"]
    R["Rejected"]
    W["Withdrawn"]
    S["Superseded"]
    D["Deprecated"]

    A -->|submit| B
    A -->|approve| C
    B -->|revisions| A
    B -->|approve| C
    B -->|reject| R
    C -->|lint --fix rollup| E
    E --> K
    E --> I
    E --> F
    K --> E
    C -->|abandon| W
    C -->|replace| S
    C -->|retire| D
```

#### REQ: status-transitions

Plan status transitions MUST follow these rules. **Prep (human-authored):** `Draft` MAY transition to `In Review`, or directly to `Approved` (a fast-track that skips review); `In Review` MAY transition back to `Draft` (revisions requested), forward to `Approved`, or to `Rejected` (the approach is rejected outright, not sent back for revisions). **Execution (derived by `lint --fix`, only from `Approved` onward):** `Approved` MAY transition to an execution-band status; execution statuses transition among themselves per the task-status rollup. **Disposition (human-authored):** `Approved` (or a later state) MAY transition to `Withdrawn`, `Superseded`, or `Deprecated` (the approach is no longer recommended, with no named successor). There is no resurrection from a disposition status — re-pursuing the work means authoring a new plan. No other transitions are permitted.

### Snapshots

A snapshot is an immutable reference point within a plan's history. Instead of freezing the entire plan on approval, snapshots record meaningful moments -- approval, checkpoints, completion -- as entries in a table with a corresponding git commit hash.

Snapshots remove the need to freeze a plan on approval. A plan can be edited freely at any time. When a reference point is needed, a snapshot captures the plan's state at that git hash. The snapshot table lives in the plan document:

```markdown
## Snapshots

| Date | Git Hash | Action | Comment |
|---|---|---|---|
| 2026-03-15 | `a1b2c3d` | approved | Initial approval by @jordan |
| 2026-03-20 | `e4f5g6h` | checkpoint | Added streaming support task |
| 2026-04-01 | `i7j8k9l` | completed | All tasks verified |
```

#### REQ: snapshot-table-format

When a plan includes snapshots, they MUST be recorded in a `## Snapshots` section containing a table with columns: Date, Git Hash, Action, and Comment.

#### REQ: snapshot-actions

Snapshot actions include `approved`, `checkpoint`, `completed`, and user-defined values. The `approved` action SHOULD correspond to setting the plan status to `Approved`.

#### REQ: snapshot-git-hash

Each snapshot MUST reference a valid git commit hash that represents the plan's state at the time the snapshot was taken.

### Recursive task and plan model

The directory hierarchy is the canonical Plan model. Inline `### Task N:` blocks remain leaf execution tasks inside a Plan document; child Plan directories decompose larger units while preserving one document format.

A plan is a composite task -- it contains other tasks. Some of those child tasks may themselves contain subtasks, making them sub-plans. This nesting is recursive with no artificial depth limit.

```text
spec/plans/
  README.md                          <- index
  chat-feature/
    README.md                        <- plan (composite task)
    chat-infrastructure/
      README.md                      <- sub-plan (also a composite task)
      set-up-database/
        README.md                    <- task (leaf)
      configure-networking/
        README.md                    <- task (leaf)
    chat-workflow-engine/
      README.md                      <- sub-plan
    send-notifications/
      README.md                      <- task (leaf, sibling of sub-plans)
  e2e-testing-framework/
    README.md                        <- plan (standalone, no sub-plans)
```

Whether something is a "plan" or a "task" is determined by structure: if it has children, it is a plan. If it has no children, it is a task. The same document format applies at every level.

#### REQ: recursive-nesting

Plans MAY nest to arbitrary depth. There is no maximum nesting level. Depth is a judgment call made by the plan author, and each nested document keeps its full logical ID. Inline `### Task N:` blocks MAY coexist with child Plan directories.

#### REQ: child-plan-format

A sub-plan (child plan) MUST follow the same format as a top-level plan, including tasks and acceptance criteria.

### Mixed children

Tasks and sub-plans can coexist at the same level within a plan. There is no requirement to separate them or force restructuring. In the example above, `send-notifications` (a leaf task) is a sibling of `chat-infrastructure` and `chat-workflow-engine` (sub-plans).

#### REQ: mixed-children

A plan MAY contain both leaf tasks and sub-plans as direct children at the same level. No restructuring is required to separate tasks from sub-plans.

### Cross-repo plan composition

The directory hierarchy nests sub-plans physically inside one source namespace. **Cross-repo plan composition** composes Plan directories in different source namespaces by reference rather than by physical nesting. It is the mechanism behind a single Idea fanning out into coordinated work across repos.

A **master plan** is an ordinary Plan whose work is carried out by **sub-plans**. Each sub-plan names its master through the **Parent** header field (`**Parent:** <plan-ref>`). The reference is either a same-source full Plan logical ID or a cross-repository `<repo-slug>:<plan-logical-id>` soft reference. On the master side, a task MAY delegate to a sub-plan via the Task entity's `sub_plan` property — the master-side expression of the same edge.

```text
specscore/spec/plans/cross-repo-master/README.md        <- master plan (root: no Parent)
specscore-cli/spec/plans/sub-cli-bootstrap/README.md    <- **Parent:** specscore:cross-repo-master
specscore/spec/plans/sub-entity-model/README.md         <- **Parent:** cross-repo-master   (same-repo)
specstudio-skills/spec/plans/sub-skills/README.md       <- **Parent:** specscore:cross-repo-master
```

The canonical, navigable link is the child's **Parent** ref (child → master); composition is single-parent (a tree) in the MVP. Lint owns reference integrity per rule `P-005`: same-repo parents are resolved and checked for acyclicity; cross-repository `<repo-slug>:<plan-logical-id>` parents are validated **syntactically only** — the linter never scans sibling repositories, so a cross-repo parent is a best-effort, unresolved reference (like an external link). Cross-repo back-link maintenance and execution ordering across the tree are out of scope for this model and belong to the consuming skills.

#### REQ: cross-repo-parent-ref

A plan MAY declare a single `**Parent:** <plan-ref>` header field naming the master plan it is a sub-plan of. The value is a same-source full Plan logical ID or a cross-repository `<repo-slug>:<plan-logical-id>` reference. Nested references MUST retain every ancestor segment. A plan with no `**Parent:**` field is a root plan. Composition is single-parent in the MVP: a plan MUST NOT declare more than one parent. Reference validity (same-repo resolution and acyclicity; cross-repo syntactic-only checks) is enforced by `specscore spec lint` rule `P-005`, not by this document.

#### REQ: prerequisite-plan-logical-ids

Each `**Prerequisite Plans:**` entry MUST be a full same-source Plan logical ID. Readiness resolves the entry relative to the resolved source namespace, preserves arbitrary nesting, and reports missing, malformed, cyclic, or not-yet-implemented prerequisites without falling back to a matching leaf name elsewhere in the hierarchy.

### Coordination branch

A plan document is a single file that many parallel actors -- human or AI agent, each often on their own feature branch or in their own git worktree -- read and want to move forward at the same time. Without an agreed rule for *where* the plan document itself is mutated, two actors editing the same plan on different branches (adding tasks, walking a status) collide when their branches merge, and the collision has to be untangled by hand (e.g. two independently numbered `### Task 16:` blocks landing from two different feature branches).

The **Coordination branch** is that agreed rule, declared once in the plan's own header:

```markdown
**Coordination:** specscore/specscore-cli@main
```

`<owner>/<repo>` is a GitHub owner/repository pair; `<branch>` is the branch name (which MAY itself contain `/`, e.g. `feature/plan-coordination-branch`). When present, every mutation to *this plan document* -- task status walks, new tasks, task renumbering, plan-status transitions, or any other edit to the file -- is authoritative only on that declared repo/branch. A plan with no `**Coordination:**` field carries no such restriction: it may be mutated from anywhere, exactly as before this field existed.

The field draws a line between two kinds of truth that otherwise get conflated on a feature branch:

- **Feature branches carry implementation.** Code changes, task completion evidence, everything a plan's tasks describe building.
- **The coordination branch carries plan truth.** The plan document's own structure and status.

An actor implementing a task on `feature/some-task` never needs to also own the authoritative copy of the plan file -- their task-completion report is applied on the coordination branch (directly, or via whatever routing the consuming tooling uses), rather than mutating a parallel copy of the plan that later has to be reconciled by hand.

A **merge conflict in a plan file that declares `**Coordination:**`** is therefore a signal, not a chore: it means the separation was violated -- some actor mutated the plan document on a branch other than the declared one. The fix is to redo the mutation on the coordination branch (or to fix the tooling that let it happen elsewhere), not to hand-merge the conflicting hunks as if this were an ordinary content collision.

Enforcement of the declared repo/branch -- checking the current invocation against it, refusing a mismatched mutation, and the override mechanism for exceptional cases -- is a CLI concern (`specscore-cli`), not this specification: this document defines only the field's syntax and intended meaning. See `specscore-cli`'s plan `change-status`/`reconcile` and task `change-status` Features for the enforcement contract.

#### REQ: coordination-branch-format

A plan MAY declare a single `**Coordination:** <owner>/<repo>@<branch>` header field. `<owner>/<repo>` MUST be a GitHub owner/repository pair (no `/` within either component); `<branch>` MUST be a non-empty git ref name (which MAY contain `/`). A plan with no `**Coordination:**` field is unrestricted. The reference is validated **syntactically only** by `specscore spec lint` (lint rule `P-010`), mirroring the cross-repo precedent of [cross-repo-parent-ref](#req-cross-repo-parent-ref) and the [implementation-commit-provenance](../implementation-commit-provenance/README.md) ref format: the linter never resolves or scans the named repository, and never checks whether the branch exists.

#### REQ: coordination-branch-semantics

When a plan declares `**Coordination:**`, all mutations to that plan document -- task status transitions, new tasks, task renumbering, plan-status transitions, and any other edit to the file -- are authoritative only on the declared `<branch>` of the declared `<owner>/<repo>`. Feature branches (and worktrees derived from them) carry implementation; the coordination branch carries plan truth. A merge conflict in a plan file that declares `**Coordination:**` indicates this separation was violated (a mutation landed on a branch other than the declared one), and MUST be treated as a process violation to fix at the source -- not as an ordinary merge chore to resolve hunk-by-hunk.

### Status rollup

Once a plan is `Approved`, its execution-band status is derived from the rollup of its task statuses. This is what `lint --fix` maintains; it never runs while the plan is in a prep state (`Draft`/`In Review`/`Approved` are human-owned).

| Condition (task-status rollup) | Derived plan status |
|---|---|
| A task is `failed`/aborted and the plan cannot complete | `Failed` |
| At least one task is `in_progress` | `Executing` |
| Tasks are blocked; none in progress and none failed | `Blocked` |
| All tasks are complete | `Implemented` |

Precedence is top-to-bottom: `Failed` wins over `Executing`, which wins over `Blocked`, which wins over `Implemented`.

#### REQ: status-rollup

When a plan is `Approved` or in an execution-band status, `specscore spec lint --fix` MUST derive the plan's execution-band status from the rollup of its task statuses per the precedence above. The rollup reads task status only — it MUST NOT write task status, and MUST NOT overwrite a human-authored prep (`Draft`/`In Review`/`Approved`) or disposition (`Withdrawn`/`Superseded`) status.

### Task count

A plan carries a derived count of its tasks, surfaced in frontmatter so external tools can size a plan without parsing its body.

#### REQ: tasks-count

A plan MUST carry a derived `tasks_count` recording the number of its direct child tasks. The value is maintained by `specscore spec lint --fix` and surfaced in the plan's frontmatter per the [artifact-frontmatter-convention](../artifact-frontmatter-convention/README.md) feature. Authors MUST NOT hand-maintain `tasks_count`; lint reconciles it from the actual task children.

### Tasks and dependencies

Tasks whose `Depends on` is `none` may execute in parallel. The dependency graph determines the critical path.

For complex plans, an optional **Dependency graph** section visualizes the parallelism:

```mermaid
graph LR
    A["Task 1"]
    B["Task 2"]
    C["Task 3"]
    D["Task 4"]
    E["Task 5<br/>(independent)"]

    A --> B
    B --> C
    B --> D
```

This section is optional -- useful for complex plans, noise for simple sequential ones.

#### REQ: parallel-eligibility

Tasks whose `Depends on` value is `none` MUST be treated as parallel-eligible. Because `Depends on` is now required on every task (see [task#req:task-required-fields](../task/README.md#req-task-required-fields)), parallel-eligibility keys off the explicit `none` value rather than an absent field. The dependency graph determines the critical path; tasks with no dependencies MAY execute concurrently.

### Task-to-feature-AC traceability

A Plan does not embed its own acceptance-criteria sections. Instead, each task declares a `**Verifies:**` line naming one or more acceptance criteria of the source Feature (by `feature-slug#ac:<ac-slug>`) that the task implements. This is the plan's traceability mechanism — every task maps to the feature ACs it satisfies — and it is enforced by lint rule `P-001` (every plan task references at least one feature AC).

#### REQ: task-verifies-feature-ac

Every task in a plan MUST declare a `**Verifies:**` line referencing one or more acceptance criteria of the plan's source Feature, in the form `feature-slug#ac:<ac-slug>` (comma-separated when more than one). A task with no such reference is rejected by lint rule `P-001`. Feature ACs that the plan does not yet cover are recorded under the optional `## Deferred AC Coverage` section. (Idea-sourced plans, which have no source Feature, are exempt from `P-001`.)

### Optional task id

A task block MAY carry an optional `**Id:** <slug>` field — a stable, hyphen-separated identifier (matching the [Task entity](task.entity.md) `id` shape) used to address the inline task from tooling. Unlike the ordinal `Task N` (which shifts on reordering) or a title-derived slug (which breaks on title edits), the `**Id:**` is durable: it is the address `specscore task change-status --plan <slug> <id>` resolves against when stamping implementation-commit provenance on a plan-inline task. The field is optional and unenforced today; tooling that addresses plan-inline tasks requires it on the targeted block.

### Optional task note and evidence

A task block MAY also carry two independent, optional annotation fields, written by `specscore task change-status --note=<text> --evidence=<ref>[,<ref>...]` immediately after `**Status:**` (and after `**Implemented-by:**` when provenance is written in the same call):

```markdown
**Status:** complete
**Implemented-by:** sneat-co/chess@cfabf5e
**Note:** shipped to production, verified live
**Evidence:** cfabf5e, https://chessraiders.com/board/
```

- `**Note:**` is a free-text justification, unstructured and unvalidated.
- `**Evidence:**` is a comma-separated list of supporting references — commit SHAs, PR URLs, file paths, deploy or monitoring links — also unstructured and unvalidated.

Both are distinct from `**Implemented-by:**` (the [implementation-commit-provenance](../implementation-commit-provenance/README.md) feature's single, syntactically validated code reference): `**Implemented-by:**` answers "which commit did the work," while `**Note:**`/`**Evidence:**` answer "what backs the claim that it's actually done" — a broader category a strict commit-ref format cannot carry (e.g. a live-URL check or a manual QA note). Unlike the provenance flags (valid only with `--to=complete`), `--note`/`--evidence` are valid on **any** legal task transition. Neither field participates in the task-status transition matrix or the [status rollup](#status-rollup) — they are pure annotations.

### Optional ROI metadata

Two optional fields can be added to the plan document header:

```markdown
**Effort:** S | M | L | XL
**Impact:** low | medium | high | critical
```

Both fields are **optional**. When absent, tooling may infer effort from task count, dependency depth, and acceptance criteria complexity. It may infer impact from feature importance and downstream dependents. During plan authoring, the tooling **suggests** values. The user accepts, declines, or overwrites.

For composite plans, effort/impact describe the aggregate. Sub-plans carry independent estimates.

#### REQ: effort-values

When present, the Effort field MUST be one of: `S`, `M`, `L`, or `XL`.

#### REQ: impact-values

When present, the Impact field MUST be one of: `low`, `medium`, `high`, or `critical`.

#### Effort scale

| Effort | Rough meaning |
|--------|---------------|
| S | A few hours of focused work, 1-3 tasks |
| M | A few days, 3-6 tasks, limited dependencies |
| L | A week or more, 5-10 tasks, cross-cutting |
| XL | Multi-week, many tasks, multiple sub-plans or deep dependencies |

#### Impact scale

| Impact | Rough meaning |
|--------|---------------|
| low | Nice-to-have, no users blocked |
| medium | Improves existing capability, some users benefit |
| high | Enables important new capability, many users benefit |
| critical | Unblocks core functionality or other critical work |

### Plans index

Every source namespace with Plans maintains an index at its resolved namespace-root `README.md`. Same-repository storage uses `spec/plans/README.md`; external storage uses the full source repository namespace. Its format — required sections, Contents-table columns, Recently Closed section, and adherence footer — is specified by the [plans-index](../plans-index/README.md) Index-Kind feature.

### Feature README back-reference

Each affected feature's README includes a **Plans** section linking to plans that touch it. Features can reference both top-level plans and sub-plans -- the path disambiguates:

```markdown
## Plans

| Plan                                                                  | Status    | Author | Approved   |
|-----------------------------------------------------------------------|-----------|--------|------------|
| [chat-feature](../../plans/chat-feature.md)                           | Draft     | @alex  | -          |
| [user-auth](../../plans/user-auth.md)                                 | Approved  | @alex  | 2026-03-15 |
| [add-batch-mode](../../plans/add-batch-mode.md)                       | In Review | @alex  | -          |
```

A feature appearing in both a plan and its sub-plan is valid -- the plan covers it broadly, the sub-plan implements a slice. A feature linked only to a top-level plan (no sub-plan yet) signals "planned but not decomposed."

#### REQ: feature-back-reference

Each affected feature's README MUST include a Plans section with a table linking to plans that touch it. The table MUST include columns for Plan, Status, Author, and Approved.

### Task artifacts

An artifact is a named output that a task produces. It is not code (code lives in code repos on branches). It is the metadata, decisions, schemas, and intermediate results that downstream tasks need to do their work.

Examples:

| Artifact | Produced by | Consumed by |
|---|---|---|
| JSON Schema definition | "Define data model" task | "Implement endpoints" task, "Build UI" task |
| API contract (OpenAPI snippet) | "Design API" task | "Implement client" task, "Write integration tests" task |
| Migration plan | "Analyze existing data" task | "Write migration script" task |
| Architecture decision record | "Evaluate auth approach" task | All downstream tasks |
| Test fixtures / seed data | "Generate test data" task | Any task running tests |

#### REQ: artifact-declaration

Plan tasks that produce outputs SHOULD declare them using the `**Produces:**` field with a bulleted list of named artifacts and their descriptions.

#### REQ: artifact-dependency-flow

When a task depends on another task, it MUST have access to that task's declared artifacts. The dependency is made explicit through the `Depends on` and `Produces` fields.

## Workflow

The planning pipeline has three stages. Each can be performed by a human or an AI agent.

```mermaid
graph LR
    A["Trigger<br/>(spec approved)"]
    B["Author<br/>plan"]
    C["Review &<br/>approve"]

    A -->|human or<br/>AI agent| B
    B -->|submit| C
```

### Stage 1: Trigger

Something initiates the need for a plan:

| Trigger | Source |
|---|---|
| New feature spec approved | `spec/features/{feature}/README.md` |
| Change request (proposal) approved | `spec/features/{feature}/proposals/{proposal}/` |
| Manual request | Human decides work is needed |

If auto-planning is enabled in the project configuration, tooling can automatically create a `Draft` plan when a feature spec or proposal reaches `Approved` status. If disabled (the default), a human or external tool initiates plan creation explicitly.

### Stage 2: Author the plan

The plan author (human or AI agent) writes the plan document following the structure defined above.

**When authored by a human:** Write the markdown directly. The spec tooling scaffolds the plan file and template (`specscore plan new`).

**When authored by an AI agent:** The agent receives the feature spec or approved proposal as input context, along with relevant codebase context, and produces the plan document. The agent should have access to:

- The feature spec or approved proposal
- Existing codebase structure (for change requests)
- Other active plans (to avoid conflicts)
- Project conventions

### Stage 3: Review and approve

```mermaid
graph LR
    A["Draft"]
    B["In Review"]
    C["Approved"]

    A -->|submit| B
    A -->|approve| C
    B -->|revisions<br/>requested| A
    B -->|approve| C
```

The review process transitions the plan from `Draft` to `In Review`, and upon approval sets the status to `Approved` and creates an `approved` snapshot. A plan MAY also be approved directly from `Draft`, skipping `In Review`, when no separate review step is needed. The plan remains editable after approval -- future changes are tracked through additional snapshots, and once execution begins `lint --fix` derives the execution-band status from task rollup.

### After approval: Execution handoff

Once approved, the plan's tasks can be consumed by execution tools to generate work items. The exact mechanism depends on the orchestration tool used. For Synchestra integration, see [synchestra.io](https://synchestra.io).

Execution is handled by the orchestration tool, not SpecScore. The plan remains a living document during execution, with snapshots marking significant milestones.

## Integration with Execution Tools

Plan tasks can be mapped to execution units (tasks, work items) by orchestration tools. SpecScore defines the plan format; execution tools consume it.

Key integration points:

- Each plan task can carry metadata (like a task identifier) that execution tools use to create and link work items.
- `Depends on` fields in the plan map to dependency relationships in the execution system.
- Acceptance criteria from plan tasks can be copied into generated work item descriptions.

## Retrospective

Once all tasks reach terminal states, a deviation report can compare planned vs actual:

- **Planned tasks** vs. **actual work items** -- were tasks added, removed, or split?
- **Planned dependencies** vs. **actual execution order**
- **Planned acceptance criteria** vs. **outcomes**
- **Time estimates** (if provided) vs. **actual durations**

The report is a learning artifact. It can be stored alongside the plan:

```
spec/plans/{plan-slug}/
  README.md             <- the plan
  reports/
    README.md           <- deviation report
```

## What's Next Report

The What's Next report is an AI-generated prioritization document that surfaces what to work on next based on plan statuses, dependencies, and ROI metadata.

### Location

`spec/plans/WHATS-NEXT.md`

### Report structure

```markdown
# What's Next

**Generated:** 2026-03-24
**Mode:** incremental | full

## Completed Since Last Update

- [chat-infrastructure](chat-feature/chat-infrastructure/) -- completed 2026-03-20

## In Progress

- [hero-scene](hero-scene/) -- 2/4 tasks done, no blockers

## Recommended Next

1. **[chat-workflow-engine](chat-feature/chat-workflow-engine/)** -- Impact: high,
   Effort: M. Unblocked by chat-infrastructure completion. Advances the
   highest-impact plan.
2. **[agent-skills-framework](agent-skills-framework/)** -- Impact: medium, Effort: L.
   No blockers, independent of current momentum.

### Reasoning

Brief AI explanation of prioritization -- dependency unlocks, ROI ratio,
momentum, competing priorities.

## Open Questions

(ambiguities the AI surfaced during analysis)
```

### Update mechanism

- **Trigger:** plan completion events or status transitions.
- **Incremental mode:** reads previous `WHATS-NEXT.md` + the completion delta. Regenerates only affected sections. Minimizes token usage.
- **Full mode:** scans all features, plans, and statuses. Used for initial generation or to correct incremental drift.
- The file is **committed to git** after each update, providing a history of how priorities evolved over time.

### Prioritization inputs

The AI considers these signals in order of priority:

1. Explicit ROI metadata (effort/impact) when present
2. Dependency graph -- what is newly unblocked by recent completions
3. Momentum -- preference for advancing plans already in progress
4. Feature status -- features closer to "stable" get a boost
5. AI inference from plan complexity when ROI metadata is absent

## Project Configuration

Planning settings are configured in the repo config file. See [Repo Config](../repo-config/README.md).

### Adherence footer

#### REQ: adherence-footer

Every plan document MUST end with an adherence footer per the [Adherence Footer feature](../adherence-footer/README.md). The footer URL MUST be `https://specscore.md/plan-specification`.

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Features are the source artifacts that trigger plans. Plans list affected features; features back-reference active plans in their README. |
| [Task](../task/README.md) | A task is the atomic unit of work. A plan is a composite task -- it contains tasks. The plan feature defines the composite structure; the task feature defines the leaf node properties. |
| [Requirement](../requirement/README.md) | Requirements are `#### REQ:` subsections within a feature's Behavior section. Plan REQs define the rules for plan structure and lifecycle. |
| [Acceptance Criteria](../acceptance-criteria/README.md) | Plan-level and task-level ACs follow the same format as feature ACs. Snapshots provide immutable references to AC state at specific points. |
| [Scenario](../scenario/README.md) | Scenarios in `_tests/` validate plan REQs with concrete Given/When/Then flows. |
| [Proposals](../proposals/README.md) | A proposal (change request) is a trigger for plan creation. Approved proposals link forward to their plan; plans link back to their source proposal. |
| [Open Questions](../open-questions/README.md) | Plan tasks may surface open questions. These follow the existing question lifecycle. |
| [Plans Index](../plans-index/README.md) | The plans-index feature specifies the resolved source namespace aggregation file that lists every Plan for that source project. Plan documents conform to `plan-specification`; the plans-index file conforms to `plans-index-specification`. |
| [Status Vocabulary](../status-vocabulary/README.md) | The canonical source of truth for the legal Plan status values. The Plan set is `Draft`, `In Review`, `Approved`, `Executing`, `Blocked`, `Implemented`, `Failed`, `Rejected`, `Withdrawn`, `Superseded`, `Deprecated`; `Executing` (in-progress role) is a documented conscious divergence from the shared `Implementing` term, governed there. |

## Acceptance Criteria

### AC: plan-document-validity

**Requirements:** plan#req:plan-title-format, plan#req:plan-required-sections, plan#req:required-header-fields, plan#req:source-binding

A plan document has a correctly formatted title (`# Plan: {Title}`), all required sections present (`Summary`, `Approach`, `Tasks`, `Open Questions`), all required header fields populated (Status, the source line, Date, Owner, Supersedes), and exactly one source declared (`Source Feature`, `Source: idea:{slug}`, or `Source: none`). A document that violates any of these is rejected by validation.

### AC: plan-location-validity

**Requirements:** plan#req:plan-file, plan#req:plan-logical-id, plan#req:plan-slug-format, plan#req:external-source-namespace

A Plan is a slug-named directory containing `README.md`. Nested Plans retain their complete logical IDs, arbitrary nesting is accepted, and a legacy flat file cannot duplicate the same logical ID.

### AC: external-plan-boundaries

**Requirements:** plan#req:external-source-namespace, plan#req:source-project-context, plan#req:read-only-plan-operations, plan#req:plan-path-conflicts, plan#req:external-store-lifecycle-lock

An externally routed Plan is stored beneath the full source host/owner/repo namespace. A Plan command launched from the generic destination with `--project` resolves Feature and AC context from that source without modifying it. Traversal, symlink escapes, ambiguous flat/directory IDs, and attempts to overwrite a parent README fail before writes. Read-only and dry-run commands change neither repository, and the destination retains the ignored lifecycle-lock inode across operations.

### AC: status-lifecycle

**Requirements:** plan#req:valid-statuses, plan#req:execution-status-derived, plan#req:status-transitions

A plan's status is always one of the eleven defined values (`Draft`, `In Review`, `Approved`, `Executing`, `Blocked`, `Implemented`, `Failed`, `Rejected`, `Withdrawn`, `Superseded`, `Deprecated`). Prep statuses are human-authored; the execution-band statuses are derived by `lint --fix` from task-status rollup and are never hand-set, and `lint --fix` only transitions from `Approved` onward. Status transitions follow the defined state machine.

### AC: snapshot-integrity

**Requirements:** plan#req:snapshot-table-format, plan#req:snapshot-actions, plan#req:snapshot-git-hash

Snapshots are recorded in a table with Date, Git Hash, Action, and Comment columns. Each snapshot references a valid git commit. Actions include `approved`, `checkpoint`, `completed`, and user-defined values. Snapshots provide immutable reference points without restricting plan mutability.

### AC: plan-structure-constraints

**Requirements:** plan#req:recursive-nesting, plan#req:mixed-children, plan#req:parallel-eligibility, plan#req:task-verifies-feature-ac, plan#req:status-rollup

Plans nest recursively with no artificial depth limit, and inline tasks and child Plans coexist at the same level. Tasks whose `Depends on` is `none` are parallel-eligible. Each task declares a `**Verifies:**` line mapping it to the source Feature's acceptance criteria (enforced by `P-001`). Once `Approved`, the plan's execution-band status derives from its task-status rollup via `lint --fix`.

### AC: tasks-count

**Requirements:** plan#req:tasks-count

A plan carries a derived `tasks_count` equal to its number of direct child tasks, maintained by `specscore spec lint --fix` and surfaced in frontmatter. The value is never hand-authored; lint reconciles it from the actual children.

### AC: cross-artifact-links

**Requirements:** plan#req:feature-back-reference, plan#req:proposal-forward-reference

Affected features back-reference plans in a Plans table. Proposals triggered by change requests include forward references to their plans. Bidirectional traceability is maintained.

## Open Questions

- How should plan tasks reference specific sections of a feature spec when the plan implements only part of a feature?
- What is the exact format for the plan task reference -- should it be structured metadata (YAML frontmatter) or a markdown convention (as shown in examples)?
- Should the deviation report be generated automatically when all tasks complete, or only on demand?
- `tasks_count` migration: existing plans gain `tasks_count` via `specscore spec lint --fix` on next touch (derived, never hand-authored), so no manual backfill is required. (Source Idea: `plan-granularity-improvement`.)
- This Feature now makes the recursive Plan-directory model canonical, and its status enum was expanded to the full prep/execution/disposition lifecycle realizing the Approved `plan-status-lifecycle` Idea. Open from that Idea: exact rollup precedence on mixed task states (encoded here as Failed > Executing > Blocked > Implemented — confirm against real `implement` runs); what triggers `lint --fix` to recompute the execution band (every lint run vs a hook vs the `implement` checkpoints); and whether `Failed` requires human acknowledgement to leave.
- The CLI enforcement of the expanded status enum and the `lint --fix` execution-band derivation are not yet implemented; until then plans use the prep-band statuses (`Draft`/`Approved`) as before. (Realizing Idea: `plan-status-lifecycle`.)

---
*This document follows the https://specscore.md/feature-specification*
