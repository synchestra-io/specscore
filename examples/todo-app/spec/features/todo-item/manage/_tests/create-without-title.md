# Scenario: Reject todo without title

**Validates:** todo-item/manage/title-required

## Steps

GIVEN an empty todo list
WHEN the user runs `todo add ""`
THEN the exit code is non-zero
AND the output contains an error message about the title
AND `todo list` shows no items

## Rehearse

```rehearse
#!/bin/bash
todo add ""
assert_exit_code 1
output=$(todo add "" 2>&1)
assert_contains "$output" "title"
list_output=$(todo list)
assert_contains "$list_output" "No todos"
```
