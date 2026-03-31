# Scenario: Rehearse block uses correct format

**Validates:** [scenario#req:rehearse-fence-tag](../README.md#req-rehearse-fence-tag), [scenario#req:rehearse-shebang](../README.md#req-rehearse-shebang), [scenario#req:rehearse-section](../README.md#req-rehearse-section)

## Steps

GIVEN a scenario file with a `## Rehearse` section containing:
  ````
  ```rehearse
  #!/bin/bash
  todo add "Buy milk"
  assert_exit_code 0
  ```
  ````
WHEN the spec validator checks the rehearse block
THEN the fence tag `rehearse` is recognized
AND the shebang `#!/bin/bash` is detected as the script language
AND the block is accepted as a valid rehearse script
