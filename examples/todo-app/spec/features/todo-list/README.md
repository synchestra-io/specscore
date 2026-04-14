# Feature: Todo List

**Status:** Draft

## Summary

Display, filter, and summarize todo items. The list is the primary view of the application — it shows todos filtered by status with a count summary.

## Problem

A todo application is only useful if you can see your todos. Users need to view active items (the default), review completed work, see everything at once, and quickly understand how much remains. Without filtering and summaries, the list becomes noise.

## Behavior

### REQ: default-filter-active

`todo list` with no flags MUST show only items with status `active`.

### REQ: filter-completed

`todo list --completed` MUST show only items with status `completed`.

### REQ: filter-all

`todo list --all` MUST show all items regardless of status.

### REQ: mutually-exclusive-filters

The flags `--all`, `--completed`, and `--overdue` are mutually exclusive. Providing more than one MUST produce an error.

### REQ: count-summary

The list output MUST include a summary line showing the count of displayed items and total items. Format: `Showing N of M items`.

### REQ: empty-list-message

When the filtered list is empty, the output MUST show a "No todos" message instead of blank output.

### Display format

Each item in the list is displayed as:

```
<id>. [<status>] <title>
```

Example:
```
1. [active] Buy milk
2. [active] Pay rent
Showing 2 of 3 items
```

## Dependencies

- [todo-item](../todo-item/README.md)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [todo-item](../todo-item/README.md) | Lists display todo items; item creation/deletion changes list contents |
| [due-dates](../due-dates/README.md) | Due dates adds `--overdue` filter and optional due date display |

## Acceptance Criteria

| AC | Requirement | Description |
|---|---|---|
| [default-shows-active](_acs/default-shows-active.md) | `todo-list#req:default-filter-active` | Default list shows active items only |
| [filter-by-status](_acs/filter-by-status.md) | `todo-list#req:filter-completed`, `todo-list#req:filter-all` | Status filters work correctly |
| [count-summary](_acs/count-summary.md) | `todo-list#req:count-summary` | Summary line shows correct counts |
| [empty-list-message](_acs/empty-list-message.md) | `todo-list#req:empty-list-message` | Empty lists show a message |

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
