# Feature: Todo Completion

**Status:** Draft

## Summary

Complete and reopen todo items. This sub-feature defines the lifecycle transitions between `active` and `completed` status, including timestamps and idempotency rules.

## Problem

Completing a todo is not a simple boolean flip — it carries side effects (timestamp recording) and edge cases (completing an already-completed item, reopening). Without explicit rules, implementations diverge on these details.

## Behavior

### REQ: complete-sets-status

`todo complete <id>` MUST set the item's status to `completed`.

### REQ: timestamp-on-complete

Completing a todo MUST set `completed_at` to the current timestamp.

### REQ: complete-idempotent

Completing an already-completed todo MUST succeed without error. The `completed_at` timestamp MUST NOT be updated.

### REQ: reopen-sets-active

`todo reopen <id>` MUST set the item's status back to `active`.

### REQ: reopen-clears-timestamp

Reopening a todo MUST clear the `completed_at` timestamp to null.

### REQ: reopen-active-error

Reopening an already-active todo MUST produce an error message and exit with a non-zero status code.

## Acceptance Criteria

| AC | Requirement | Description |
|---|---|---|
| [status-transition](_acs/status-transition.md) | `todo-item/completion#req:complete-sets-status`, `todo-item/completion#req:reopen-sets-active` | Status transitions work correctly |
| [completion-timestamp](_acs/completion-timestamp.md) | `todo-item/completion#req:timestamp-on-complete` | Completion records a timestamp |
| [reopen-clears-timestamp](_acs/reopen-clears-timestamp.md) | `todo-item/completion#req:reopen-clears-timestamp` | Reopening clears the timestamp |

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
