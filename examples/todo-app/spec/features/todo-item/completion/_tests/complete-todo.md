# Scenario: Complete a todo item

**Validates:** todo-item/completion/status-transition, todo-item/completion/completion-timestamp

## Steps

GIVEN a todo list with one active item titled "Buy milk"
WHEN the user runs `todo complete 1`
THEN the exit code is 0
AND `todo list` shows no items (default shows active only)
AND `todo list --completed` shows "Buy milk" with status `completed`

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo complete 1
assert_exit_code 0
active=$(todo list)
assert_contains "$active" "No todos"
completed=$(todo list --completed)
assert_contains "$completed" "Buy milk"
assert_contains "$completed" "[completed]"
```

---
*This document follows the https://specscore.md/scenario-specification*
