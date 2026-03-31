# Scenario: Completed todo with past due date is not overdue

**Validates:** due-dates/completed-not-overdue

## Steps

GIVEN a todo "Pay rent" with due date 2026-01-01
AND the todo is completed
WHEN the user runs `todo list --overdue`
THEN the output shows no items
AND "Pay rent" does not appear

## Rehearse

```rehearse
#!/bin/bash
todo add "Pay rent" --due 2026-01-01
todo complete 1
output=$(todo list --overdue)
assert_contains "$output" "No todos"
assert_not_contains "$output" "Pay rent"
```
