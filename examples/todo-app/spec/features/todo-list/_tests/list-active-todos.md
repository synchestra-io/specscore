# Scenario: List active todos (default)

**Validates:** todo-list/default-shows-active, todo-list/count-summary

## Steps

GIVEN a todo list with two active items and one completed item
WHEN the user runs `todo list`
THEN the output shows exactly two items
AND both items have status `active`
AND the summary line reads "Showing 2 of 3 items"

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo add "Pay rent"
todo add "Walk dog"
todo complete 3
output=$(todo list)
assert_contains "$output" "Buy milk"
assert_contains "$output" "Pay rent"
assert_not_contains "$output" "Walk dog"
assert_contains "$output" "Showing 2 of 3 items"
```

---
*This document follows the https://specscore.md/scenario-specification*
