# Feature: Manage Todos

**Status:** Draft

## Summary

Create, edit, and delete todo items. This sub-feature defines the CRUD operations and validation rules for todo item management.

## Problem

Users need to create, modify, and remove todo items. Without clear rules for validation (what makes a valid todo?) and behavior (what happens on delete?), the CLI cannot give consistent, predictable responses.

## Behavior

### REQ: title-required

Creating a todo MUST require a non-empty title. The CLI MUST reject `todo add ""` and `todo add` (no argument) with an error message.

### REQ: added-to-list

When a todo is created with a valid title, it MUST appear in the active list with status `active` and the exact title provided.

### REQ: title-max-length

A todo title MUST NOT exceed 200 characters. The CLI MUST reject titles longer than 200 characters with an error message.

### REQ: edit-title

`todo edit <id> --title <text>` MUST update the title of the specified item. The new title MUST follow the same validation rules as creation (non-empty, max 200 characters).

### REQ: edit-nonexistent

Editing a non-existent ID MUST produce an error message and exit with a non-zero status code.

### REQ: delete-removes

`todo delete <id>` MUST remove the item from all listings. The ID MUST NOT be reused.

### REQ: delete-nonexistent

Deleting a non-existent ID MUST produce an error message and exit with a non-zero status code.

## Acceptance Criteria

| AC | Requirement | Description |
|---|---|---|
| [title-required](_acs/title-required.md) | `todo-item/manage#req:title-required` | Empty titles are rejected |
| [added-to-list](_acs/added-to-list.md) | `todo-item/manage#req:added-to-list` | Valid todos appear in the list |
| [edit-updates-fields](_acs/edit-updates-fields.md) | `todo-item/manage#req:edit-title` | Editing updates the correct fields |
| [delete-removes-item](_acs/delete-removes-item.md) | `todo-item/manage#req:delete-removes` | Deleted items disappear from all views |

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
