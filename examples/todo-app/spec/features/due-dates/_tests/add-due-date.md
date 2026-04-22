# Scenario: Add a todo with a due date

**Validates:** due-dates/due-date-optional

## Steps

GIVEN an empty todo list
WHEN the user runs `todo add "Pay rent" --due 2026-04-15`
THEN the exit code is 0
AND `todo list` shows "Pay rent" with due date "2026-04-15"

## Rehearse

```rehearse
#!/bin/bash
todo add "Pay rent" --due 2026-04-15
assert_exit_code 0
output=$(todo list)
assert_contains "$output" "Pay rent"
assert_contains "$output" "due: 2026-04-15"
```

---
*This document follows the https://specscore.md/scenario-specification*
