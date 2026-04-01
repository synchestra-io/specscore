# Unified Plan/Task Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Unify SpecScore and Synchestra terminology around a single recursive task/plan model. Rename "development plan" to "plan," introduce the task concept into SpecScore, remove artificial nesting limits, and replace immutability-after-approval with mutable plans + snapshots.

**Spec:** [Unified Plan/Task Model Design](../specs/2026-04-01-unified-plan-task-model-design.md)

**Scope:** SpecScore repository only. Synchestra alignment is a separate follow-up.

---

## File Map

**Renamed directories:**
- `spec/features/development-plan/` -> `spec/features/plan/`

**Major rewrites:**
- `spec/features/plan/README.md` (was `development-plan/README.md`) — rewrite to unified model
- `spec/features/plan/_tests/*.md` — update/add/remove scenarios for new model

**New files:**
- `spec/features/task/README.md` — new feature: task concept
- `spec/features/task/_tests/README.md` — test index
- `spec/features/task/_tests/*.md` — scenarios for task behavior

**Modified files (reference updates):**
- `spec/features/README.md` — rename entry, add task feature
- `spec/features/feature/README.md` — update references from "development plan" to "plan"
- `spec/features/acceptance-criteria/README.md` — update references
- `spec/features/source-references/README.md` — update references
- `spec/README.md` — update references
- `README.md` — update references
- `AGENTS.md` — update architecture description
- `tools/site-generator/site-config.json` — rename slug/path
- `tools/site-generator/landing.html` — update text
- `docs/for/*.md` — update terminology in audience docs

---

### Task 1: Create the Task feature spec

**Files:**
- New: `spec/features/task/README.md`
- New: `spec/features/task/_tests/README.md`

This is the foundation — defines what a task **is** in SpecScore.

- [ ] **Step 1: Write `spec/features/task/README.md`**

  Define the task feature following SpecScore conventions. Key sections:

  - **Summary:** A task is the atomic unit of work in SpecScore. It is a leaf node — actionable work that an agent or human picks up and completes.
  - **Problem:** No standard definition of executable work items in the spec methodology.
  - **Behavior:**
    - Task structure: directory with `README.md`, properties (status, dependencies, acceptance criteria, artifacts)
    - Task statuses and lifecycle: `planning`, `queued`, `in_progress`, `blocked`, `complete`, `failed`, `aborted`
    - Status transitions (from Synchestra's current model)
    - Dependency references: sibling (bare slug), cousin (relative path), cross-project (URL)
    - Task status board format: columns, layout, recently finished section
    - Task properties: `depends_on`, `produces`, acceptance criteria
  - **Interaction with Other Features:** Plan (task is a leaf within a plan), Feature, Requirement, Acceptance Criteria
  - **Acceptance Criteria:** task-structure, status-lifecycle, dependency-references, board-format
  - **Outstanding Questions** section

  Use REQ: prefixed requirements throughout the Behavior section, following existing convention.

- [ ] **Step 2: Write `spec/features/task/_tests/README.md`**

  Create test index listing all scenario files. Start with these scenarios:
  - `task-status-transitions.md` — valid and invalid status transitions
  - `dependency-resolution.md` — sibling, cousin, and cross-project reference resolution
  - `board-format-validation.md` — board table structure and column requirements

- [ ] **Step 3: Write initial test scenarios**

  Create the scenario files listed in Step 2, using Given/When/Then format consistent with existing `_tests/` files in other features.

---

### Task 2: Rewrite the Development Plan feature as Plan

**Files:**
- Rename: `spec/features/development-plan/` -> `spec/features/plan/`
- Rewrite: `spec/features/plan/README.md`

Depends on: Task 1

- [ ] **Step 1: Rename directory from `development-plan` to `plan`**

  ```bash
  git mv spec/features/development-plan spec/features/plan
  ```

- [ ] **Step 2: Rewrite `spec/features/plan/README.md`**

  Major changes from the current spec:

  **Terminology:**
  - "Development plan" -> "Plan" everywhere
  - "Steps" -> "Tasks" (plan steps become tasks)
  - "Sub-steps" -> "Subtasks"

  **Structural changes:**
  - Remove the 2-level nesting limit for steps (REQ: nesting-limit) — tasks nest recursively
  - Remove the 2-level hierarchy limit for roadmap/child plans (REQ: hierarchy-nesting-limit) — plans nest recursively
  - Merge roadmap and plan concepts: a plan with sub-plans is just a composite task whose children are also composite. Remove the separate "Roadmap" concept and its special rules (REQ: roadmap-no-steps, REQ: roadmap-status-derivation)
  - Remove `superseded` status — plans are mutable, snapshots track history
  - Remove immutability-after-approval (REQ: immutability-after-approval) — replace with snapshots

  **New concepts to add:**
  - **Recursive task/plan model:** A plan is a composite task. A task with subtasks is a plan. No depth limit.
  - **Snapshots:** Table in the plan README with date, git hash, action, and comment. Replaces immutability. Actions include `approved`, `checkpoint`, `completed`, and user-defined.
  - **Mixed children:** Tasks and sub-plans can coexist at the same level. No forced restructuring.
  - **Status rollup:** A composite task's status derives from its children.

  **Preserved (renamed):**
  - Plan location under `spec/plans/`
  - Plan document structure (header fields, context, acceptance criteria) — but "Steps" section becomes "Tasks" section
  - `Depends on` and `Produces` fields — now on tasks instead of steps
  - Plan-level and task-level acceptance criteria (was step-level)
  - Plan statuses: `draft`, `in_review`, `approved` (remove `superseded`)
  - Status transitions (simplified: remove `superseded` transitions)
  - Plan index (`spec/plans/README.md`)
  - Feature README back-references
  - ROI metadata (Effort/Impact)
  - What's Next report

  **Updated plan document example:**
  - Show tasks instead of steps, with a nested sub-plan example
  - Show snapshot table
  - Show mixed children (tasks and sub-plans at same level)

- [ ] **Step 3: Update or remove test scenarios**

  Review each file in `spec/features/plan/_tests/`:

  - **Remove:** `immutability-after-approval.md`, `nesting-limit-enforcement.md`, `roadmap-structure.md`, `roadmap-status-derivation.md`
  - **Update:** All remaining scenarios to use "plan"/"task" terminology instead of "development plan"/"step"
  - **Add:** `snapshot-lifecycle.md`, `recursive-nesting.md`, `mixed-children.md`, `status-rollup.md`

---

### Task 3: Update cross-references across the repository

**Files:**
- Modify: `spec/features/README.md`
- Modify: `spec/features/feature/README.md`
- Modify: `spec/features/feature/_tests/path-based-identification.md`
- Modify: `spec/features/acceptance-criteria/README.md`
- Modify: `spec/features/acceptance-criteria/_tests/plan-ac-references-feature-ac.md`
- Modify: `spec/features/source-references/README.md`
- Modify: `spec/README.md`
- Modify: `README.md`
- Modify: `AGENTS.md`

Depends on: Task 2

- [ ] **Step 1: Update `spec/features/README.md`**

  - Rename the `development-plan` row to `plan` with updated path
  - Add a new row for the `task` feature
  - Update the directory tree to show `plan/` and `task/`

- [ ] **Step 2: Update `spec/features/feature/README.md`**

  Replace all references to "development plan" with "plan" and update link paths from `development-plan/` to `plan/`.

- [ ] **Step 3: Update `spec/features/acceptance-criteria/README.md`**

  Replace all references to "development plan" with "plan" and update link paths.

- [ ] **Step 4: Update `spec/features/source-references/README.md`**

  Replace references to "development plan" with "plan" and update link paths.

- [ ] **Step 5: Update `spec/README.md`, `README.md`, and `AGENTS.md`**

  - `spec/README.md`: Update any feature list or references
  - `README.md`: Update terminology
  - `AGENTS.md`: Update architecture description — change "development plans" to "plans" in the `spec/` description

---

### Task 4: Update site generator configuration and landing page

**Files:**
- Modify: `tools/site-generator/site-config.json`
- Modify: `tools/site-generator/landing.html`

Depends on: Task 2

- [ ] **Step 1: Update `site-config.json`**

  Rename the `development-plan` entry: update `source` path, `slug`, `title`, and `navLabel`. Add an entry for the new `task` feature page.

- [ ] **Step 2: Update `landing.html`**

  Replace "Development Plan" text with "Plan" in any feature lists or descriptions. Add mention of the Task concept if the landing page lists features.

---

### Task 5: Update audience documentation

**Files:**
- Modify: `docs/for/developers.md`
- Modify: `docs/for/architects.md`
- Modify: `docs/for/project-managers.md`
- Modify: `docs/for/qas.md`

Depends on: Task 2

- [ ] **Step 1: Update all docs/for/*.md files**

  Replace "development plan" with "plan" in all audience-facing documentation. Update any descriptions to reflect the unified task/plan model.

---

## Dependency Graph

```
Task 1 (Create Task feature spec)
  |
  v
Task 2 (Rewrite Plan feature spec)
  |
  +--> Task 3 (Update cross-references)
  |
  +--> Task 4 (Update site generator)
  |
  +--> Task 5 (Update audience docs)
```

Tasks 3, 4, and 5 are parallel-eligible after Task 2 completes.

## Risks and Open Decisions

- **Synchestra alignment:** This plan covers SpecScore only. Synchestra's task-status-board spec references SpecScore's development-plan feature via URL. Those references will break and need updating in a follow-up.
- **Existing plans in spec/plans/:** If any concrete plans exist using the old step format, they need migration. Check for `spec/plans/` directories.
- **Site build:** After renaming, the site generator must be re-run to verify no broken links. Include a build verification step.

## Outstanding Questions

None at this time.
