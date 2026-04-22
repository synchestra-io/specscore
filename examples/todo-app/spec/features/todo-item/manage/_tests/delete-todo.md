# Scenario: Delete a todo item

**Validates:** todo-item/manage/delete-removes-item

## Steps

GIVEN a todo list with one item titled "Buy milk"
WHEN the user runs `todo delete 1`
THEN the exit code is 0
AND `todo list` shows no items
AND `todo list --all` shows no items

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo delete 1
assert_exit_code 0
output=$(todo list --all)
assert_contains "$output" "No todos"
```

---
*This document follows the https://specscore.md/scenario-specification*
