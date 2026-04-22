# Scenario: Create a todo item

**Validates:** todo-item/manage/added-to-list

## Steps

GIVEN an empty todo list
WHEN the user runs `todo add "Buy milk"`
THEN the exit code is 0
AND the output confirms the item was created
AND `todo list` shows exactly one item with title "Buy milk" and status `active`

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
assert_exit_code 0
output=$(todo list)
assert_contains "$output" "Buy milk"
assert_contains "$output" "[active]"
count=$(echo "$output" | grep -c "Buy milk")
assert_eq "$count" "1"
```

---
*This document follows the https://specscore.md/scenario-specification*
