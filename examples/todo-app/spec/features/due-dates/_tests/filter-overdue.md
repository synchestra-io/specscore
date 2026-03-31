# Scenario: Filter overdue todos (cross-feature)

**Validates:** due-dates/overdue-detection, due-dates/no-due-date-not-overdue, todo-list/filter-by-status

## Steps

GIVEN a todo "Pay rent" with due date 2026-01-01 (past)
AND a todo "Buy milk" with no due date
AND a todo "File taxes" with due date 2026-12-31 (future)
WHEN the user runs `todo list --overdue`
THEN the output shows exactly one item: "Pay rent"
AND "Buy milk" does not appear (no due date)
AND "File taxes" does not appear (future due date)

## Rehearse

```rehearse
#!/bin/bash
todo add "Pay rent" --due 2026-01-01
todo add "Buy milk"
todo add "File taxes" --due 2026-12-31
output=$(todo list --overdue)
assert_contains "$output" "Pay rent"
assert_not_contains "$output" "Buy milk"
assert_not_contains "$output" "File taxes"
assert_contains "$output" "Showing 1 of 3 items"
```
