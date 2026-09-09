---
format: https://specscore.md/scenario-specification
---

# Scenario: External Plan store preserves source boundaries

**Validates:** [plan#req:external-source-namespace](../README.md#req-external-source-namespace), [plan#req:source-project-context](../README.md#req-source-project-context), [plan#req:read-only-plan-operations](../README.md#req-read-only-plan-operations), [plan#req:plan-path-conflicts](../README.md#req-plan-path-conflicts), [plan#req:external-store-lifecycle-lock](../README.md#req-external-store-lifecycle-lock)

## Steps

GIVEN source project `github.com/datatug/datatug` routes Plans to `sneat-co/workbench`
AND the destination checkout ignores `/.specscore-lifecycle.lock`
WHEN Plan `phase-1/core-loop` is created with `--project` naming the source checkout
THEN its document is `spec/plans/github.com/datatug/datatug/phase-1/core-loop/README.md` in the destination
AND Feature and acceptance-criteria references are resolved from the source checkout
AND no source artifact is modified

GIVEN the same external namespace contains a symlinked ancestor that resolves outside `spec/plans/github.com/datatug/datatug`
WHEN any Plan read or mutation follows that path
THEN it fails before reading or writing the escaped target

GIVEN a dry-run or read-only Plan operation
WHEN it completes or reports validation errors
THEN both repositories have the same bytes as before the command

GIVEN a parent Plan at `phase-1/README.md`
AND a child Plan at `phase-1/core-loop/README.md`
WHEN the source namespace index is regenerated
THEN only the namespace-root `README.md` is updated
AND the parent and child Plan documents are unchanged

GIVEN an existing lifecycle lock inode in the destination
WHEN Plan operations acquire and release the lock
THEN the inode remains present for later operations

---
*This document follows the https://specscore.md/scenario-specification*
