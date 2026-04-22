# Scenario: Empty list shows message

**Validates:** todo-list/empty-list-message

## Steps

GIVEN an empty todo list
WHEN the user runs `todo list`
THEN the output contains "No todos"
AND the exit code is 0

## Rehearse

```rehearse
#!/bin/bash
output=$(todo list)
assert_exit_code 0
assert_contains "$output" "No todos"
```

---
*This document follows the https://specscore.md/scenario-specification*
