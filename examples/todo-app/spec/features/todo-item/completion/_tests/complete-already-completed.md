# Scenario: Complete an already-completed todo

**Validates:** todo-item/completion/completion-timestamp

## Steps

GIVEN a completed todo item titled "Buy milk"
WHEN the user runs `todo complete 1` again
THEN the exit code is 0
AND the item remains completed
AND the completion timestamp is unchanged

## Rehearse

```rehearse
#!/bin/bash
todo add "Buy milk"
todo complete 1
timestamp1=$(todo list --completed --format=json | jq -r '.[0].completed_at')
todo complete 1
assert_exit_code 0
timestamp2=$(todo list --completed --format=json | jq -r '.[0].completed_at')
assert_eq "$timestamp1" "$timestamp2"
```

---
*This document follows the https://specscore.md/scenario-specification*
