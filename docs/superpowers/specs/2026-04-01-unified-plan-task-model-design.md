# Unified Plan/Task Model

## Goal

Unify SpecScore and Synchestra terminology around a single recursive concept: **tasks**. A task with subtasks is a **plan**. This replaces the current split between "development plans" (SpecScore) and "tasks" (Synchestra) with one coherent model that spans both systems.

## Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Core concept | Task | One term, universally understood |
| Composite task | Plan | A task with subtasks is called a plan — emergent from structure, not declared |
| Nesting depth | Unlimited (recursive) | No artificial 2-level limits; depth is a judgment call |
| "Development plan" rename | Plan | Simpler, aligns with the unified model |
| Plan mutability | Always mutable, snapshots on demand | Real plans change; frozen immutability is a fiction |
| Snapshot mechanism | Git hash + action + comment | No file duplication; git stores the state |
| Task concept ownership | SpecScore (methodology) | Synchestra implements; SpecScore defines |

## The Unified Model

### Tasks all the way down

There is one structural concept: the **task**. A task is a unit of work. Some tasks have subtasks — these are called **plans**. A subtask can itself be a plan. The type is determined by structure, not declaration:

- **Has children** -> it's a plan
- **No children** -> it's a task
- **Add subtasks to a task** -> it becomes a plan automatically

```
Plan (root)
  |-- Task: set up project structure
  |-- Task: Build auth middleware (plan)
  |    |-- Task: implement JWT
  |    |-- Task: refresh tokens
  |    +-- Task: session management
  |-- Task: write integration tests
  +-- ...
```

There is no need to choose upfront whether something is a "task" or a "plan." A leaf task becomes a plan the moment subtasks are added. No restructuring required.

### Task properties

Every task (leaf or composite) carries:

- **Dependencies** (`depends_on`) — which sibling tasks must complete first
- **Parallelism** — tasks without dependencies are parallel-eligible
- **Acceptance criteria** — what "done" looks like
- **Artifacts** (`produces`) — named outputs downstream tasks need

For composite tasks (plans), status rolls up: a plan is "done" when all its subtasks are done, "in progress" when any are, "blocked" when dependencies aren't met.

### Mixed children

Tasks and sub-plans can coexist at the same level. This avoids forced restructuring when a user groups a few tasks into a sub-plan:

```
Plan
  |-- Task: set up project structure
  |-- Plan: Auth middleware (3 subtasks)
  |    |-- Task: implement JWT
  |    |-- Task: refresh tokens
  |    +-- Task: session management
  |-- Task: write integration tests
  +-- Task: deploy to staging
```

No rule that says "branches contain only branches or only leaves." Practical flexibility over theoretical cleanliness.

### Plans are mutable, snapshots are immutable

Plans are living documents — they evolve as execution reveals complexity. There is no "frozen after approval" rule. Instead, users take **snapshots** when they need a reference point.

A snapshot is minimal:

| Date | Hash | Action | Comment |
|---|---|---|---|
| 2026-03-15 | `a3f7c2d` | approved | Ready for execution |
| 2026-03-20 | `e1b4a8f` | checkpoint | Mid-sprint, added auth subtasks |
| 2026-03-28 | `9d2f1c3` | completed | All tasks done |

The git hash **is** the snapshot. `git show a3f7c2d:path/to/plan/README.md` recovers the exact state. `git diff a3f7c2d..9d2f1c3 -- path/to/plan/` produces the retrospective diff. No separate snapshot files, no duplication.

"Approved" becomes just one possible snapshot action — not a special status with special rules.

### Execution as a shared work queue

An implementation plan (the execution-phase plan) is not per-machine or per-agent. It is a **shared task queue**. Multiple agents on multiple machines pull from the same queue:

```
Plan (shared task queue)
  |-- Task 1  <-- picked by agent on machine A
  |-- Task 2  <-- picked by agent on machine B
  |-- Task 3  <-- picked by agent in cloud
  +-- Task 4  <-- waiting in queue
```

With shared state management and a message bus that crosses machine boundaries, agents are equivalent workers regardless of where they run. The plan owns the **what**; infrastructure handles the **where**.

## Separation of Concerns: SpecScore vs Synchestra

### SpecScore = methodology (open source)

SpecScore defines what things **are** and how they relate:

- What a task is (leaf work item)
- What a plan is (composite task)
- Task/plan structure: dependencies, acceptance criteria, artifacts, ordering
- Snapshots and their format
- Task statuses and lifecycle
- Task status board format and layout
- Dependency reference model (sibling, cousin, cross-project)

### Synchestra = tooling (freemium product)

Synchestra **implements** the methodology with operational tooling:

- Agent claim protocol (optimistic locking via git)
- Remote execution across machines
- Message bus and shared state management
- CLI commands (`task claim`, `task start`, etc.)
- Conflict resolution mechanics
- Board retention settings and automation
- Agent/Branch/Time execution metadata

### The litmus test

If someone follows SpecScore methodology with just markdown and git (no Synchestra), do they need it? If yes, it belongs in SpecScore. If it's about automation, orchestration, or multi-agent coordination, it belongs in Synchestra.

The analogy: SpecScore is the HTTP spec. Synchestra is Nginx. The spec is public and open — that's how people understand and adopt the methodology. The product adds the operational layer.

Synchestra's value proposition: "You **could** follow SpecScore manually with markdown and git. Synchestra automates it."

## What Changes

### In SpecScore

| Current | New |
|---|---|
| "Development plan" feature | "Plan" feature |
| Steps and sub-steps (max 2 levels) | Tasks (unlimited recursive nesting) |
| Plan hierarchy: roadmap -> child plan (max 2 levels) | Single recursive task/plan model |
| Immutable after approval | Mutable always, snapshots on demand |
| `superseded` status | Removed — just edit the plan, snapshot history shows what changed |
| No task concept | Task concept added (from Synchestra) |
| No task status board concept | Task status board format defined (from Synchestra) |

### In Synchestra

| Current | New |
|---|---|
| Task defined locally | Task definition references SpecScore |
| Task status board defined locally | Board format references SpecScore; implementation stays local |
| Independent terminology | Aligned with SpecScore |

## Outstanding Questions

- Should the snapshot table live in the plan's README or in a separate `snapshots.md` file?
- How should the "What's Next" report adapt to the recursive model — does it operate at a specific depth or traverse the full tree?
- What is the migration path for existing plans that use the current step/sub-step format?
