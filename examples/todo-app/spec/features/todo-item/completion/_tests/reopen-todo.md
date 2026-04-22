# Scenario: Reopen a completed todo

**Validates:** todo-item/completion/status-transition, todo-item/completion/reopen-clears-timestamp

## Steps

GIVEN a completed todo item titled "Buy milk"
WHEN the user runs `todo reopen 1`
THEN the exit code is 0
AND `todo list` shows "Buy milk" with status `active`
AND `todo list --completed` shows no items

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo complete 1
todo reopen 1
assert_exit_code 0
active=$(todo list)
assert_contains "$active" "Buy milk"
assert_contains "$active" "[active]"
completed=$(todo list --completed)
assert_contains "$completed" "No todos"
```

---
*This document follows the https://specscore.md/scenario-specification*
