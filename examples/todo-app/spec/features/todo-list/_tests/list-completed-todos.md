# Scenario: List completed todos

**Validates:** todo-list/filter-by-status

## Steps

GIVEN a todo list with two active items and one completed item
WHEN the user runs `todo list --completed`
THEN the output shows exactly one item
AND the item has status `completed`
AND the summary line reads "Showing 1 of 3 items"

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo add "Pay rent"
todo add "Walk dog"
todo complete 3
output=$(todo list --completed)
assert_contains "$output" "Walk dog"
assert_contains "$output" "[completed]"
assert_not_contains "$output" "Buy milk"
assert_contains "$output" "Showing 1 of 3 items"
```
