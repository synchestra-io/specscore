# Scenario: Edit a todo title

**Validates:** todo-item/manage/edit-updates-fields

## Steps

GIVEN a todo list with one item titled "Buy milk"
WHEN the user runs `todo edit 1 --title "Buy oat milk"`
THEN the exit code is 0
AND `todo list` shows the item with title "Buy oat milk"
AND the item ID is still 1

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo edit 1 --title "Buy oat milk"
assert_exit_code 0
output=$(todo list)
assert_contains "$output" "Buy oat milk"
assert_not_contains "$output" "Buy milk"
```
