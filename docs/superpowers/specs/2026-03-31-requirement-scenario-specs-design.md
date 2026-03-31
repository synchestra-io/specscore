# Requirement & Scenario Spec Definitions + Demo Todo App

## Goal

Introduce two new SpecScore spec features — **Requirement** and **Scenario** — that formalize the layers between Feature and execution. Validate the concepts by creating a demo CLI todo-app with full four-layer specifications (Feature → Requirements → ACs → Scenarios).

## Background

SpecScore currently has Feature and Acceptance Criteria as spec concepts. Analysis of OpenSpec.dev and internal discussion identified two gaps:

1. **Requirements** — Features embed behavioral rules in prose Behavior sections, but individual rules are not addressable. ACs cannot trace back to specific obligations.
2. **Scenarios** — ACs blur the line between abstract success conditions and concrete testable flows. Separating them improves clarity and enables Rehearse automation.

## Concept Definitions

### Requirement

A discrete, testable rule or condition the system must satisfy. Requirements live as named subsections within a feature's Behavior section — they are a naming convention, not a separate file artifact.

**Format in Behavior section:**

```markdown
### REQ: {slug}

{Rule description using RFC 2119 language: MUST, MUST NOT, SHOULD, etc.}
```

**ID scheme:** `{feature-path}#req:{slug}`

Examples:
- `todo-item#req:title-required`
- `todo-item/completion#req:timestamp-on-complete`
- `todo-list#req:default-filter-active`

**Rules:**
- Slugs are lowercase, hyphen-separated, URL-safe
- Each requirement is a single testable obligation
- ACs reference requirement IDs via `**Requirement:**` metadata
- Subsection groupings (### headings without REQ: prefix) remain valid for organization — not every subsection is a requirement

### Acceptance Criteria

Tightened from current definition: an AC is an **abstract verification condition** — a success/failure statement that can be checked, but does not prescribe specific inputs or flows.

**Example AC** (`_acs/title-required.md`):
```markdown
# AC: title-required

**Requirement:** todo-item/manage#req:title-required

Creating a todo without a title is rejected. Creating a todo with a title succeeds.
```

ACs answer: "What must be true?" — not "How do we test it?"

### Scenario

A concrete example of system behavior written in Given/When/Then format. Scenarios live in `_tests/` and are executable by Rehearse.

**Format** (`_tests/{scenario-slug}.md`):
```markdown
# Scenario: {Title}

**Validates:** {feature-slug}/{ac-slug}, {feature-slug}/{ac-slug-2}

## Steps

GIVEN {initial condition}
WHEN {action}
THEN {expected outcome}
AND {additional outcome}
```

**Rules:**
- `**Validates:**` is optional, links to ACs (many-to-many)
- Scenarios without AC links are valid (exploratory, integration, cross-feature)
- Each scenario includes a `rehearse` fenced code block for execution
- Scenarios can reference setup/teardown flows in `_tests/flows/`

### Relationship Summary

```
Feature
├── Behavior section
│     └── Requirements (### REQ: named rules)
│           └── _acs/ (abstract conditions verifying requirements)
└── _tests/ (concrete Given/When/Then scenarios validating ACs)
```

Links: Scenario →validates→ AC →requirement→ Requirement →defined-in→ Feature

## Part 1: New Spec Feature Documents

### `spec/features/requirement/README.md`

New feature spec defining:
- What a requirement is (discrete testable rule)
- Convention format (`### REQ: {slug}`)
- ID scheme and referencing
- Relationship to features and ACs
- Structural rules

### `spec/features/scenario/README.md`

New feature spec defining:
- What a scenario is (concrete behavior example)
- Given/When/Then format
- File location (`_tests/`)
- Validates metadata (links to ACs)
- Relationship to Rehearse test runner
- Distinction from ACs

### Updates to Existing Specs

**`spec/features/feature/README.md`:**
- Add `_tests/` to the directory structure diagram
- Add Requirements convention to the Behavior section description
- Add Requirement and Scenario to the Interaction with Other Features table

**`spec/features/acceptance-criteria/README.md`:**
- Tighten definition to "abstract verification conditions"
- Add explicit AC vs Scenario distinction
- Add Requirement reference metadata (`**Requirement:**` field)
- Add Scenario to the Interaction table

**`spec/features/README.md`:**
- Add requirement and scenario to the features index table

## Part 2: Demo Todo CLI App

Location: `examples/todo-app/`

A simple CLI todo application specified using all four SpecScore layers. No implementation code — specs only.

### Features

```
examples/todo-app/spec/features/
  README.md                    # Feature index
  todo-item/
    README.md                  # Parent: entity definition, shared reqs
    manage/
      README.md                # Sub-feature: create, edit, delete
      _acs/
        title-required.md
        added-to-list.md
        edit-updates-fields.md
        delete-removes-item.md
      _tests/
        create-todo.md
        create-without-title.md
        edit-todo.md
        delete-todo.md
    completion/
      README.md                # Sub-feature: complete/reopen lifecycle
      _acs/
        status-transition.md
        completion-timestamp.md
        reopen-clears-timestamp.md
      _tests/
        complete-todo.md
        reopen-todo.md
        complete-already-completed.md
  todo-list/
    README.md                  # Filtering, counting, display
    _acs/
      default-shows-active.md
      filter-by-status.md
      count-summary.md
    _tests/
      list-active-todos.md
      list-completed-todos.md
      list-all-todos.md
      empty-list.md
  due-dates/
    README.md                  # Optional due dates, overdue detection
    _acs/
      due-date-optional.md
      overdue-detection.md
      no-due-date-not-overdue.md
    _tests/
      add-due-date.md
      overdue-todo.md
      filter-overdue.md
      completed-past-due-not-overdue.md
```

### Todo CLI Interface

The demo app is a CLI with these commands:

```
todo add "Buy milk"                    # create
todo add "Pay rent" --due 2026-04-15   # create with due date
todo list                              # list active (default)
todo list --all                        # list all
todo list --completed                  # list completed
todo list --overdue                    # list overdue
todo edit 1 --title "Buy oat milk"     # edit
todo complete 1                        # mark complete
todo reopen 1                          # reopen
todo delete 1                          # delete
```

### Cross-Feature Scenario

`due-dates/_tests/filter-overdue.md` exercises all three features:
- Creates todos (todo-item/manage)
- Sets due dates (due-dates)
- Filters the list (todo-list)
- Validates overdue detection (due-dates)

This validates that scenarios can span feature boundaries.

### Rehearse Script Format

Each scenario includes a fenced `rehearse` block:

````markdown
```rehearse
#!/bin/bash
todo add "Pay rent" --due 2026-01-01
output=$(todo list --overdue)
assert_contains "$output" "Pay rent"
assert_not_contains "$output" "Buy milk"
```
````

These are executable shell scripts that the Rehearse test runner can invoke.

## What This Does NOT Include

- Actual Go/code implementation of the todo CLI
- Rehearse framework changes
- Changes to `pkg/` tooling (lint rules for requirements, etc.)
- specscore-project.yaml for the demo (could add later)

## File Changes Summary

**New files:**
- `spec/features/requirement/README.md`
- `spec/features/scenario/README.md`
- `examples/todo-app/` — ~25 spec files

**Modified files:**
- `spec/features/README.md` — add to index
- `spec/features/feature/README.md` — add requirement/scenario references
- `spec/features/acceptance-criteria/README.md` — tighten definition
