# Feature: Due Dates

**Status:** Draft

## Summary

Optional due dates for todo items with overdue detection. Items past their due date and still active are flagged as overdue. The list supports an `--overdue` filter to show only overdue items.

## Problem

Todos without deadlines have no urgency signal. Users need to assign due dates, see which items are overdue, and filter for overdue items to prioritize their work. The challenge is handling edge cases: items with no due date, items completed after their due date, and the interaction between overdue filtering and the existing status filters.

## Behavior

### REQ: due-date-optional

The `--due` flag on `todo add` and `todo edit` is optional. Items without a due date function normally.

### REQ: due-date-format

Due dates MUST be specified in `YYYY-MM-DD` format. Invalid formats MUST be rejected with an error message.

### REQ: overdue-detection

A todo item is overdue when ALL of the following are true:
- It has a due date
- The due date is in the past (before today)
- Its status is `active`

### REQ: completed-not-overdue

A completed todo item MUST NOT be considered overdue, even if its due date is in the past.

### REQ: no-due-date-not-overdue

A todo item without a due date MUST NOT be considered overdue.

### REQ: overdue-filter

`todo list --overdue` MUST show only items that are overdue (as defined by the overdue detection rule).

### Display format

Items with due dates show the date after the status:

```
1. [active] Buy milk (due: 2026-04-15)
2. [active] Pay rent (due: 2026-03-01) [OVERDUE]
```

Overdue items are tagged with `[OVERDUE]`.

## Dependencies

- [todo-item](../todo-item/README.md)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [todo-item](../todo-item/README.md) | Extends the todo entity with an optional `due_date` field |
| [todo-item/completion](../todo-item/completion/README.md) | Completion status affects overdue detection — completed items are never overdue |
| [todo-list](../todo-list/README.md) | Adds `--overdue` filter to the list command; `--overdue` is mutually exclusive with `--completed` and `--all` |

## Acceptance Criteria

| AC | Requirement | Description |
|---|---|---|
| [due-date-optional](_acs/due-date-optional.md) | `due-dates#req:due-date-optional` | Due dates are optional |
| [overdue-detection](_acs/overdue-detection.md) | `due-dates#req:overdue-detection` | Overdue is correctly determined |
| [no-due-date-not-overdue](_acs/no-due-date-not-overdue.md) | `due-dates#req:no-due-date-not-overdue` | Items without due dates are not overdue |
| [completed-not-overdue](_acs/completed-not-overdue.md) | `due-dates#req:completed-not-overdue` | Completed items are never overdue |

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
