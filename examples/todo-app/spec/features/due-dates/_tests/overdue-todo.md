# Scenario: Overdue todo is flagged

**Validates:** due-dates/overdue-detection

## Steps

GIVEN a todo "Pay rent" with due date 2026-01-01
AND today is after 2026-01-01
WHEN the user runs `todo list`
THEN the output shows "Pay rent" with an `[OVERDUE]` tag

## Rehearse

```rehearse
#!/bin/bash
todo add "Pay rent" --due 2026-01-01
output=$(todo list)
assert_contains "$output" "Pay rent"
assert_contains "$output" "[OVERDUE]"
```

---
*This document follows the https://specscore.md/scenario-specification*
