# Scenario: List all todos

**Validates:** todo-list/filter-by-status, todo-list/count-summary

## Steps

GIVEN a todo list with two active items and one completed item
WHEN the user runs `todo list --all`
THEN the output shows all three items
AND the summary line reads "Showing 3 of 3 items"

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo add "Pay rent"
todo add "Walk dog"
todo complete 3
output=$(todo list --all)
assert_contains "$output" "Buy milk"
assert_contains "$output" "Pay rent"
assert_contains "$output" "Walk dog"
assert_contains "$output" "Showing 3 of 3 items"
```
