# Feature: Todo Item

**Status:** Conceptual

## Summary

A todo item is the core entity of the application. It has a title (required), an optional description, and a status that tracks its lifecycle. This parent feature defines the entity structure; sub-features define operations (manage) and lifecycle transitions (completion).

## Contents

| Directory | Description |
|---|---|
| [manage/](manage/README.md) | Create, edit, and delete todo items |
| [completion/](completion/README.md) | Complete and reopen todo items |

### manage

CRUD operations for todo items — creating new items with validation, editing existing items, and deleting items. Defines the rules for what makes a valid todo item.

### completion

Lifecycle transitions for todo items — marking items as completed (with timestamp) and reopening completed items. Defines the rules for status transitions and their side effects.

## Problem

A todo application needs a well-defined entity at its core. Without clear rules for what a todo item is and what fields it carries, sub-features (management, completion, listing, due dates) have no stable foundation to build on.

## Behavior

### REQ: has-title

Every todo item MUST have a `title` field. The title is a non-empty string.

### REQ: has-status

Every todo item MUST have a `status` field with value `active` or `completed`. New items are created with status `active`.

### REQ: has-id

Every todo item MUST have a unique numeric `id` assigned on creation. IDs are sequential starting from 1 and MUST NOT be reused after deletion.

### Entity fields

| Field | Type | Required | Default | Notes |
|---|---|---|---|---|
| `id` | integer | Yes | Auto-assigned | Sequential, never reused |
| `title` | string | Yes | — | Non-empty |
| `description` | string | No | empty | Free-form text |
| `status` | enum | Yes | `active` | `active` or `completed` |
| `created_at` | timestamp | Yes | Auto-assigned | Set on creation |
| `completed_at` | timestamp | No | null | Set on completion, cleared on reopen |
| `due_date` | date | No | null | See [due-dates](../due-dates/README.md) |

## Dependencies

- [due-dates](../due-dates/README.md) (optional — extends the entity with `due_date` field)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [todo-list](../todo-list/README.md) | Todo list displays and filters todo items |
| [due-dates](../due-dates/README.md) | Due dates adds an optional `due_date` field to the entity |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

- Acceptance criteria not yet defined for this feature.
